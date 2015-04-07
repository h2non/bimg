package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
)

func Resize(buf []byte, o Options) ([]byte, error) {
	defer C.vips_thread_shutdown()

	if len(buf) == 0 {
		return nil, errors.New("Image buffer is empty")
	}

	image, imageType, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}

	// Defaults
	if o.Quality == 0 {
		o.Quality = QUALITY
	}
	if o.Compression == 0 {
		o.Compression = 6
	}
	if o.Type == 0 {
		o.Type = imageType
	}

	if IsTypeSupported(o.Type) == false {
		return nil, errors.New("Unsupported image output type")
	}

	debug("Options: %#v", o)

	inWidth := int(image.Xsize)
	inHeight := int(image.Ysize)

	// image calculations
	factor := imageCalculations(&o, inWidth, inHeight)
	shrink := int(math.Max(math.Floor(factor), 1))
	residual := float64(shrink) / factor

	// Do not enlarge the output if the input width *or* height are already less than the required dimensions
	if o.Enlarge == false {
		if inWidth < o.Width && inHeight < o.Height {
			factor = 1.0
			shrink = 1
			residual = 0
			o.Width = inWidth
			o.Height = inHeight
		}
	}

	// Try to use libjpeg shrink-on-load
	if imageType == JPEG && shrink >= 2 {
		tmpImage, factor, err := shrinkJpegImage(buf, image, factor, shrink)
		if err != nil {
			return nil, err
		}
		if tmpImage != nil {
			image = tmpImage
			factor = math.Max(factor, 1.0)
			shrink = int(math.Floor(factor))
			residual = float64(shrink) / factor
		}
	}

	// Calculate integral box shrink
	windowSize := vipsWindowSize(o.Interpolator.String())
	if factor >= 2 && windowSize > 3 {
		// Shrink less, affine more with interpolators that use at least 4x4 pixel window, e.g. bicubic
		shrink = int(math.Max(float64(math.Floor(factor*3.0/windowSize)), 1))
	}

	// Transform image if necessary
	shouldTransform := o.Width != inWidth || o.Height != inHeight || o.AreaWidth > 0 || o.AreaHeight > 0
	if shouldTransform {
		// Use vips_shrink with the integral reduction
		if shrink > 1 {
			image, residual, err = shrinkImage(image, o, residual, shrink)
			if err != nil {
				return nil, err
			}
		}
		// Use vips_affine with the remaining float part
		if residual != 0 {
			image, err = vipsAffine(image, residual, o.Interpolator)
			if err != nil {
				return nil, err
			}
		}

		debug("Transform image: factor=%v, shrink=%v, residual=%v", factor, shrink, residual)
		// Extract area from image
		image, err = extractImage(image, o)
		if err != nil {
			return nil, err
		}
	}

	// Rotate / flip image if necessary based on EXIF metadata
	image, err = rotateImage(image, o)
	if err != nil {
		return nil, err
	}

	saveOptions := vipsSaveOptions{
		Quality:     o.Quality,
		Type:        o.Type,
		Compression: o.Compression,
	}

	// Finally save as buffer
	buf, err = vipsSave(image, saveOptions)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func extractImage(image *C.struct__VipsImage, o Options) (*C.struct__VipsImage, error) {
	var err error = nil
	inWidth := int(image.Xsize)
	inHeight := int(image.Ysize)

	switch {
	case o.Crop:
		width := int(math.Min(float64(inWidth), float64(o.Width)))
		height := int(math.Min(float64(inHeight), float64(o.Height)))
		left, top := calculateCrop(inWidth, inHeight, o.Width, o.Height, o.Gravity)
		image, err = vipsExtract(image, left, top, width, height)
		break
	case o.Embed:
		left, top := (o.Width-inWidth)/2, (o.Height-inHeight)/2
		image, err = vipsEmbed(image, left, top, o.Width, o.Height, o.Extend)
		break
	case o.Top > 0 || o.Left > 0:
		if o.AreaWidth == 0 || o.AreaHeight == 0 {
			err = errors.New(fmt.Sprintf("Invalid area to extract %dx%d", o.AreaWidth, o.AreaHeight))
		} else {
			image, err = vipsExtract(image, o.Left, o.Top, o.AreaWidth, o.AreaHeight)
		}
		break
	}

	return image, err
}

