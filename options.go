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
	D0   Angle = 0
	D90  Angle = 90
	D180 Angle = 180
	D270 Angle = 270
)

type Direction int

const (
	HORIZONTAL Direction = C.VIPS_DIRECTION_HORIZONTAL
	VERTICAL   Direction = C.VIPS_DIRECTION_VERTICAL
)

// Image interpretation type
// See: http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/VipsImage.html#VipsInterpretation
type Interpretation int

const (
	INTERPRETATION_ERROR     Interpretation = C.VIPS_INTERPRETATION_ERROR
	INTERPRETATION_MULTIBAND Interpretation = C.VIPS_INTERPRETATION_MULTIBAND
	INTERPRETATION_B_W       Interpretation = C.VIPS_INTERPRETATION_B_W
	INTERPRETATION_CMYK      Interpretation = C.VIPS_INTERPRETATION_CMYK
	INTERPRETATION_RGB       Interpretation = C.VIPS_INTERPRETATION_RGB
	INTERPRETATION_sRGB      Interpretation = C.VIPS_INTERPRETATION_sRGB
	INTERPRETATION_RGB16     Interpretation = C.VIPS_INTERPRETATION_RGB16
	INTERPRETATION_GREY16    Interpretation = C.VIPS_INTERPRETATION_GREY16
	INTERPRETATION_scRGB     Interpretation = C.VIPS_INTERPRETATION_scRGB
)

const WATERMARK_FONT = "sans 10"

// Color represents a traditional RGB color scheme
type Color struct {
	R, G, B uint8
}

type Watermark struct {
	Width       int
	DPI         int
	Margin      int
	Opacity     float32
	NoReplicate bool
	Text        string
	Font        string
	Background  Color
}

type Options struct {
	Height         int
	Width          int
	AreaHeight     int
	AreaWidth      int
	Top            int
	Left           int
	Extend         int
	Quality        int
	Compression    int
	Zoom           int
	Crop           bool
	Enlarge        bool
	Embed          bool
	Flip           bool
	Flop           bool
	NoAutoRotate   bool
	NoProfile      bool
	Interlace      bool
	Rotate         Angle
	Gravity        Gravity
	Watermark      Watermark
	Type           ImageType
	Interpolator   Interpolator
	Interpretation Interpretation
}
