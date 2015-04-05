package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import (
	"math"
)

const (
	BICUBIC Interpolator = iota
	BILINEAR
	NOHALO
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

	image, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}

	// defaults
	if o.Quality == 0 {
		o.Quality = QUALITY
	}

	// get WxH
	inWidth := int(image.Xsize)
	inHeight := int(image.Ysize)

	if o.Crop {
		left, top := calculateCrop(inWidth, inHeight, o.Width, o.Height, o.Gravity)
		o.Width = int(math.Min(float64(inWidth), float64(o.Width)))
		o.Height = int(math.Min(float64(inHeight), float64(o.Height)))
		image, err = vipsExtract(image, left, top, o.Width, o.Height)
		if err != nil {
			return nil, err
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

	buf, err = vipsSave(image, vipsSaveOptions{Quality: o.Quality})
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
