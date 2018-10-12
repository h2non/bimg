// +build go1.7

package bimg

/*
 #cgo pkg-config: vips
 #include "vips/vips.h"
*/
import "C"

import (
	"errors"
	"runtime"
)

// Resize is used to transform a given image as byte buffer
// with the passed options.
func Resize(buf []byte, o Options) ([]byte, error) {
	// Required in order to prevent premature garbage collection. See:
	// https://github.com/h2non/bimg/pull/162
	defer runtime.KeepAlive(buf)
	return resizer(buf, o)
}

// AutoRotate is a directly calling vips.
func AutoRotate(buf []byte, o Options) ([]byte, error) {
	defer C.vips_thread_shutdown()
	image, imageType, err := loadImage(buf)
	if err != nil {
		return nil, err
	}
	// Clone and define default options
	o = applyDefaults(o, imageType)
	if !IsTypeSupported(o.Type) {
		return nil, errors.New("Unsupported image output type")
	}
	image, err = vipsAutoRotate(image)
	if err != nil {
		return nil, err
	}
	return saveImage(image, o)
}
