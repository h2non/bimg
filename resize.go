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

func Resize(buf []byte, o Options) ([]byte, error) {
	// detect (if possible) the file type
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

	// crop
	if o.Crop {
		left, top := calculateCrop(inWidth, inHeight, o.Width, o.Height, o.Gravity)
		o.Width = int(math.Min(float64(inWidth), float64(o.Width)))
		o.Height = int(math.Min(float64(inHeight), float64(o.Height)))
		image, err = vipsExtract(image, left, top, o.Width, o.Height)
		if err != nil {
			return nil, err
		}
	}

	// rotate
	if o.Rotate > 0 {
		image, err = Rotate(image, Rotation{o.Rotate})
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

func Rotate(image *C.struct__VipsImage, r Rotation) (*C.struct__VipsImage, error) {
	//vips := &Vips{}
	return vipsRotate(image, r.calculate())
}

const (
	CENTRE Gravity = iota
	NORTH
	EAST
	SOUTH
	WEST
)

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
