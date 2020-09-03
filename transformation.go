package bimg

import "C"
import "runtime"

type ImageTransformation struct {
	image     *C.VipsImage
	imageType ImageType
}

func NewImageTransformation(buf []byte) (*ImageTransformation, error) {
	image, imageType, err := vipsRead(buf)
	if err != nil {
		return nil, err
	}
	it := &ImageTransformation{
		image:     image,
		imageType: imageType,
	}
	runtime.SetFinalizer(it, finalizeImageTransformation)
	return it, nil
}

func (it *ImageTransformation) Clone() *ImageTransformation {
	clone := &ImageTransformation{
		image:     it.image,
		imageType: it.imageType,
	}
	C.g_object_ref(C.gpointer(clone.image))
	runtime.SetFinalizer(it, finalizeImageTransformation)
	return clone
}

func (it *ImageTransformation) Close() {
	C.g_object_unref(C.gpointer(it.image))
	it.image = nil
}

func finalizeImageTransformation(it *ImageTransformation) {
	it.Close()
}

func (it *ImageTransformation) updateImage(image *C.VipsImage) {
	C.g_object_unref(C.gpointer(it.image))
	it.image = image
}
