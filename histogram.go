package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import "errors"

func HistogramFind(buf []byte) ([]byte, error) {
	defer C.vips_thread_shutdown()

	if len(buf) == 0 {
		return nil, errors.New("Image buffer is empty")
	}

	image, _, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}

	imageHist, err := vipsHistogramFind(image)
	if err != nil {
		return nil, err
	}

	return getImageBuffer(imageHist)
}

func HistogramNorm(buf []byte) ([]byte, error) {
	defer C.vips_thread_shutdown()

	if len(buf) == 0 {
		return nil, errors.New("Image buffer is empty")
	}

	image, _, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}

	imageHist, err := vipsHistogramNorm(image)
	if err != nil {
		return nil, err
	}

	return getImageBuffer(imageHist)
}
