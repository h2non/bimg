package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"
import (
	"fmt"
	"math"
)

type ImageTransformation struct {
	buf        []byte
	bufTainted bool
	image      *vipsImage
	imageType  ImageType
}

func NewImageTransformation(buf []byte) (*ImageTransformation, error) {
	image, imageType, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}
	it := &ImageTransformation{
		buf:        buf,
		bufTainted: false,
		image:      image,
		imageType:  imageType,
	}
	return it, nil
}

func (it *ImageTransformation) Clone() *ImageTransformation {
	return &ImageTransformation{
		buf:        it.buf,
		bufTainted: it.bufTainted,
		image:      it.image.clone(),
		imageType:  it.imageType,
	}
}

func (it *ImageTransformation) Close() {
	it.image.close()
	it.image = nil
	it.buf = nil
}

func (it *ImageTransformation) updateImage(image *vipsImage) {
	if it.image == image {
		return
	}

	if it.image != nil {
		it.image.close()
	}
	it.image = image
	// We replaced the image, so the buffer is no longer the same content.
	it.bufTainted = true
}

type ResizeOptions struct {
	Height         int
	Width          int
	AreaHeight     int
	AreaWidth      int
	Top            int
	Left           int
	Zoom           int
	Crop           bool
	Enlarge        bool
	Embed          bool
	Force          bool
	Trim           bool
	Extend         Extend
	Background     RGBAProvider
	Gravity        Gravity
	Interpolator   Interpolator
	Interpretation Interpretation
	Threshold      float64
}

func calculateResizeFactor(opts *ResizeOptions, inWidth, inHeight int) float64 {
	factor := 1.0
	xfactor := float64(inWidth) / float64(opts.Width)
	yfactor := float64(inHeight) / float64(opts.Height)

	switch {
	// Fixed width and height
	case opts.Width > 0 && opts.Height > 0:
		if opts.Crop {
			factor = math.Min(xfactor, yfactor)
		} else {
			factor = math.Max(xfactor, yfactor)
		}
	// Fixed width, auto height
	case opts.Width > 0:
		if opts.Crop {
			opts.Height = inHeight
		} else {
			factor = xfactor
			opts.Height = roundFloat(float64(inHeight) / factor)
		}
	// Fixed height, auto width
	case opts.Height > 0:
		if opts.Crop {
			opts.Width = inWidth
		} else {
			factor = yfactor
			opts.Width = roundFloat(float64(inWidth) / factor)
		}
	// Identity transform
	default:
		opts.Width = inWidth
		opts.Height = inHeight
		break
	}

	return factor
}

func (it *ImageTransformation) Resize(opts ResizeOptions) error {
	if opts.Interpretation == 0 {
		opts.Interpretation = InterpretationSRGB
	}
	if !opts.Force && !opts.Crop && !opts.Embed && !opts.Enlarge && (opts.Width > 0 || opts.Height > 0) {
		opts.Force = true
	}

	inWidth := int(it.image.c.Xsize)
	inHeight := int(it.image.c.Ysize)

	// image calculations
	factor := calculateResizeFactor(&opts, inWidth, inHeight)
	shrink := calculateShrink(factor, opts.Interpolator)
	residual := calculateResidual(factor, shrink)

	// Do not enlarge the output if the input width or height
	// are already less than the required dimensions
	if !opts.Enlarge && !opts.Force {
		if inWidth < opts.Width && inHeight < opts.Height {
			factor = 1.0
			shrink = 1
			residual = 0
			opts.Width = inWidth
			opts.Height = inHeight
		}
	}

	// Try to use libjpeg/libwebp shrink-on-load, if the buffer is still usable.
	// If we performed "destructive" transformations already, this will no longer
	// be the case.
	isShrinkableWebP := it.imageType == WEBP && VipsMajorVersion >= 8 && VipsMinorVersion >= 3
	isShrinkableJpeg := it.imageType == JPEG
	supportsShrinkOnLoad := !it.bufTainted && (isShrinkableWebP || isShrinkableJpeg)

	if supportsShrinkOnLoad && shrink >= 2 {
		tmpImage, factor, err := shrinkOnLoad(it.buf, it.imageType, factor, shrink)
		if err != nil {
			return fmt.Errorf("cannot shrink-on-load: %w", err)
		}

		it.updateImage(tmpImage)
		factor = math.Max(factor, 1.0)
		shrink = int(math.Floor(factor))
		residual = float64(shrink) / factor
	}

	// Zoom image, if necessary
	if image, err := zoomImage(it.image, opts.Zoom); err != nil {
		return fmt.Errorf("cannot zoom image: %w", err)
	} else {
		it.updateImage(image)
	}

	// Transform image, if necessary
	if shouldTransformImage(opts, inWidth, inHeight) {
		if image, err := transformImage(it.image, opts, shrink, residual); err != nil {
			return err
		} else {
			it.updateImage(image)
		}
	}

	return nil
}

type RotateOptions struct {
	Angle        Angle
	Flip         bool
	Flop         bool
	NoAutoRotate bool
}

func (it *ImageTransformation) Rotate(opts RotateOptions) error {
	image, _, err := rotateAndFlipImage(it.image, opts)
	if err != nil {
		return err
	}
	it.updateImage(image)

	// TODO is the following a good idea? Isn't saving always destructive and so time consuming
	//      that it outweighs the potential time saving when using the optimized shrink?
	// If it's a JPEG or HEIF image, the rotation might have been non-destructive. So we
	// update our buffer (which involves saving).
	//if rotated && (it.imageType == JPEG || it.imageType == HEIF) {
	//	buf, err := getImageBuffer(image)
	//	if err != nil {
	//		return fmt.Errorf("cannot get the rotated image buffer: %w", err)
	//	}
	//}

	return nil
}

func (it *ImageTransformation) Blur(opts GaussianBlur) error {
	if image, err := vipsGaussianBlur(it.image, opts); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

func (it *ImageTransformation) Sharpen(opts Sharpen) error {
	if image, err := vipsSharpen(it.image, opts); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

func (it *ImageTransformation) WatermarkText(opts Watermark) error {
	if image, err := watermarkImageWithText(it.image, opts); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

func (it *ImageTransformation) WatermarkImage(opts WatermarkImage) error {
	if image, err := watermarkImageWithAnotherImage(it.image, opts); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

func (it *ImageTransformation) Flatten(background RGBAProvider) error {
	if image, err := vipsFlattenBackground(it.image, background); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

func (it *ImageTransformation) Gamma(gamma float64) error {
	if image, err := vipsGamma(it.image, gamma); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

type SaveOptions vipsSaveOptions

func (it *ImageTransformation) Save(opts SaveOptions) ([]byte, error) {
	// Normalize options first.
	if opts.Quality == 0 {
		opts.Quality = Quality
	}
	if opts.Compression == 0 {
		opts.Compression = 6
	}
	if opts.Type == 0 {
		opts.Type = it.imageType
	}

	return vipsSave(it.image, vipsSaveOptions(opts))
}

// TODO convert o.SmartCrop to o.Gravity
