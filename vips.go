package bimg

/*
#cgo pkg-config: vips
#include "vips.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"strings"
	"unsafe"
)

type vipsImage C.struct__VipsImage

func init() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := C.vips_init(C.CString("bimg"))
	if err != 0 {
		C.vips_shutdown()
		panic("unable to start vips!")
	}

	C.vips_concurrency_set(1)                   // default
	C.vips_cache_set_max_mem(100 * 1024 * 1024) // 100 MB
	C.vips_cache_set_max(500)                   // 500 operations
}

type Vips struct {
	buf []byte
}

func vipsRotate(image *C.struct__VipsImage, angle Angle) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage

	err := C.vips_rotate(image, &out, C.int(angle))
	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, catchVipsError()
	}
	defer C.g_object_unref(C.gpointer(out))

	return out, nil
}

func vipsFlip(image *C.struct__VipsImage, direction Direction) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage

	err := C.vips_flip_seq(image, &out)
	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, catchVipsError()
	}
	defer C.g_object_unref(C.gpointer(out))

	return out, nil
}

func vipsRead(buf []byte) (*C.struct__VipsImage, ImageType, error) {
	var image *C.struct__VipsImage
	imageType := vipsImageType(buf)

	if imageType == UNKNOWN {
		return nil, UNKNOWN, errors.New("Input buffer contains unsupported image format")
	}

	// feed it
	length := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])
	imageTypeC := C.int(imageType)

	err := C.vips_init_image(imageBuf, length, imageTypeC, &image)
	if err != 0 {
		return nil, UNKNOWN, catchVipsError()
	}

	return image, imageType, nil
}

func vipsExtract(image *C.struct__VipsImage, left int, top int, width int, height int) (*C.struct__VipsImage, error) {
	var buf *C.struct__VipsImage

	err := C.vips_extract_area_0(image, &buf, C.int(left), C.int(top), C.int(width), C.int(height))
	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, catchVipsError()
	}

	return buf, nil
}

func vipsImageType(buf []byte) ImageType {
	imageType := UNKNOWN

	length := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])
	bufferType := C.GoString(C.vips_foreign_find_load_buffer(imageBuf, length))

	switch {
	case strings.HasSuffix(bufferType, "JpegBuffer"):
		imageType = JPEG
		break
	case strings.HasSuffix(bufferType, "PngBuffer"):
		imageType = PNG
		break
	case strings.HasSuffix(bufferType, "TiffBuffer"):
		imageType = TIFF
		break
	case strings.HasSuffix(bufferType, "WebpBuffer"):
		imageType = WEBP
		break
	case strings.HasSuffix(bufferType, "MagickBuffer"):
		imageType = MAGICK
		break
	}

	return imageType
}

func vipsExifOrientation(image *C.struct__VipsImage) int {
	return int(C.vips_exif_orientation(image))
}

type vipsSaveOptions struct {
	Quality     int
	Compression int
	Type        ImageType
}

func vipsSave(image *C.struct__VipsImage, o vipsSaveOptions) ([]byte, error) {
	var ptr unsafe.Pointer
	length := C.size_t(0)
	err := C.int(0)

	switch {
	case o.Type == PNG:
		err = C.vips_pngsave_custom(image, &ptr, &length, 1, C.int(o.Compression), C.int(o.Quality), 0)
		break
	case o.Type == WEBP:
		err = C.vips_webpsave_custom(image, &ptr, &length, 1, C.int(o.Quality), 0)
		break
	default:
		err = C.vips_jpegsave_custom(image, &ptr, &length, 1, C.int(o.Quality), 0)
		break
	}

	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, catchVipsError()
	}

	buf := C.GoBytes(ptr, C.int(length))
	// cleanup
	C.g_free(C.gpointer(ptr))

	return buf, nil
}

func catchVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}
