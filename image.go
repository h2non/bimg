package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

type Image struct {
	buf   []byte
	image *C.struct__VipsImage
}
