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

// resizer is used to transform a given image as byte buffer
// with the passed options.
func resizer(buf []byte, o Options) ([]byte, error) {
	defer C.vips_thread_shutdown()

	t, err := NewImageTransformation(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot load image: %w", err)
	}

	// Clone and define default options
	o = applyDefaults(o, t.imageType)
	saveOptions := buildSaveOptions(o)

	// Ensure supported type
	if !IsTypeSupportedSave(o.Type) {
		return nil, errors.New("unsupported image output type")
	}

	// Autorotate only
	if o.autoRotateOnly {
		if err := t.Rotate(RotateOptions{}); err != nil {
			return nil, err
		}
		return t.Save(saveOptions)
	}

	// Auto rotate image based on EXIF orientation header
	if err := t.Rotate(RotateOptions{
		Angle:        o.Rotate,
		Flip:         o.Flip,
		Flop:         o.Flop,
		NoAutoRotate: o.NoAutoRotate,
	}); err != nil {
		return nil, fmt.Errorf("cannot rotate image: %w", err)
	}

	if err := resizeIfNecessary(t, o); err != nil {
		return nil, err
	}

	if err := extractOrEmbedImage(t, o); err != nil {
		return nil, fmt.Errorf("cannot extract or embed: %w", err)
	}

	// Apply effects, if necessary
	if o.GaussianBlur.Sigma > 0 || o.GaussianBlur.MinAmpl > 0 {
		if err := t.Blur(o.GaussianBlur); err != nil {
			return nil, fmt.Errorf("cannot apply blur: %w", err)
		}
	}
	if o.Sharpen.Radius > 0 && o.Sharpen.Y2 > 0 || o.Sharpen.Y3 > 0 {
		if err := t.Sharpen(o.Sharpen); err != nil {
			return nil, fmt.Errorf("cannot sharpen image: %w", err)
		}
	}

	// Add watermark, if necessary
	if err := t.WatermarkText(o.Watermark); err != nil {
		return nil, fmt.Errorf("cannot add watermark text: %w", err)
	}

	// Add watermark, if necessary
	if len(o.WatermarkImage.Buf) > 0 {
		watermarkImage, err := NewImageTransformation(o.WatermarkImage.Buf)
		if err != nil {
			return nil, fmt.Errorf("cannot decode watermark image: %w", err)
		}
		if err := t.WatermarkImage(WatermarkImageOptions{
			Left:    o.WatermarkImage.Left,
			Top:     o.WatermarkImage.Top,
			Image:   watermarkImage,
			Opacity: o.WatermarkImage.Opacity,
		}); err != nil {
			return nil, fmt.Errorf("cannot add watermark image: %w", err)
		}
	}

	// Flatten image on a background, if necessary
	if shouldFlatten(o) {
		if err := t.Flatten(o.Background); err != nil {
			return nil, fmt.Errorf("cannot flatten image: %w", err)
		}
	}

	// Apply Gamma filter, if necessary
	if o.Gamma > 0 {
		if err := t.Gamma(o.Gamma); err != nil {
			return nil, fmt.Errorf("cannot apply gamma: %w", err)
		}
	}

	return t.Save(saveOptions)
}

func resizeIfNecessary(t *ImageTransformation, o Options) error {
	imageSize := t.Size()

	if o.Crop || o.Embed {
		if o.Width == 0 {
			o.Width = imageSize.Width
		}
		if o.Height == 0 {
			o.Height = imageSize.Height
		}

		if !o.Force && !o.Enlarge && imageSize.Width < o.Width && imageSize.Height < o.Height {
			// Nothing to do.
			return nil
		}
	}

	var resizeMode ResizeMode

	switch {
	case o.Force:
		resizeMode = ResizeModeForce
	case o.Crop, o.Embed, o.Trim:
		resizeMode = ResizeModeFitUp
	case o.Zoom > 0, o.Left > 0, o.Top > 0, o.AreaHeight > 0, o.AreaWidth > 0:
		resizeMode = ResizeModeFit
	default:
		resizeMode = ResizeModeForce
	}

	if err := t.Resize(ResizeOptions{
		Height:         o.Height,
		Width:          o.Width,
		Top:            o.Top,
		Left:           o.Left,
		Zoom:           o.Zoom,
		Mode:           resizeMode,
		Interpolator:   o.Interpolator,
		Interpretation: o.Interpretation,
	}); err != nil {
		return fmt.Errorf("cannot resize image: %w", err)
	}
	return nil
}

