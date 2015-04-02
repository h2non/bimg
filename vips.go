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

type vipsImage *C.struct__VipsImage

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

func vipsRotate(image *C.struct__VipsImage, degrees int) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage

	err := C.vips_rotate(image, &out, C.int(degrees))
	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, vipsError()
	}
	defer C.g_object_unref(C.gpointer(out))

	return out, nil
}

func vipsFlip(image *C.struct__VipsImage, direction Direction) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage

	err := C.vips_flip_seq(image, &out)
	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, vipsError()
	}
	defer C.g_object_unref(C.gpointer(out))

	return out, nil
}

func vipsRead(buf []byte) (*C.struct__VipsImage, error) {
	var image *C.struct__VipsImage

	debug("Format: %s", vipsImageType(buf))

	// feed it
	length := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	err := C.vips_jpegload_buffer_seq(imageBuf, length, &image)
	if err != 0 {
		return nil, vipsError()
	}

	return image, nil
}

func vipsExtract(image *C.struct__VipsImage, left int, top int, width int, height int) (*C.struct__VipsImage, error) {
	var buf *C.struct__VipsImage

	err := C.vips_extract_area_0(image, &buf, C.int(left), C.int(top), C.int(width), C.int(height))
	C.g_object_unref(C.gpointer(image))
	if err != 0 {
		return nil, vipsError()
	}

	return buf, nil
}

func vipsImageType(buf []byte) string {
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

type vipsSaveOptions struct {
	Quality int
}

func vipsSave(image *C.struct__VipsImage, o vipsSaveOptions) ([]byte, error) {
	var ptr unsafe.Pointer
	length := C.size_t(0)

	err := C.vips_jpegsave_custom(image, &ptr, &length, 1, C.int(o.Quality), 0)
	if err != 0 {
		return nil, vipsError()
	}
	C.g_object_unref(C.gpointer(image))

	buf := C.GoBytes(ptr, C.int(length))
	// cleanup
	C.g_free(C.gpointer(ptr))

	return buf, nil
}

func vipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}
