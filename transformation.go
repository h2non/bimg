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

type ResizeMode int

const (
	ResizeModeFit ResizeMode = iota
	ResizeModeUp
	ResizeModeDown
	ResizeModeForce
)

type ResizeOptions struct {
	Height         int
	Width          int
	Top            int
	Left           int
	Zoom           int
	Mode           ResizeMode
	Interpolator   Interpolator
	Interpretation Interpretation
}

func calculateResizeFactor(opts *ResizeOptions, inWidth, inHeight int) float64 {
	factor := 1.0
	xfactor := float64(inWidth) / float64(opts.Width)
	yfactor := float64(inHeight) / float64(opts.Height)

	switch {
	// Fixed width and height
	case opts.Width > 0 && opts.Height > 0:
		if opts.Mode == ResizeModeForce {
			factor = math.Max(xfactor, yfactor)
		} else {
			factor = math.Min(xfactor, yfactor)
		}
	// Fixed width, auto height
	case opts.Width > 0:
		if opts.Mode == ResizeModeForce {
			factor = xfactor
			opts.Height = roundFloat(float64(inHeight) / factor)
		} else {
			opts.Height = inHeight
		}
	// Fixed height, auto width
	case opts.Height > 0:
		if opts.Mode == ResizeModeForce {
			factor = yfactor
			opts.Width = roundFloat(float64(inWidth) / factor)
		} else {
			opts.Width = inWidth
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

	inWidth := int(it.image.c.Xsize)
	inHeight := int(it.image.c.Ysize)

	// image calculations
	factor := calculateResizeFactor(&opts, inWidth, inHeight)
	shrink := calculateShrink(factor, opts.Interpolator)
	residual := calculateResidual(factor, shrink)

	// Do not enlarge the output if the input width or height
	// are already less than the required dimensions
	if (opts.Mode == ResizeModeDown || opts.Mode == ResizeModeFit) &&
		inWidth < opts.Width && inHeight < opts.Height {

		factor = 1.0
		shrink = 1
		residual = 0
		opts.Width = inWidth
		opts.Height = inHeight
	}

	// Try to use libjpeg/libwebp shrink-on-load, if the buffer is still usable.
	// If we performed "destructive" transformations already, this will no longer
	// be the case.
	isShrinkableWebP := it.imageType == WEBP && vipsVersionMin(8, 3)
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
	if image, err := transformImage(it.image, opts, shrink, residual); err != nil {
		return err
	} else {
		it.updateImage(image)
	}

	return nil
}

type CropOptions struct {
	Width   int
	Height  int
	Gravity Gravity
}

func (it *ImageTransformation) Crop(opts CropOptions) error {
	inWidth := int(it.image.c.Xsize)
	inHeight := int(it.image.c.Ysize)

	// it's already at an appropriate size, return immediately
	if inWidth <= opts.Width && inHeight <= opts.Height {
		return nil
	}

	if opts.Gravity == GravitySmart {
		width := int(math.Min(float64(inWidth), float64(opts.Width)))
		height := int(math.Min(float64(inHeight), float64(opts.Height)))

		if image, err := vipsSmartCrop(it.image, width, height); err != nil {
			return err
		} else {
			it.updateImage(image)
			return nil
		}
	} else {
		width := int(math.Min(float64(inWidth), float64(opts.Width)))
		height := int(math.Min(float64(inHeight), float64(opts.Height)))
		left, top := calculateCrop(inWidth, inHeight, opts.Width, opts.Height, opts.Gravity)
		left, top = int(math.Max(float64(left), 0)), int(math.Max(float64(top), 0))

		if image, err := vipsExtract(it.image, left, top, width, height); err != nil {
			return err
		} else {
			it.updateImage(image)
			return nil
		}
	}
}

type TrimOptions struct {
	Background RGBAProvider
	Threshold  float64
}

func (it *ImageTransformation) Trim(opts TrimOptions) error {
	left, top, width, height, err := vipsTrim(it.image, opts.Background, opts.Threshold)
	if err != nil {
		return fmt.Errorf("cannot determine trim area: %w", err)
	}

	if image, err := vipsExtract(it.image, left, top, width, height); err != nil {
		return fmt.Errorf("cannot extract trim area: %w", err)
	} else {
		it.updateImage(image)
		return nil
	}
}

type EmbedOptions struct {
	Width      int
	Height     int
	Extend     Extend
	Background RGBAProvider
}

func (it *ImageTransformation) Embed(opts EmbedOptions) error {
	inWidth := int(it.image.c.Xsize)
	inHeight := int(it.image.c.Ysize)

	left, top := (opts.Width-inWidth)/2, (opts.Height-inHeight)/2
	if image, err := vipsEmbed(it.image, left, top, opts.Width, opts.Height, opts.Extend, opts.Background); err != nil {
		return err
	} else {
		it.updateImage(image)
		return err
	}
}

type ExtractOptions struct {
	Left   int
	Top    int
	Width  int
	Height int
}

func (it *ImageTransformation) Extract(opts ExtractOptions) error {
	if opts.Width == 0 || opts.Height == 0 {
		return errors.New("extract area width/height params are required")
	}
	if image, err := vipsExtract(it.image, opts.Left, opts.Top, opts.Width, opts.Height); err != nil {
		return err
	} else {
		it.updateImage(image)
		return nil
	}
}

type RotateOptions struct {
	Angle        Angle
	Flip         bool
	Flop         bool
	NoAutoRotate bool
}

func (it *ImageTransformation) Rotate(opts RotateOptions) error {
	var image *vipsImage
	var err error

	if opts.NoAutoRotate {
		image = it.image
	} else {
		image, err = vipsAutoRotate(it.image)
		if err != nil {
			return fmt.Errorf("cannot autorotate image: %w", err)
		}
	}

	image, err = rotateAndFlipImage(image, opts)
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

type WatermarkImageOptions struct {
	Left    int
	Top     int
	Image   *ImageTransformation
	Opacity float32
}

func (it *ImageTransformation) WatermarkImage(opts WatermarkImageOptions) error {
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

func (it *ImageTransformation) Size() ImageSize {
	return ImageSize{
		Width:  int(it.image.c.Xsize),
		Height: int(it.image.c.Ysize),
	}
}

func (it *ImageTransformation) Metadata() ImageMetadata {
	size := it.Size()

	orientation := vipsExifIntTag(it.image, Orientation)

	return ImageMetadata{
		Size:        size,
		Channels:    int(it.image.c.Bands),
		Orientation: orientation,
		Alpha:       vipsHasAlpha(it.image),
		Profile:     vipsHasProfile(it.image),
		Space:       vipsSpace(it.image),
		Type:        ImageTypeName(it.imageType),
		EXIF: EXIF{
			Make:                    vipsExifStringTag(it.image, Make),
			Model:                   vipsExifStringTag(it.image, Model),
			Orientation:             orientation,
			XResolution:             vipsExifStringTag(it.image, XResolution),
			YResolution:             vipsExifStringTag(it.image, YResolution),
			ResolutionUnit:          vipsExifIntTag(it.image, ResolutionUnit),
			Software:                vipsExifStringTag(it.image, Software),
			Datetime:                vipsExifStringTag(it.image, Datetime),
			YCbCrPositioning:        vipsExifIntTag(it.image, YCbCrPositioning),
			Compression:             vipsExifIntTag(it.image, Compression),
			ExposureTime:            vipsExifStringTag(it.image, ExposureTime),
			FNumber:                 vipsExifStringTag(it.image, FNumber),
			ExposureProgram:         vipsExifIntTag(it.image, ExposureProgram),
			ISOSpeedRatings:         vipsExifIntTag(it.image, ISOSpeedRatings),
			ExifVersion:             vipsExifStringTag(it.image, ExifVersion),
			DateTimeOriginal:        vipsExifStringTag(it.image, DateTimeOriginal),
			DateTimeDigitized:       vipsExifStringTag(it.image, DateTimeDigitized),
			ComponentsConfiguration: vipsExifStringTag(it.image, ComponentsConfiguration),
			ShutterSpeedValue:       vipsExifStringTag(it.image, ShutterSpeedValue),
			ApertureValue:           vipsExifStringTag(it.image, ApertureValue),
			BrightnessValue:         vipsExifStringTag(it.image, BrightnessValue),
			ExposureBiasValue:       vipsExifStringTag(it.image, ExposureBiasValue),
			MeteringMode:            vipsExifIntTag(it.image, MeteringMode),
			Flash:                   vipsExifIntTag(it.image, Flash),
			FocalLength:             vipsExifStringTag(it.image, FocalLength),
			SubjectArea:             vipsExifStringTag(it.image, SubjectArea),
			MakerNote:               vipsExifStringTag(it.image, MakerNote),
			SubSecTimeOriginal:      vipsExifStringTag(it.image, SubSecTimeOriginal),
			SubSecTimeDigitized:     vipsExifStringTag(it.image, SubSecTimeDigitized),
			ColorSpace:              vipsExifIntTag(it.image, ColorSpace),
			PixelXDimension:         vipsExifIntTag(it.image, PixelXDimension),
			PixelYDimension:         vipsExifIntTag(it.image, PixelYDimension),
			SensingMethod:           vipsExifIntTag(it.image, SensingMethod),
			SceneType:               vipsExifStringTag(it.image, SceneType),
			ExposureMode:            vipsExifIntTag(it.image, ExposureMode),
			WhiteBalance:            vipsExifIntTag(it.image, WhiteBalance),
			FocalLengthIn35mmFilm:   vipsExifIntTag(it.image, FocalLengthIn35mmFilm),
			SceneCaptureType:        vipsExifIntTag(it.image, SceneCaptureType),
			GPSLatitudeRef:          vipsExifStringTag(it.image, GPSLatitudeRef),
			GPSLatitude:             vipsExifStringTag(it.image, GPSLatitude),
			GPSLongitudeRef:         vipsExifStringTag(it.image, GPSLongitudeRef),
			GPSLongitude:            vipsExifStringTag(it.image, GPSLongitude),
			GPSAltitudeRef:          vipsExifStringTag(it.image, GPSAltitudeRef),
			GPSAltitude:             vipsExifStringTag(it.image, GPSAltitude),
			GPSSpeedRef:             vipsExifStringTag(it.image, GPSSpeedRef),
			GPSSpeed:                vipsExifStringTag(it.image, GPSSpeed),
			GPSImgDirectionRef:      vipsExifStringTag(it.image, GPSImgDirectionRef),
			GPSImgDirection:         vipsExifStringTag(it.image, GPSImgDirection),
			GPSDestBearingRef:       vipsExifStringTag(it.image, GPSDestBearingRef),
			GPSDestBearing:          vipsExifStringTag(it.image, GPSDestBearing),
			GPSDateStamp:            vipsExifStringTag(it.image, GPSDateStamp),
		},
	}
}

// TODO convert o.SmartCrop to o.Gravity
