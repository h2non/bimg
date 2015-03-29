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

	//var tmpImage *C.struct__VipsImage
	/*
		// feed it
		imageLength := C.size_t(len(buf))
		imageBuf := unsafe.Pointer(&buf[0])
		debug("buffer: %s", buf[0])

		C.vips_jpegload_buffer_seq(imageBuf, imageLength, &image)
	*/

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
		//err := C.vips_extract_area_0(image, &tmpImage, C.int(left), C.int(top), C.int(o.Width), C.int(o.Height))
		//C.g_object_unref(C.gpointer(image))
		//image = tmpImage
	}

	// rotate
	r := Rotation{180}
	image, err = Rotate(image, r)
	if err != nil {
		return nil, err
	}

	// Finally save
	//var ptr unsafe.Pointer
	//length := C.size_t(0)

	//C.vips_jpegsave_custom(image, &ptr, &length, 1, C.int(o.Quality), 0)
	//C.g_object_unref(C.gpointer(image))
	//C.g_object_unref(C.gpointer(newImage))

	// get back the buffer
	//buf = C.GoBytes(ptr, C.int(length))
	// cleanup
	//C.g_free(C.gpointer(ptr))

	buf, err = vipsSave(image, vipsSaveOptions{Quality: o.Quality})
	if err != nil {
		return nil, err
	}
	C.vips_error_clear()

	return buf, nil
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