func loadImage(buf []byte) (*vipsImage, ImageType, error) {
	if len(buf) == 0 {
		return nil, JPEG, errors.New("Image buffer is empty")
	}

	image, imageType, err := vipsRead(buf)
	if err != nil {
		return nil, JPEG, err
	}

	return image, imageType, nil
}

func applyDefaults(o Options, imageType ImageType) Options {
	if o.Quality == 0 {
		o.Quality = Quality
	}
	if o.Compression == 0 {
		o.Compression = 6
	}
	if o.Type == 0 {
		o.Type = imageType
	}
	if o.Interpretation == 0 {
		o.Interpretation = InterpretationSRGB
	}
	if o.SmartCrop {
		o.Gravity = GravitySmart
	}
	return o
}

func buildSaveOptions(o Options) SaveOptions {
	return SaveOptions{
		Quality:        o.Quality,
		Type:           o.Type,
		Compression:    o.Compression,
		Interlace:      o.Interlace,
		NoProfile:      o.NoProfile,
		Interpretation: o.Interpretation,
		InputICC:       o.InputICC,
		OutputICC:      o.OutputICC,
		StripMetadata:  o.StripMetadata,
		Lossless:       o.Lossless,
		Palette:        o.Palette,
		Speed:          o.Speed,
	}
}

func transformImage(image *vipsImage, o ResizeOptions, shrink int, residual float64) (*vipsImage, error) {
	var err error
	// Use vips_shrink with the integral reduction
	if shrink > 1 {
		image, residual, err = shrinkImage(image, o, residual, shrink)
		if err != nil {
			return nil, err
		}
	}

	residualx, residualy := residual, residual
	if o.Mode == ResizeModeForce {
		residualx = float64(o.Width) / float64(image.c.Xsize)
		residualy = float64(o.Height) / float64(image.c.Ysize)
	}

	if o.Mode == ResizeModeForce || residual != 0 {
		if residualx < 1 && residualy < 1 {
			image, err = vipsReduce(image, 1/residualx, 1/residualy)
		} else {
			// TODO extend configurable?
			image, err = vipsAffine(image, residualx, residualy, o.Interpolator, ExtendBlack)
		}
		if err != nil {
			return nil, err
		}
	}

	return image, nil
}

func extractOrEmbedImage(it *ImageTransformation, o Options) error {
	if o.Force {
		o.Crop = false
		o.Embed = false
	}
	// If (one of) the size values are not given, assume they stay
	// as in the current state of transformation.
	if o.Width == 0 {
		o.Width = int(it.image.c.Xsize)
	}
	if o.Height == 0 {
		o.Height = int(it.image.c.Ysize)
	}

	switch {
	case o.Crop || o.Gravity == GravitySmart:
		return it.Crop(CropOptions{
			Width:   o.Width,
			Height:  o.Height,
			Gravity: o.Gravity,
		})
	case o.Embed:
		return it.Embed(EmbedOptions{
			Width:      o.Width,
			Height:     o.Height,
			Extend:     o.Extend,
			Background: o.Background,
		})
	case o.Trim:
		return it.Trim(TrimOptions{
			Background: o.Background,
			Threshold:  o.Threshold,
		})
	case o.Top != 0 || o.Left != 0 || o.AreaWidth != 0 || o.AreaHeight != 0:
		if o.AreaWidth == 0 {
			o.AreaWidth = o.Width
		}
		if o.AreaHeight == 0 {
			o.AreaHeight = o.Height
		}
		return it.Extract(ExtractOptions{
			Left:   o.Left,
			Top:    o.Top,
			Width:  o.AreaWidth,
			Height: o.AreaHeight,
		})
	}

	return nil
}

