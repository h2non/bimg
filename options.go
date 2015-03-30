package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

const QUALITY = 80

type Gravity int

type Interpolator int

var interpolations = map[Interpolator]string{
	BICUBIC:  "bicubic",
	BILINEAR: "bilinear",
	NOHALO:   "nohalo",
}

func (i Interpolator) String() string {
	return interpolations[i]
}

type Rotation struct {
	angle int
}

func (a Rotation) calculate() int {
	angle := a.angle
	divisor := angle % 90
	if divisor != 0 {
		angle = a.angle - divisor
	}
	return angle
}

type Direction int

const (
	HORIZONTAL Direction = C.VIPS_DIRECTION_HORIZONTAL
	VERTICAL   Direction = C.VIPS_DIRECTION_VERTICAL
)

type Options struct {
	Height       int
	Width        int
	Crop         bool
	Enlarge      bool
	Extend       int
	Embed        bool
	Quality      int
	Rotate       int
	Flip         Direction
	Gravity      Gravity
	Interpolator Interpolator
}