func rotateImage(image *C.struct__VipsImage, o Options) (*C.struct__VipsImage, error) {
	var err error
	var direction Direction = -1

	rotation, flip := calculateRotationAndFlip(image, o.Rotate)
	if flip {
		o.Flip = flip
	}
	if rotation > D0 && o.Rotate == 0 {
		o.Rotate = rotation
	}

	if o.Rotate > 0 {
		image, err = vipsRotate(image, getAngle(o.Rotate))
	}

	if o.Flip {
		direction = HORIZONTAL
	} else if o.Flop {
		direction = VERTICAL
	}

	if direction != -1 {
		image, err = vipsFlip(image, direction)
	}

	return image, err
}

func shrinkImage(image *C.struct__VipsImage, o Options, residual float64, shrink int) (*C.struct__VipsImage, float64, error) {
	// Use vips_shrink with the integral reduction
	image, err := vipsShrink(image, shrink)
	if err != nil {
		return nil, 0, err
	}

	// Recalculate residual float based on dimensions of required vs shrunk images
	residualx := float64(o.Width) / float64(image.Xsize)
	residualy := float64(o.Height) / float64(image.Ysize)

	if o.Crop {
		residual = math.Max(residualx, residualy)
	} else {
		residual = math.Min(residualx, residualy)
	}

	return image, residual, nil
}

func shrinkJpegImage(buf []byte, input *C.struct__VipsImage, factor float64, shrink int) (*C.struct__VipsImage, float64, error) {
	var image *C.struct__VipsImage
	var err error
	shrinkOnLoad := 1

	// Recalculate integral shrink and double residual
	switch {
	case shrink >= 8:
		factor = factor / 8
		shrinkOnLoad = 8
	case shrink >= 4:
		factor = factor / 4
		shrinkOnLoad = 4
	case shrink >= 2:
		factor = factor / 2
		shrinkOnLoad = 2
	}

	// Reload input using shrink-on-load
	if shrinkOnLoad > 1 {
		image, err = vipsShrinkJpeg(buf, input, shrinkOnLoad)
	}

	return image, factor, err
}

func imageCalculations(o *Options, inWidth, inHeight int) float64 {
	factor := 1.0
	xfactor := float64(inWidth) / float64(o.Width)
	yfactor := float64(inHeight) / float64(o.Height)
	defer debug("Image calculations: %dx%d", o.Width, o.Height)

	switch {
	// Fixed width and height
	case o.Width > 0 && o.Height > 0:
		if o.Crop {
			factor = math.Min(xfactor, yfactor)
		} else {
			factor = math.Max(xfactor, yfactor)
		}
	// Fixed width, auto height
	case o.Width > 0:
		factor = xfactor
		o.Height = int(math.Floor(float64(inHeight) / factor))
	// Fixed height, auto width
	case o.Height > 0:
		factor = yfactor
		o.Width = int(math.Floor(float64(inWidth) / factor))
	default:
		// Identity transform
		o.Width = inWidth
		o.Height = inHeight
		break
	}

	return factor
}

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

func calculateRotationAndFlip(image *C.struct__VipsImage, angle Angle) (Angle, bool) {
	rotate := D0
	flip := false

	if angle == -1 {
		switch vipsExifOrientation(image) {
		case 6:
			rotate = D90
			break
		case 3:
			rotate = D180
			break
		case 8:
			rotate = D270
			break
		case 2:
			flip = true
			break // flip 1
		case 7:
			flip = true
			rotate = D90
			break // flip 6
		case 4:
			flip = true
			rotate = D180
			break // flip 3
		case 5:
			flip = true
			rotate = D270
			break // flip 8
		}
	} else {
		if angle == 90 {
			rotate = D90
		} else if angle == 180 {
			rotate = D180
		} else if angle == 270 {
			rotate = D270
		}
	}

	return rotate, flip
}

func getAngle(angle Angle) Angle {
	divisor := angle % 90
	if divisor != 0 {
		angle = angle - divisor
	}
	return Angle(math.Min(float64(angle), 270))
}
