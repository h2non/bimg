package bimg

/*
#cgo pkg-config: vips
#include "vips.h"
*/
import "C"

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

const (
	maxCacheMem  = 100 * 1024 * 1024
	maxCacheSize = 500
)

var (
	m           sync.Mutex
	initialized bool
)

type VipsMemoryInfo struct {
	Memory          int64
	MemoryHighwater int64
	Allocations     int64
}

type vipsSaveOptions struct {
	Quality     int
	Compression int
	Type        ImageType
}

type vipsWatermarkOptions struct {
	Width       C.int
	DPI         C.int
	Margin      C.int
	NoReplicate C.int
	Opacity     C.float
	Background  [3]C.double
}

type vipsWatermarkTextOptions struct {
	Text *C.char
	Font *C.char
}

func init() {
	Initialize()
}

// Explicit thread-safe start of libvips.
// Only call this function if you've previously shutdown libvips
func Initialize() {
	if C.VIPS_MAJOR_VERSION <= 7 && C.VIPS_MINOR_VERSION < 40 {
		panic("unsupported libvips version!")
	}

	m.Lock()
	runtime.LockOSThread()
	defer m.Unlock()
	defer runtime.UnlockOSThread()

	err := C.vips_init(C.CString("bimg"))
	if err != 0 {
		panic("unable to start vips!")
	}

	// Set libvips cache params
	C.vips_cache_set_max_mem(maxCacheMem)
	C.vips_cache_set_max(maxCacheSize)

	// Explicit concurrency limit to avoid thread-unsafe issues.
	// See: https://github.com/jcupitt/libvips/issues/261#issuecomment-92850414
	if os.Getenv("VIPS_CONCURRENCY") == "" {
		C.vips_concurrency_set(1)
	}

	if os.Getenv("VIPS_TRACE") != "" {
		C.vips_enable_cache_set_trace()
	}

	initialized = true
}

// Thread-safe function to shutdown libvips. You could call this to drop caches as well.
// If libvips was already initialized, the function is no-op
func Shutdown() {
	m.Lock()
	defer m.Unlock()

	if initialized {
		C.vips_shutdown()
		initialized = false
	}
}

// Output to stdout vips collected data. Useful for debugging
func VipsDebugInfo() {
	C.im__print_all()
}

// Get memory info stats from vips (cache size, memory allocs...)
func VipsMemory() VipsMemoryInfo {
	return VipsMemoryInfo{
		Memory:          int64(C.vips_tracked_get_mem()),
		MemoryHighwater: int64(C.vips_tracked_get_mem_highwater()),
		Allocations:     int64(C.vips_tracked_get_allocs()),
	}
}

func vipsExifOrientation(image *C.struct__VipsImage) int {
	return int(C.vips_exif_orientation(image))
}

func vipsHasAlpha(image *C.struct__VipsImage) bool {
	return int(C.has_alpha_channel(image)) > 0
}

func vipsHasProfile(image *C.struct__VipsImage) bool {
	return int(C.has_profile_embed(image)) > 0
}

func vipsWindowSize(name string) float64 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return float64(C.interpolator_window_size(cname))
}

func vipsSpace(image *C.struct__VipsImage) string {
	return C.GoString(C.vips_enum_nick_bridge(image))
}

func vipsRotate(image *C.struct__VipsImage, angle Angle) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(image))

	err := C.vips_rotate(image, &out, C.int(angle))
	if err != 0 {
		return nil, catchVipsError()
	}

	return out, nil
}

func vipsFlip(image *C.struct__VipsImage, direction Direction) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(image))

	err := C.vips_flip_bridge(image, &out, C.int(direction))
	if err != 0 {
		return nil, catchVipsError()
	}

	return out, nil
}

func vipsZoom(image *C.struct__VipsImage, zoom int) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(image))

	err := C.vips_zoom_bridge(image, &out, C.int(zoom), C.int(zoom))
	if err != 0 {
		return nil, catchVipsError()
	}

	return out, nil
}

