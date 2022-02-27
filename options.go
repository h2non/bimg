package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

const (
	// Quality defines the default JPEG quality to be used.
	Quality = 75
	// MaxSize defines the maximum pixels width or height supported.
	MaxSize = 16383
)

// Gravity represents the image gravity value.
type Gravity int

const (
	// GravityCentre represents the centre value used for image gravity orientation.
	GravityCentre Gravity = iota
	// GravityNorth represents the north value used for image gravity orientation.
	GravityNorth
	// GravityEast represents the east value used for image gravity orientation.
	GravityEast
	// GravitySouth represents the south value used for image gravity orientation.
	GravitySouth
	// GravityWest represents the west value used for image gravity orientation.
	GravityWest
	// GravitySmart enables libvips Smart Crop algorithm for image gravity orientation.
	GravitySmart
)

// Interpolator represents the image interpolation value.
type Interpolator int

const (
	// Bicubic interpolation value.
	Bicubic Interpolator = iota
	// Bilinear interpolation value.
	Bilinear
	// Nohalo interpolation value.
	Nohalo
	// Nearest neighbour interpolation value.
	Nearest
)

var interpolations = map[Interpolator]string{
	Bicubic:  "bicubic",
	Bilinear: "bilinear",
	Nohalo:   "nohalo",
	Nearest:  "nearest",
}

func (i Interpolator) String() string {
	return interpolations[i]
}

// Direction represents the image direction value.
type Direction int

const (
	// Horizontal represents the orizontal image direction value.
	Horizontal Direction = C.VIPS_DIRECTION_HORIZONTAL
	// Vertical represents the vertical image direction value.
	Vertical Direction = C.VIPS_DIRECTION_VERTICAL
)

// Interpretation represents the image interpretation type.
// See: https://libvips.github.io/libvips/API/current/VipsImage.html#VipsInterpretation
type Interpretation int

const (
	// InterpretationError points to the libvips interpretation error type.
	InterpretationError Interpretation = C.VIPS_INTERPRETATION_ERROR
	// InterpretationMultiband points to its libvips interpretation equivalent type.
	InterpretationMultiband Interpretation = C.VIPS_INTERPRETATION_MULTIBAND
	// InterpretationBW points to its libvips interpretation equivalent type.
	InterpretationBW Interpretation = C.VIPS_INTERPRETATION_B_W
	// InterpretationCMYK points to its libvips interpretation equivalent type.
	InterpretationCMYK Interpretation = C.VIPS_INTERPRETATION_CMYK
	// InterpretationRGB points to its libvips interpretation equivalent type.
	InterpretationRGB Interpretation = C.VIPS_INTERPRETATION_RGB
	// InterpretationSRGB points to its libvips interpretation equivalent type.
	InterpretationSRGB Interpretation = C.VIPS_INTERPRETATION_sRGB
	// InterpretationRGB16 points to its libvips interpretation equivalent type.
	InterpretationRGB16 Interpretation = C.VIPS_INTERPRETATION_RGB16
	// InterpretationGREY16 points to its libvips interpretation equivalent type.
	InterpretationGREY16 Interpretation = C.VIPS_INTERPRETATION_GREY16
	// InterpretationScRGB points to its libvips interpretation equivalent type.
	InterpretationScRGB Interpretation = C.VIPS_INTERPRETATION_scRGB
	// InterpretationLAB points to its libvips interpretation equivalent type.
	InterpretationLAB Interpretation = C.VIPS_INTERPRETATION_LAB
	// InterpretationXYZ points to its libvips interpretation equivalent type.
	InterpretationXYZ Interpretation = C.VIPS_INTERPRETATION_XYZ
)

// Extend represents the image extend mode, used when the edges
// of an image are extended, you can specify how you want the extension done.
// See: https://libvips.github.io/libvips/API/current/libvips-conversion.html#VIPS-EXTEND-BACKGROUND:CAPS
type Extend int

const (
	// ExtendBlack extend with black (all 0) pixels mode.
	ExtendBlack Extend = C.VIPS_EXTEND_BLACK
	// ExtendCopy copy the image edges.
	ExtendCopy Extend = C.VIPS_EXTEND_COPY
	// ExtendRepeat repeat the whole image.
	ExtendRepeat Extend = C.VIPS_EXTEND_REPEAT
	// ExtendMirror mirror the whole image.
	ExtendMirror Extend = C.VIPS_EXTEND_MIRROR
	// ExtendWhite extend with white (all bits set) pixels.
	ExtendWhite Extend = C.VIPS_EXTEND_WHITE
	// ExtendBackground with colour from the background property.
	ExtendBackground Extend = C.VIPS_EXTEND_BACKGROUND
	// ExtendLast extend with last pixel.
	ExtendLast Extend = C.VIPS_EXTEND_LAST
)

// WatermarkFont defines the default watermark font to be used.
var WatermarkFont = "sans 10"

// Color represents a traditional RGB color scheme.
type Color struct {
	R, G, B uint8
}

// RGBA returns the color with full opacity.
func (c Color) RGBA() (r uint8, g uint8, b uint8, a uint8) {
	return c.R, c.G, c.B, 0xFF
}

// ColorWithAlpha represents a traditional RGBA color scheme.
// It uses Color as base for easier reusability of already defined colors.
type ColorWithAlpha struct {
	Color
	A uint8
}

// RGBA returns the color with the specified alpha value.
func (c ColorWithAlpha) RGBA() (r uint8, g uint8, b uint8, a uint8) {
	return c.R, c.G, c.B, c.A
}

// RGBAProvider defines the interface required to get RGBA compatible
// color channels.
type RGBAProvider interface {
	RGBA() (r uint8, g uint8, b uint8, a uint8)
}

// ColorBlack is a shortcut to black RGB color representation.
var ColorBlack = Color{0x00, 0x00, 0x00}

// ColorWhite is a shortcut to white RGB color representation.
var ColorWhite = Color{0xFF, 0xFF, 0xFF}

// WatermarkOptions represents the text-based watermark supported options.
type WatermarkOptions struct {
	Width       int
	DPI         int
	Margin      int
	Opacity     float32
	NoReplicate bool
	Text        string
	Font        string
	Background  RGBAProvider
}

// WatermarkImage represents the image-based watermark supported options.
type WatermarkImage struct {
	Left    int
	Top     int
	Buf     []byte
	Opacity float32
}

// GaussianBlurOptions represents the gaussian image transformation values.
type GaussianBlurOptions struct {
	Sigma   float64
	MinAmpl float64
}

// SharpenOptions represents the image sharp transformation options.
type SharpenOptions struct {
	Radius int
	X1     float64
	Y2     float64
	Y3     float64
	M1     float64
	M2     float64
}
