package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

type ImageSize struct {
	Width  int
	Height int
}

type ImageMetadata struct {
	Orientation int
	Channels    int
	Alpha       bool
	Profile     bool
	Type        string
	Space       string
	Size        ImageSize
}

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
		Type:        getImageTypeName(imageType),
	}

	return metadata, nil
}
