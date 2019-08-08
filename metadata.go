package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

// ImageSize represents the image width and height values
type ImageSize struct {
	Width  int
	Height int
}

// ImageMetadata represents the basic metadata fields
type ImageMetadata struct {
	Orientation int
	Channels    int
	Alpha       bool
	Profile     bool
	Type        string
	Space       string
	Colourspace string
	Size        ImageSize
}

// Size returns the image size by width and height pixels.
func Size(buf []byte) (ImageSize, error) {
	metadata, err := Metadata(buf)
	if err != nil {
		return ImageSize{}, err
	}

	return ImageSize{
		Width:  int(metadata.Size.Width),
		Height: int(metadata.Size.Height),
	}, nil
}

// ColourspaceIsSupported checks if the image colourspace is supported by libvips.
func ColourspaceIsSupported(buf []byte) (bool, error) {
	return vipsColourspaceIsSupportedBuffer(buf)
}

// ImageInterpretation returns the image interpretation type.
// See: https://jcupitt.github.io/libvips/API/current/VipsImage.html#VipsInterpretation
func ImageInterpretation(buf []byte) (Interpretation, error) {
	return vipsInterpretationBuffer(buf)
}

// Metadata returns the image metadata (size, type, alpha channel, profile, EXIF orientation...).
func Metadata(buf []byte) (ImageMetadata, error) {
	defer C.vips_thread_shutdown()

	image, imageType, err := vipsRead(buf)
	if err != nil {
		return ImageMetadata{}, err
	}
	defer C.g_object_unref(C.gpointer(image))

	size := ImageSize{
		Width:  int(image.Xsize),
		Height: int(image.Ysize),
	}

	metadata := ImageMetadata{
		Size:        size,
		Channels:    int(image.Bands),
		Orientation: vipsExifOrientation(image),
		Alpha:       vipsHasAlpha(image),
		Profile:     vipsHasProfile(image),
		Space:       vipsSpace(image),
		Type:        ImageTypeName(imageType),
	}

	return metadata, nil
}

// Npages returns the number of the given image.
func Npages(buf []byte) (int, error) {
	image, _, err := loadImage(buf)
	if err != nil {
		return -1, err
	}
	defer C.g_object_unref(C.gpointer(image))

	return int(C.vips_image_get_n_pages(image)), nil
}

// GetPageHeight returns height of the given image.
func GetPageHeight(buf []byte) (int, error) {
	image, _, err := loadImage(buf)
	if err != nil {
		return -1, err
	}
	defer C.g_object_unref(C.gpointer(image))

	return int(C.vips_image_get_page_height(image)), nil
}

// SetImageHeight returns ...
func SetImageHeight(buf []byte, h int) []byte {
	defer C.vips_thread_shutdown()

	img, imageType, err := loadImage(buf)
	if err != nil {
		return nil
	}

	o := applyDefaults(Options{}, imageType)

	buffer, err := vipsSetImageHeight(img, h, o)
	if err != nil {
		return nil
	}

	return buffer
}
