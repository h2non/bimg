package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

const (
	QUALITY  = 80
	MAX_SIZE = 16383
)

type Gravity int

const (
	CENTRE Gravity = iota
	NORTH
	EAST
	SOUTH
	WEST
)

type Interpolator int

const (
	BICUBIC Interpolator = iota
	BILINEAR
	NOHALO
)

var interpolations = map[Interpolator]string{
	BICUBIC:  "bicubic",
	BILINEAR: "bilinear",
	NOHALO:   "nohalo",
}

func (i Interpolator) String() string {
	return interpolations[i]
}

type Angle int

const (
	D0   Angle = C.VIPS_ANGLE_D0
	D90  Angle = C.VIPS_ANGLE_D90
	D180 Angle = C.VIPS_ANGLE_D180
	D270 Angle = C.VIPS_ANGLE_D270
)

type Direction int

const (
	HORIZONTAL Direction = C.VIPS_DIRECTION_HORIZONTAL
	VERTICAL   Direction = C.VIPS_DIRECTION_VERTICAL
)

type Options struct {
	Height       int
	Width        int
	AreaHeight   int
	AreaWidth    int
	Top          int
	Left         int
	Crop         bool
	Enlarge      bool
	Extend       int
	Embed        bool
	Quality      int
	Compression  int
	Type         ImageType
	Rotate       Angle
	Flip         Direction
	Gravity      Gravity
	Interpolator Interpolator
}