func vipsWatermark(image *C.struct__VipsImage, w Watermark) (*C.struct__VipsImage, error) {
	var out *C.struct__VipsImage

	// Defaults
	noReplicate := 0
	if w.NoReplicate {
		noReplicate = 1
	}

	text := C.CString(w.Text)
	font := C.CString(w.Font)
	background := [3]C.double{C.double(w.Background.R), C.double(w.Background.G), C.double(w.Background.B)}

	textOpts := vipsWatermarkTextOptions{text, font}
	opts := vipsWatermarkOptions{C.int(w.Width), C.int(w.DPI), C.int(w.Margin), C.int(noReplicate), C.float(w.Opacity), background}

	defer C.free(unsafe.Pointer(text))
	defer C.free(unsafe.Pointer(font))

	err := C.vips_watermark(image, &out, (*C.WatermarkTextOptions)(unsafe.Pointer(&textOpts)), (*C.WatermarkOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, catchVipsError()
	}

	return out, nil
}

func vipsRead(buf []byte) (*C.struct__VipsImage, ImageType, error) {
	var image *C.struct__VipsImage
	imageType := vipsImageType(buf)

	if imageType == UNKNOWN {
		return nil, UNKNOWN, errors.New("Unsupported image format")
	}

	length := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	err := C.vips_init_image(imageBuf, length, C.int(imageType), &image)
	if err != 0 {
		return nil, UNKNOWN, catchVipsError()
	}

	return image, imageType, nil
}

func vipsSave(image *C.struct__VipsImage, o vipsSaveOptions) ([]byte, error) {
	var ptr unsafe.Pointer
	length := C.size_t(0)
	err := C.int(0)

	defer C.g_object_unref(C.gpointer(image))

	switch {
	case o.Type == PNG:
		err = C.vips_pngsave_bridge(image, &ptr, &length, 1, C.int(o.Compression), C.int(o.Quality), 0)
		break
	case o.Type == WEBP:
		err = C.vips_webpsave_bridge(image, &ptr, &length, 1, C.int(o.Quality), 0)
		break
	default:
		debug("Save JPEG options: Q: %s", o.Quality)
		err = C.vips_jpegsave_bridge(image, &ptr, &length, 1, C.int(o.Quality), 0)
		break
	}

	if int(err) != 0 {
		return nil, catchVipsError()
	}

	buf := C.GoBytes(ptr, C.int(length))

	// Clean up
	C.g_free(C.gpointer(ptr))
	C.vips_error_clear()

	return buf, nil
}

func vipsExtract(image *C.struct__VipsImage, left, top, width, height int) (*C.struct__VipsImage, error) {
	var buf *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(image))

	if width > MAX_SIZE || height > MAX_SIZE {
		return nil, errors.New("Maximum image size exceeded")
	}

	err := C.vips_extract_area_bridge(image, &buf, C.int(left), C.int(top), C.int(width), C.int(height))
	if err != 0 {
		return nil, catchVipsError()
	}

	return buf, nil
}

func vipsShrinkJpeg(buf []byte, input *C.struct__VipsImage, shrink int) (*C.struct__VipsImage, error) {
	var image *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(input))

	err := C.vips_jpegload_buffer_shrink(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), &image, C.int(shrink))
	if err != 0 {
		return nil, catchVipsError()
	}

	return image, nil
}

func vipsShrink(input *C.struct__VipsImage, shrink int) (*C.struct__VipsImage, error) {
	var image *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(input))

	err := C.vips_shrink_bridge(input, &image, C.double(float64(shrink)), C.double(float64(shrink)))
	if err != 0 {
		return nil, catchVipsError()
	}

	return image, nil
}

func vipsEmbed(input *C.struct__VipsImage, left, top, width, height, extend int) (*C.struct__VipsImage, error) {
	var image *C.struct__VipsImage
	defer C.g_object_unref(C.gpointer(input))

	err := C.vips_embed_bridge(input, &image, C.int(left), C.int(top), C.int(width), C.int(height), C.int(extend))
	if err != 0 {
		return nil, catchVipsError()
	}

	return image, nil
}

func vipsAffine(input *C.struct__VipsImage, residual float64, i Interpolator) (*C.struct__VipsImage, error) {
	var image *C.struct__VipsImage
	istring := C.CString(i.String())
	interpolator := C.vips_interpolate_new(istring)

	defer C.free(unsafe.Pointer(istring))
	defer C.g_object_unref(C.gpointer(input))
	defer C.g_object_unref(C.gpointer(interpolator))

	err := C.vips_affine_interpolator(input, &image, C.double(residual), 0, 0, C.double(residual), interpolator)
	if err != 0 {
		return nil, catchVipsError()
	}

	return image, nil
}

func vipsImageType(buf []byte) ImageType {
	imageType := UNKNOWN

	if len(buf) == 0 {
		return imageType
	}

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

func catchVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}
