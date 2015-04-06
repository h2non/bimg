package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import (
	"errors"
	"math"
)

const (
	CENTRE Gravity = iota
	NORTH
	EAST
	SOUTH
	WEST
)

func Resize(buf []byte, o Options) ([]byte, error) {
	defer C.vips_thread_shutdown()

	image, imageType, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}

	// defaults
	if o.Quality == 0 {
		o.Quality = QUALITY
	}
	if o.Compression == 0 {
		o.Compression = 6
	}
	if o.Type == 0 {
		o.Type = imageType
	}

	if IsTypeSupported(o.Type) == false {
		return nil, errors.New("Unsupported image output type")
	}

	// get WxH
	inWidth := int(image.Xsize)
	inHeight := int(image.Ysize)

	// prepare for factor
	factor := 0.0

	// image calculations
	switch {
	// Fixed width and height
	case o.Width > 0 && o.Height > 0:
		xf := float64(inWidth) / float64(o.Width)
		yf := float64(inHeight) / float64(o.Height)
		if o.Crop {
			factor = math.Min(xf, yf)
		} else {
			factor = math.Max(xf, yf)
		}
	// Fixed width, auto height
	case o.Width > 0:
		factor = float64(inWidth) / float64(o.Width)
		o.Height = int(math.Floor(float64(inHeight) / factor))
	// Fixed height, auto width
	case o.Height > 0:
		factor = float64(inHeight) / float64(o.Height)
		o.Width = int(math.Floor(float64(inWidth) / factor))
	// Identity transform
	default:
		factor = 1
		o.Width = inWidth
		o.Height = inHeight
	}

	debug("transform from %dx%d to %dx%d", inWidth, inHeight, o.Width, o.Height)

	shrink := int(math.Max(math.Floor(factor), 1))
	residual := float64(shrink) / factor

	// Do not enlarge the output if the input width *or* height are already less than the required dimensions
	if o.Enlarge == false {
		if inWidth < o.Width && inHeight < o.Height {
			factor = 1
			shrink = 1
			residual = 0
			o.Width = inWidth
			o.Height = inHeight
		}
	}

	debug("factor: %v, shrink: %v, residual: %v", factor, shrink, residual)

	// Try to use libjpeg shrink-on-load
	shrinkOnLoad := 1
	if imageType == JPEG && shrink >= 2 {
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

		if shrinkOnLoad > 1 {
			// Recalculate integral shrink and double residual
			factor = math.Max(factor, 1.0)
			shrink = int(math.Floor(factor))
			residual = float64(shrink) / factor
			// Reload input using shrink-on-load
			image, err = vipsShrinkJpeg(buf, shrinkOnLoad)
			if err != nil {
				return nil, err
			}
		}
	}

	if shrink > 1 {
		debug("shrink %d", shrink)
		// Use vips_shrink with the integral reduction
		image, err = vipsShrink(image, shrink)
		if err != nil {
			return nil, err
		}

		// Recalculate residual float based on dimensions of required vs shrunk images
		shrunkWidth := int(image.Xsize)
		shrunkHeight := int(image.Ysize)

		residualx := float64(o.Width) / float64(shrunkWidth)
		residualy := float64(o.Height) / float64(shrunkHeight)
		if o.Crop {
			residual = math.Max(residualx, residualy)
		} else {
			residual = math.Min(residualx, residualy)
		}
	}

	// Use vips_affine with the remaining float part
	if residual != 0 {
		debug("residual %.2f", residual)
		image, err = vipsAffine(image, residual, o.Interpolator)
		if err != nil {
			return nil, err
		}
	}

	// Crop/embed
	affinedWidth := int(image.Xsize)
	affinedHeight := int(image.Ysize)

	if affinedWidth != o.Width || affinedHeight != o.Height {
		if o.Crop {
			left, top := calculateCrop(affinedWidth, affinedHeight, o.Width, o.Height, o.Gravity)
			o.Width = int(math.Min(float64(inWidth), float64(o.Width)))
			o.Height = int(math.Min(float64(inHeight), float64(o.Height)))
			debug("crop image to %dx%d to %dx%d", left, top, o.Width, o.Height)
			image, err = vipsExtract(image, left, top, o.Width, o.Height)
			if err != nil {
				return nil, err
			}
		} else if o.Embed {
			left := (o.Width - affinedWidth) / 2
			top := (o.Height - affinedHeight) / 2
			image, err = vipsEmbed(image, left, top, o.Width, o.Height, o.Extend)
			if err != nil {
				return nil, err
			}
		}
	}

	rotation, flip := calculateRotationAndFlip(image, o.Rotate)
	if flip {
		o.Flip = HORIZONTAL
	}
	if rotation != D0 {
		o.Rotate = rotation
	}

	if o.Rotate > 0 {
		rotation := calculateRotation(o.Rotate)
		image, err = vipsRotate(image, rotation)
		if err != nil {
			return nil, err
		}
	}

	if o.Flip > 0 {
		image, err = vipsFlip(image, o.Flip)
		if err != nil {
			return nil, err
		}
	}

	saveOptions := vipsSaveOptions{
		Quality:     o.Quality,
		Type:        o.Type,
		Compression: o.Compression,
	}

	buf, err = vipsSave(image, saveOptions)
	if err != nil {
		return nil, err
	}
	C.vips_error_clear()

	return buf, nil
}

func calculateCrop(inWidth, inHeight, outWidth, outHeight int, gravity Gravity) (int, int) {
	left, top := 0, 0

	switch gravity {
	case NORTH:
		left = (inWidth - outWidth + 1) / 2
	case EAST:
		left = inWidth - outWidth
		top = (inHeight - outHeight + 1) / 2
	case SOUTH:
		left = (inWidth - outWidth + 1) / 2
		top = inHeight - outHeight
	case WEST:
		top = (inHeight - outHeight + 1) / 2
	default:
		left = (inWidth - outWidth + 1) / 2
		top = (inHeight - outHeight + 1) / 2
	}

	return left, top
}

func calculateRotationAndFlip(image *C.struct__VipsImage, angle Angle) (Angle, bool) {
	rotate := D0
	flip := false

	if angle == -1 {
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
			rotate = D90
			break // flip 6
		case 4:
			flip = true
			rotate = D180
			break // flip 3
		case 5:
			flip = true
			rotate = D270
			break // flip 8
		}
	} else {
		if angle == 90 {
			rotate = D90
		} else if angle == 180 {
			rotate = D180
		} else if angle == 270 {
			rotate = D270
		}
	}

	return rotate, flip
}

func calculateRotation(angle Angle) Angle {
	divisor := angle % 90
	if divisor != 0 {
		angle = angle - divisor
	}
	return angle
}
