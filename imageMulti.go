package bimg
/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import(
	"errors"
	"runtime"
)

type ImageMulti struct {
	buffer [][]byte
	outputOptions *Options
}

func NewImageFrom(inputBuffer [][]byte, outputOptions *Options) *ImageMulti {
	return &ImageMulti{inputBuffer, outputOptions}
}

func (x *ImageMulti) ArrayJoin(o ArrayJoin) ([]byte, error) {
	optionsMulti := OptionsMulti{ArrayJoin: o}
	return x.process(optionsMulti, "arrayjoin")
}

func (x *ImageMulti) Mosaic(o Mosaic) ([]byte, error) {
	optionsMulti := OptionsMulti{Mosaic: o}
	return x.process(optionsMulti, "mosaic")
}

func (x *ImageMulti) Composite(o Composite) ([]byte, error) {
	optionsMulti := OptionsMulti{Composite: o}
	return x.process(optionsMulti, "composite")
}

func (x *ImageMulti) Composite2(o Composite2) ([]byte, error) {
	optionsMulti := OptionsMulti{Composite2: o}
	return x.process(optionsMulti, "composite2")
}

func (x *ImageMulti) process(om OptionsMulti, process string) ([]byte, error) {
	// Required in order to prevent premature garbage collection. See:
	// https://github.com/h2non/bimg/pull/162
	defer runtime.KeepAlive(x.buffer)
	defer C.vips_thread_shutdown()

	imageType := x.outputOptions.Type
	if imageType == UNKNOWN {
		imageType = JPEG
	}
	
	// Clone and define default options
	o := applyDefaults(*x.outputOptions, imageType)

	// Ensure supported type
	if !IsTypeSupportedSave(imageType) {
		return nil, errors.New("Unsupported image output type")
	}

	var images []*C.VipsImage

	for _, buf := range x.buffer {
		img, _, err := vipsRead(buf)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	var err error
	var image *C.VipsImage
	
/*
// https://github.com/libvips/libvips/
	add 			Add two images 	vips_add()
	arrayjoin	Join an array of images	vips_arrayjoin()
	bandjoin 	Bandwise join a set of images 	vips_bandjoin(), vips_bandjoin2()
	bandrank 	Band-wise rank of a set of images 	vips_bandrank()
	boolean 		Boolean operation on two images 	vips_boolean(), vips_andimage(), vips_orimage(), vips_eorimage(), vips_lshift(), vips_rshift()
	case 			Use pixel values to pick cases from an array of images 	vips_case()
	complex2 	Complex binary operations on two images 	vips_complex2(), vips_cross_phase()
	complexform 	Form a complex image from two real images 	vips_complexform()
	composite 	Blend an array of images with an array of blend modes 	vips_composite()
	composite2 	Blend a pair of images with a blend mode 	vips_composite2()
	divide 		Divide two images 	vips_divide()
	draw_image 	Paint an image into another image 	vips_draw_image()
	join 			Join a pair of images 	vips_join()
	match 		First-order match of two images 	vips_match()
	merge 		Merge two images 	vips_merge()
	mosaic 		Mosaic two images 	vips_mosaic()
	mosaic1 		First-order mosaic of two images 	vips_mosaic1()
	multiply 	Multiply two images 	vips_multiply()
	remainder 	Remainder after integer division of two images 	vips_remainder()
	subtract 	Subtract two images 	vips_subtract()
	sum 			Sum an array of images 	vips_sum()
*/

	switch (process) {
	case "arrayjoin" :
		image, err = arrayjoin(images, om.ArrayJoin)

	case "mosaic" :
		image, err = mosaic(images, om.Mosaic)

	case "composite" :
		image, err = composite(images, om.Composite)

	case "composite2" :
		image, err = composite2(images, om.Composite2)

	default :
		return nil, errors.New("no multi image process found")
 	}

	if err != nil {
		return nil, err
	}

	return saveImage(image, o)
}

func arrayjoin(images []*C.VipsImage, o ArrayJoin) (*C.VipsImage, error) {
	o.Num = len(images)
	if o.Num <= 1 {
		return nil, errors.New("arrayjoin requires at least one image")
	}
	if (o.Across == 0) {
		o.Across = o.Num
	}
	
	return vipsArrayJoin(images, o)
}

func mosaic(images []*C.VipsImage, o Mosaic) (*C.VipsImage, error) {
	num := len(images)
	if num != 2 {
		return nil, errors.New("mosaic requires two images")
	}

	return vipsMosaic(images[0], images[1], o)
}

func composite(images []*C.VipsImage, o Composite) (*C.VipsImage, error) {
	o.Num = len(images)
	
	if len(o.Mode) != o.Num - 1 {
		return nil, errors.New("composite requires n - 1 modes")
	}

	return vipsComposite(images, o)
}

func composite2(images []*C.VipsImage, o Composite2) (*C.VipsImage, error) {	
	num := len(images)
	if num != 2 {
		return nil, errors.New("composite2 requires two images")
	}
	if o.Mode == 0 {
		o.Mode = Overlay
	}	

	return vipsComposite2(images[0], images[1], o)
}