package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

type Align int
const (
	Low Align = C.VIPS_ALIGN_LOW
	Center Align = C.VIPS_ALIGN_CENTRE
	High Align  = C.VIPS_ALIGN_HIGH
)

type BlendMode int
const (
	Clear		BlendMode = C.VIPS_BLEND_MODE_CLEAR
	Source	BlendMode = C.VIPS_BLEND_MODE_SOURCE
	Over		BlendMode = C.VIPS_BLEND_MODE_OVER
	In			BlendMode = C.VIPS_BLEND_MODE_IN
	Out		BlendMode = C.VIPS_BLEND_MODE_OUT
	Atop		BlendMode = C.VIPS_BLEND_MODE_ATOP
	Dest		BlendMode = C.VIPS_BLEND_MODE_DEST
	DestOver	BlendMode = C.VIPS_BLEND_MODE_DEST_OVER
	DestIn	BlendMode = C.VIPS_BLEND_MODE_DEST_IN
	DestOut	BlendMode = C.VIPS_BLEND_MODE_DEST_OUT
	DestAtop	BlendMode = C.VIPS_BLEND_MODE_DEST_ATOP
	Xor		BlendMode = C.VIPS_BLEND_MODE_XOR
	Add		BlendMode = C.VIPS_BLEND_MODE_ADD
	Saturate	BlendMode = C.VIPS_BLEND_MODE_SATURATE
	Mulitply	BlendMode = C.VIPS_BLEND_MODE_MULTIPLY
	Screen	BlendMode = C.VIPS_BLEND_MODE_SCREEN
	Overlay	BlendMode = C.VIPS_BLEND_MODE_OVERLAY
	Darken	BlendMode = C.VIPS_BLEND_MODE_DARKEN
	Lighten	BlendMode = C.VIPS_BLEND_MODE_LIGHTEN
	ColorDodge	BlendMode = C.VIPS_BLEND_MODE_COLOUR_DODGE
	ColorBurn	BlendMode = C.VIPS_BLEND_MODE_COLOUR_BURN
	HardLight	BlendMode = C.VIPS_BLEND_MODE_HARD_LIGHT
	SoftLight	BlendMode = C.VIPS_BLEND_MODE_SOFT_LIGHT
	Difference	BlendMode = C.VIPS_BLEND_MODE_DIFFERENCE
	Exclusion	BlendMode = C.VIPS_BLEND_MODE_EXCLUSION
)

type OptionsMulti struct {
	ArrayJoin		ArrayJoin
	Mosaic			Mosaic
	Composite		Composite
	Composite2		Composite2
}

type ArrayJoin struct {
	Num 		int
	Across	int
	Shim		int
	HSpacing int
	VSpacing int

	//TODO:
	HAlign	Align
	VAlign	Align
	Background	Color
}
/*
 * * @background: #VipsArrayDouble, background ink colour
 * * @halign: #VipsAlign, low, centre or high alignment
 * * @valign: #VipsAlign, low, centre or high alignment
*/

 type Mosaic struct {
	Direction 	Direction //TODO:
	Xref 			int
	Yref			int
	Xsec			int
	Ysec			int
 }

 type Composite struct {
	Mode	[]BlendMode
	Num	int
 }

 type Composite2 struct {
	Mode	BlendMode
 }