func rotateAndFlipImage(image *vipsImage, o RotateOptions) (*vipsImage, error) {
	var err error

	if o.Angle > 0 {
		image, err = vipsRotate(image, getAngle(o.Angle))
	}

	if o.Flip {
		image, err = vipsFlip(image, Horizontal)
	}

	if o.Flop {
		image, err = vipsFlip(image, Vertical)
	}

	return image, err
}

func watermarkImageWithText(image *vipsImage, w Watermark) (*vipsImage, error) {
	if w.Text == "" {
		return image, nil
	}

	// Defaults
	if w.Font == "" {
		w.Font = WatermarkFont
	}
	if w.Width == 0 {
		w.Width = int(math.Floor(float64(image.c.Xsize / 6)))
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

func shouldFlatten(o Options) bool {
	// If no background is set, we cannot flatten anything. Just skip.
	if o.Background == nil {
		return false
	}

	// If an alpha channel is set, but is not full opacity, we should not flatten, since it would
	// purge the alpha channel.
	_, _, _, a := o.Background.RGBA()
	if a < 0xFF {
		return false
	}

	return true
}

func zoomImage(image *vipsImage, zoom int) (*vipsImage, error) {
	if zoom == 0 {
		return image, nil
	}
	return vipsZoom(image, zoom+1)
}

func shrinkImage(image *vipsImage, o ResizeOptions, residual float64, shrink int) (*vipsImage, float64, error) {
	// Use vips_shrink with the integral reduction
	image, err := vipsShrink(image, shrink)
	if err != nil {
		return nil, 0, err
	}

	// Recalculate residual float based on dimensions of required vs shrunk images
	residualx := float64(o.Width) / float64(image.c.Xsize)
	residualy := float64(o.Height) / float64(image.c.Ysize)

	if o.Mode == ResizeModeFit {
		residual = math.Max(residualx, residualy)
	} else {
		residual = math.Min(residualx, residualy)
	}

	return image, residual, nil
}

func shrinkOnLoad(buf []byte, imageType ImageType, factor float64, shrink int) (*vipsImage, float64, error) {
	var image *vipsImage
	var err error

	// Reload input using shrink-on-load
	if imageType == JPEG && shrink >= 2 {
		shrinkOnLoad := 1
		// Recalculate integral shrink and double residual
		switch {
		case shrink >= 8:
			factor = factor / 8
			shrinkOnLoad = 8
		case shrink >= 4:
			factor = factor / 4
			shrinkOnLoad = 4
		case shrink >= 2:
			factor = factor / 2
			shrinkOnLoad = 2
		}

		image, err = vipsShrinkJpeg(buf, shrinkOnLoad)
	} else if imageType == WEBP {
		image, err = vipsShrinkWebp(buf, shrink)
	} else {
		return nil, 0, fmt.Errorf("%v doesn't support shrink on load", ImageTypeName(imageType))
	}

	return image, factor, err
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

func calculateRotationAndFlip(image *vipsImage, angle Angle) (Angle, bool) {
	rotate := D0
	flip := false

	if angle > 0 {
		return rotate, flip
	}

	switch vipsExifOrientation(image) {
	case 6:
		rotate = D90
		break
	case 3:
		rotate = D180
		break
	case 8:
		rotate = D270
		break
	case 2:
		flip = true
		break // flip 1
	case 7:
		flip = true
		rotate = D270
		break // flip 6
	case 4:
		flip = true
		rotate = D180
		break // flip 3
	case 5:
		flip = true
		rotate = D90
		break // flip 8
	}

	return rotate, flip
}

func calculateShrink(factor float64, i Interpolator) int {
	var shrink float64

	// Calculate integral box shrink
	windowSize := vipsWindowSize(i.String())
	if factor >= 2 && windowSize > 3 {
		// Shrink less, affine more with interpolators that use at least 4x4 pixel window, e.g. bicubic
		shrink = float64(math.Floor(factor * 3.0 / windowSize))
	} else {
		shrink = math.Floor(factor)
	}

	return int(math.Max(shrink, 1))
}

func calculateResidual(factor float64, shrink int) float64 {
	return float64(shrink) / factor
}

func getAngle(angle Angle) Angle {
	divisor := angle % 90
	if divisor != 0 {
		angle = angle - divisor
	}
	return Angle(math.Min(float64(angle), 270))
}
