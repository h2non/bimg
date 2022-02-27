package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrExtractAreaParamsRequired defines a generic extract area error
	ErrExtractAreaParamsRequired = errors.New("extract area width/height params are required")
)

func resizeImage(image *vipsImage, o ResizeOptions) (*vipsImage, error) {
	xscale := float64(o.Width) / float64(image.c.Xsize)
	yscale := float64(o.Height) / float64(image.c.Ysize)

	return vipsResize(image, xscale, yscale)
}

func watermarkImageWithText(image *vipsImage, w WatermarkOptions) (*vipsImage, error) {
	if w.Text == "" {
		return image, nil
	}

	// Defaults
	if w.Font == "" {
		w.Font = WatermarkFont
	}
	if w.Width == 0 {
		w.Width = int(image.c.Xsize) / 6
	}
	if w.DPI == 0 {
		w.DPI = 150
	}
	if w.Margin == 0 {
		w.Margin = w.Width
	}
	if w.Opacity == 0 {
		w.Opacity = 0.25
	} else if w.Opacity > 1 {
		w.Opacity = 1
	}

	image, err := vipsWatermark(image, w)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func watermarkImageWithAnotherImage(image *vipsImage, w WatermarkImageOptions) (*vipsImage, error) {
	if w.Image == nil || w.Image.image == nil {
		return image, errors.New("no image to watermark with given")
	}

	if w.Opacity == 0.0 {
		w.Opacity = 1.0
	}

	image, err := vipsDrawWatermark(image, drawWatermarkOptions{
		Left:    w.Left,
		Top:     w.Top,
		Image:   w.Image.image,
		Opacity: w.Opacity,
	})

	if err != nil {
		return nil, err
	}

	return image, nil
}

func zoomImage(image *vipsImage, zoom int) (*vipsImage, error) {
	if zoom == 0 {
		return image, nil
	}
	return vipsZoom(image, zoom+1)
}

func shrinkOnLoad(buf []byte, imageType ImageType, factor float64, shrink float64) (*vipsImage, error) {
	var image *vipsImage
	var err error

	if shrink < 2 {
		return nil, fmt.Errorf("only available for shrink >=2")
	}

	shrinkOnLoad := 1
	// Recalculate integral shrink and double residual
	switch {
	case shrink >= 8:
		shrinkOnLoad = 8
	case shrink >= 4:
		shrinkOnLoad = 4
	case shrink >= 2:
		shrinkOnLoad = 2
	}

	// Reload input using shrink-on-load
	switch imageType {
	case JPEG:
		image, err = vipsShrinkJpeg(buf, shrinkOnLoad)
	case WEBP:
		image, err = vipsShrinkWebp(buf, shrinkOnLoad)
	default:
		return nil, fmt.Errorf("%v doesn't support shrink on load", ImageTypeName(imageType))
	}

	return image, err
}

func roundFloat(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

func calculateCrop(inWidth, inHeight, outWidth, outHeight int, gravity Gravity) (int, int) {
	left, top := 0, 0

	switch gravity {
	case GravityNorth:
		left = (inWidth - outWidth + 1) / 2
	case GravityEast:
		left = inWidth - outWidth
		top = (inHeight - outHeight + 1) / 2
	case GravitySouth:
		left = (inWidth - outWidth + 1) / 2
		top = inHeight - outHeight
	case GravityWest:
		top = (inHeight - outHeight + 1) / 2
	default:
		left = (inWidth - outWidth + 1) / 2
		top = (inHeight - outHeight + 1) / 2
	}

	return left, top
}
func calculateShrink(factor float64, i Interpolator) float64 {
	var shrink float64

	// Calculate integral box shrink
	windowSize := vipsWindowSize(i.String())
	if factor >= 2 && windowSize > 3 {
		// Shrink less, affine more with interpolators that use at least 4x4 pixel window, e.g. bicubic
		shrink = float64(math.Floor(factor * 3.0 / windowSize))
	} else {
		shrink = math.Floor(factor)
	}

	return math.Max(shrink, 1)
}
