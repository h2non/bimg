package bimg

/*
#cgo pkg-config: vips
#include "vips.h"
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

// VipsVersion exposes the current libvips semantic version
const VipsVersion = string(C.VIPS_VERSION)

// VipsMajorVersion exposes the current libvips major version number
const VipsMajorVersion = int(C.VIPS_MAJOR_VERSION)

// VipsMinorVersion exposes the current libvips minor version number
const VipsMinorVersion = int(C.VIPS_MINOR_VERSION)

const (
	maxCacheMem  = 100 * 1024 * 1024
	maxCacheSize = 500
)

var (
	m           sync.Mutex
	initialized bool
)

// VipsMemoryInfo represents the memory stats provided by libvips.
type VipsMemoryInfo struct {
	Memory          int64
	MemoryHighwater int64
	Allocations     int64
}

type vipsImage struct {
	c *C.VipsImage
}

// vipsSaveOptions represents the internal option used to talk with libvips.
type vipsSaveOptions struct {
	Quality        int
	Compression    int
	Speed          int // Speed defines the AVIF encoders CPU effort. Valid values are 0-8.
	Type           ImageType
	MagickFormat   string // Format to use when saving using the ImageType MAGICK
	Interlace      bool
	NoProfile      bool
	StripMetadata  bool
	Lossless       bool
	InputICC       string // Absolute path to the input ICC profile
	OutputICC      string // Absolute path to the output ICC profile
	Interpretation Interpretation
	Palette        bool
}

type vipsWatermarkOptions struct {
	Width       C.int
	DPI         C.int
	Margin      C.int
	NoReplicate C.int
	Opacity     C.float
	Background  [3]C.double
}

type vipsWatermarkImageOptions struct {
	Left    C.int
	Top     C.int
	Opacity C.float
}

type vipsWatermarkTextOptions struct {
	Text *C.char
	Font *C.char
}

func init() {
	Initialize()
}

func wrapVipsImage(cImage *C.VipsImage) *vipsImage {
	vipsImage := &vipsImage{cImage}
	runtime.SetFinalizer(vipsImage, finalizeVipsImage)
	return vipsImage
}

func (vi *vipsImage) isNil() bool {
	return vi.c == nil
}

func (vi *vipsImage) clone() *vipsImage {
	C.g_object_ref(C.gpointer(vi.c))
	return wrapVipsImage(vi.c)
}

func (vi *vipsImage) close() {
	if vi.c != nil {
		C.g_object_unref(C.gpointer(vi.c))
		vi.c = nil
	}
}

func finalizeVipsImage(vi *vipsImage) {
	vi.close()
}

func vipsVersionMin(major, minor int) bool {
	return VipsMajorVersion > major || (VipsMajorVersion == major && VipsMinorVersion >= minor)
}

// Initialize is used to explicitly start libvips in thread-safe way.
// Only call this function if you have previously turned off libvips.
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

	// Define a custom thread concurrency limit in libvips (this may generate thread-unsafe issues)
	// See: https://github.com/jcupitt/libvips/issues/261#issuecomment-92850414
	if os.Getenv("VIPS_CONCURRENCY") == "" {
		C.vips_concurrency_set(1)
	}

	// Enable libvips cache tracing
	if os.Getenv("VIPS_TRACE") != "" {
		C.vips_enable_cache_set_trace()
	}

	initialized = true
}

// Shutdown is used to shutdown libvips in a thread-safe way.
// You can call this to drop caches as well.
// If libvips was already initialized, the function is no-op
func Shutdown() {
	m.Lock()
	defer m.Unlock()

	if initialized {
		C.vips_shutdown()
		initialized = false
	}
}

// VipsCacheSetMaxMem Sets the maximum amount of tracked memory allowed before the vips operation cache
// begins to drop entries.
func VipsCacheSetMaxMem(maxCacheMem int) {
	C.vips_cache_set_max_mem(C.size_t(maxCacheMem))
}

// VipsCacheSetMax sets the maximum number of operations to keep in the vips operation cache.
func VipsCacheSetMax(maxCacheSize int) {
	C.vips_cache_set_max(C.int(maxCacheSize))
}

// VipsCacheDropAll drops the vips operation cache, freeing the allocated memory.
func VipsCacheDropAll() {
	C.vips_cache_drop_all()
}

// VipsVectorSetEnabled enables or disables SIMD vector instructions. This can give speed-up,
// but can also be unstable on some systems and versions.
func VipsVectorSetEnabled(enable bool) {
	flag := 0
	if enable {
		flag = 1
	}

	C.vips_vector_set_enabled(C.int(flag))
}

// VipsDebugInfo outputs to stdout libvips collected data. Useful for debugging.
func VipsDebugInfo() {
	C.im__print_all()
}

// VipsMemory gets memory info stats from libvips (cache size, memory allocs...)
func VipsMemory() VipsMemoryInfo {
	return VipsMemoryInfo{
		Memory:          int64(C.vips_tracked_get_mem()),
		MemoryHighwater: int64(C.vips_tracked_get_mem_highwater()),
		Allocations:     int64(C.vips_tracked_get_allocs()),
	}
}

// VipsIsTypeSupported returns true if the given image type
// is supported by the current libvips compilation.
func VipsIsTypeSupported(t ImageType) bool {
	switch t {
	case JPEG:
		return int(C.vips_type_find_bridge(C.JPEG)) != 0
	case WEBP:
		return int(C.vips_type_find_bridge(C.WEBP)) != 0
	case PNG:
		return int(C.vips_type_find_bridge(C.PNG)) != 0
	case GIF:
		return int(C.vips_type_find_bridge(C.GIF)) != 0
	case PDF:
		return int(C.vips_type_find_bridge(C.PDF)) != 0
	case SVG:
		return int(C.vips_type_find_bridge(C.SVG)) != 0
	case TIFF:
		return int(C.vips_type_find_bridge(C.TIFF)) != 0
	case MAGICK:
		return int(C.vips_type_find_bridge(C.MAGICK)) != 0
	case HEIF, AVIF:
		return int(C.vips_type_find_bridge(C.HEIF)) != 0
	case JP2K:
		return int(C.vips_type_find_bridge(C.JP2K)) != 0
	default:
		return false
	}
}

// VipsIsTypeSupportedSave returns true if the given image type
// is supported by the current libvips compilation for the
// save operation.
func VipsIsTypeSupportedSave(t ImageType) bool {
	switch t {
	case JPEG:
		return int(C.vips_type_find_save_bridge(C.JPEG)) != 0
	case WEBP:
		return int(C.vips_type_find_save_bridge(C.WEBP)) != 0
	case PNG:
		return int(C.vips_type_find_save_bridge(C.PNG)) != 0
	case TIFF:
		return int(C.vips_type_find_save_bridge(C.TIFF)) != 0
	case HEIF, AVIF:
		return int(C.vips_type_find_save_bridge(C.HEIF)) != 0
	case GIF:
		return int(C.vips_type_find_save_bridge(C.MAGICK)) != 0
	case MAGICK:
		return int(C.vips_type_find_save_bridge(C.MAGICK)) != 0
	case JP2K:
		return int(C.vips_type_find_save_bridge(C.JP2K)) != 0
	default:
		return false
	}
}

func vipsExifStringTag(image *vipsImage, tag string) string {
	return vipsExifShort(C.GoString(C.vips_exif_tag(image.c, C.CString(tag))))
}

func vipsExifIntTag(image *vipsImage, tag string) int {
	return int(C.vips_exif_tag_to_int(image.c, C.CString(tag)))
}

func vipsExifOrientation(image *vipsImage) int {
	return int(C.vips_exif_orientation(image.c))
}

func vipsExifShort(s string) string {
	i := strings.Index(s, " (")
	if i > 0 {
		return s[:i]
	}
	return s
}

func vipsHasAlpha(image *vipsImage) bool {
	return int(C.has_alpha_channel(image.c)) > 0
}

func vipsHasProfile(image *vipsImage) bool {
	return int(C.has_profile_embed(image.c)) > 0
}

func vipsWindowSize(name string) float64 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return float64(C.interpolator_window_size(cname))
}

func vipsSpace(image *vipsImage) string {
	return C.GoString(C.vips_enum_nick_bridge(image.c))
}

func vipsRotate(image *vipsImage, angle Angle) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_rotate_bridge(image.c, &out, C.int(angle))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsAutoRotate(image *vipsImage) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_autorot_bridge(image.c, &out)
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsFlip(image *vipsImage, direction Direction) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_flip_bridge(image.c, &out, C.int(direction))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsZoom(image *vipsImage, zoom int) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_zoom_bridge(image.c, &out, C.int(zoom), C.int(zoom))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsWatermark(image *vipsImage, w Watermark) (*vipsImage, error) {
	var out *C.VipsImage

	// Defaults
	noReplicate := 0
	if w.NoReplicate {
		noReplicate = 1
	}

	text := C.CString(w.Text)
	font := C.CString(w.Font)

	var r, g, b uint8
	if w.Background != nil {
		r, g, b, _ = w.Background.RGBA()
	}

	background := [3]C.double{C.double(r), C.double(g), C.double(b)}

	textOpts := vipsWatermarkTextOptions{text, font}
	opts := vipsWatermarkOptions{C.int(w.Width), C.int(w.DPI), C.int(w.Margin), C.int(noReplicate), C.float(w.Opacity), background}

	defer C.free(unsafe.Pointer(text))
	defer C.free(unsafe.Pointer(font))

	err := C.vips_watermark(image.c, &out, (*C.WatermarkTextOptions)(unsafe.Pointer(&textOpts)), (*C.WatermarkOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsRead(buf []byte) (*vipsImage, ImageType, error) {
	var image *C.VipsImage
	imageType := vipsImageType(buf)

	if imageType == UNKNOWN {
		return nil, UNKNOWN, errors.New("unsupported image format")
	}

	length := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	err := C.vips_init_image(imageBuf, length, C.int(imageType), &image)
	if err != 0 {
		return nil, UNKNOWN, catchVipsError()
	}

	return wrapVipsImage(image), imageType, nil
}

func vipsColourspaceIsSupportedBuffer(buf []byte) (bool, error) {
	image, _, err := vipsRead(buf)
	if err != nil {
		return false, err
	}
	return vipsColourspaceIsSupported(image), nil
}

func vipsColourspaceIsSupported(image *vipsImage) bool {
	return int(C.vips_colourspace_issupported_bridge(image.c)) == 1
}

func vipsInterpretationBuffer(buf []byte) (Interpretation, error) {
	image, _, err := vipsRead(buf)
	if err != nil {
		return InterpretationError, err
	}
	return vipsInterpretation(image), nil
}

func vipsInterpretation(image *vipsImage) Interpretation {
	return Interpretation(C.vips_image_guess_interpretation_bridge(image.c))
}

func vipsFlattenBackground(image *vipsImage, background RGBAProvider) (*vipsImage, error) {
	if background == nil {
		return nil, errors.New("cannot flatten without background")
	}

	if !vipsHasAlpha(image) {
		return image, nil
	}

	if newImage, err := vipsAdjustColourspaceToColour(image, background); err != nil {
		return nil, err
	} else {
		image = newImage
	}

	r, g, b, _ := background.RGBA()
	backgroundC := [3]C.double{
		C.double(r),
		C.double(g),
		C.double(b),
	}

	var outImage *C.VipsImage
	err := C.vips_flatten_background_bridge(image.c, &outImage,
		backgroundC[0], backgroundC[1], backgroundC[2])
	if int(err) != 0 {
		return nil, catchVipsError()
	}
	return wrapVipsImage(outImage), nil
}

func vipsPreSave(image *vipsImage, o *vipsSaveOptions) (*vipsImage, error) {
	var outImage *C.VipsImage
	// Remove ICC profile metadata
	if o.NoProfile {
		C.remove_profile(image.c)
	}

	if o.Interpretation > 0 && vipsColourspaceIsSupported(image) {
		// Apply the proper color space.
		if newImage, err := vipsColourspace(image, o.Interpretation); err != nil {
			return nil, err
		} else {
			image = newImage
		}
	}

	if o.OutputICC != "" && o.InputICC != "" {
		outputIccPath := C.CString(o.OutputICC)
		defer C.free(unsafe.Pointer(outputIccPath))

		inputIccPath := C.CString(o.InputICC)
		defer C.free(unsafe.Pointer(inputIccPath))

		err := C.vips_icc_transform_with_default_bridge(image.c, &outImage, outputIccPath, inputIccPath)
		if int(err) != 0 {
			return nil, catchVipsError()
		}
		return wrapVipsImage(outImage), nil
	}

	if o.OutputICC != "" && vipsHasProfile(image) {
		outputIccPath := C.CString(o.OutputICC)
		defer C.free(unsafe.Pointer(outputIccPath))

		err := C.vips_icc_transform_bridge(image.c, &outImage, outputIccPath)
		if int(err) != 0 {
			return nil, catchVipsError()
		}
		image = wrapVipsImage(outImage)
	}

	return image, nil
}

func vipsSave(image *vipsImage, o vipsSaveOptions) ([]byte, error) {
	image, err := vipsPreSave(image, &o)
	if err != nil {
		return nil, err
	}

	length := C.size_t(0)
	saveErr := C.int(0)
	interlace := C.int(boolToInt(o.Interlace))
	quality := C.int(o.Quality)
	strip := C.int(boolToInt(o.StripMetadata))
	lossless := C.int(boolToInt(o.Lossless))
	palette := C.int(boolToInt(o.Palette))
	speed := C.int(o.Speed)

	if o.Type != 0 && !IsTypeSupportedSave(o.Type) {
		return nil, fmt.Errorf("VIPS cannot save to %#v", ImageTypes[o.Type])
	}
	var ptr unsafe.Pointer
	switch o.Type {
	case WEBP:
		saveErr = C.vips_webpsave_bridge(image.c, &ptr, &length, strip, quality, lossless)
	case PNG:
		saveErr = C.vips_pngsave_bridge(image.c, &ptr, &length, strip, C.int(o.Compression), quality, interlace, palette)
	case TIFF:
		saveErr = C.vips_tiffsave_bridge(image.c, &ptr, &length)
	case HEIF:
		saveErr = C.vips_heifsave_bridge(image.c, &ptr, &length, strip, quality, lossless)
	case AVIF:
		saveErr = C.vips_avifsave_bridge(image.c, &ptr, &length, strip, quality, lossless, speed)
	case JP2K:
		saveErr = C.vips_jp2ksave_bridge(image.c, &ptr, &length, strip, quality, lossless)
	case GIF:
		formatString := C.CString("GIF")
		defer C.free(unsafe.Pointer(formatString))
		saveErr = C.vips_magicksave_bridge(image.c, &ptr, &length, formatString, quality)
	case MAGICK:
		formatString := C.CString(o.MagickFormat)
		defer C.free(unsafe.Pointer(formatString))
		saveErr = C.vips_magicksave_bridge(image.c, &ptr, &length, formatString, quality)
	default:
		saveErr = C.vips_jpegsave_bridge(image.c, &ptr, &length, strip, quality, interlace)
	}

	if int(saveErr) != 0 {
		return nil, catchVipsError()
	}

	buf := C.GoBytes(ptr, C.int(length))

	// Clean up
	C.g_free(C.gpointer(ptr))
	C.vips_error_clear()

	return buf, nil
}

func vipsExtract(image *vipsImage, left, top, width, height int) (*vipsImage, error) {
	var out *C.VipsImage

	if width > MaxSize || height > MaxSize {
		return nil, errors.New("maximum image size exceeded")
	}

	top, left = max(top), max(left)
	err := C.vips_extract_area_bridge(image.c, &out, C.int(left), C.int(top), C.int(width), C.int(height))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsSmartCrop(image *vipsImage, width, height int) (*vipsImage, error) {
	var out *C.VipsImage

	if width > MaxSize || height > MaxSize {
		return nil, errors.New("maximum image size exceeded")
	}

	err := C.vips_smartcrop_bridge(image.c, &out, C.int(width), C.int(height))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsTrim(image *vipsImage, background RGBAProvider, threshold float64) (int, int, int, int, error) {
	top := C.int(0)
	left := C.int(0)
	width := C.int(0)
	height := C.int(0)

	if background == nil {
		//return 0, 0, 0, 0, errors.New("cannot trim without a background to look for")
		background = ColorBlack
	}

	r, g, b, _ := background.RGBA()

	err := C.vips_find_trim_bridge(image.c,
		&top, &left, &width, &height,
		C.double(r), C.double(g), C.double(b),
		C.double(threshold))
	if err != 0 {
		return 0, 0, 0, 0, catchVipsError()
	}

	return int(top), int(left), int(width), int(height), nil
}

func vipsShrinkJpeg(buf []byte, shrink int) (*vipsImage, error) {
	var image *C.VipsImage
	var ptr = unsafe.Pointer(&buf[0])

	err := C.vips_jpegload_buffer_shrink(ptr, C.size_t(len(buf)), &image, C.int(shrink))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(image), nil
}

func vipsShrinkWebp(buf []byte, shrink int) (*vipsImage, error) {
	var image *C.VipsImage
	var ptr = unsafe.Pointer(&buf[0])

	err := C.vips_webpload_buffer_shrink(ptr, C.size_t(len(buf)), &image, C.int(shrink))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(image), nil
}

func vipsShrink(input *vipsImage, shrink int) (*vipsImage, error) {
	var image *C.VipsImage

	err := C.vips_shrink_bridge(input.c, &image, C.double(float64(shrink)), C.double(float64(shrink)))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(image), nil
}

func vipsReduce(input *vipsImage, xshrink float64, yshrink float64) (*vipsImage, error) {
	var image *C.VipsImage

	err := C.vips_reduce_bridge(input.c, &image, C.double(xshrink), C.double(yshrink))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(image), nil
}

func vipsColourspace(input *vipsImage, interpretation Interpretation) (*vipsImage, error) {
	var out *C.VipsImage
	err := C.vips_colourspace_bridge(input.c, &out, C.VipsInterpretation(interpretation))
	if int(err) != 0 {
		return nil, catchVipsError()
	}
	return wrapVipsImage(out), nil
}

func vipsAdjustColourspaceToColour(input *vipsImage, color RGBAProvider) (*vipsImage, error) {
	var result = input

	r, g, b, a := color.RGBA()

	hasAlpha := vipsHasAlpha(input)
	if !hasAlpha && a < 0xFF {
		// No alpha channel but the background color is not opaque? Try to add an alpha channel then.
		var withAlpha *C.VipsImage = C.vips_image_new()
		C.vips_addalpha_bridge(result.c, &withAlpha)
		result = wrapVipsImage(withAlpha)
	}

	// In case it's a grayscale image and our desired color is also grayscale, we can keep the
	// colorspace intact.
	channels := input.c.Bands
	if (channels == 1 || channels == 2) && (r == g && g == b && r == b) {
		return result, nil
	}

	// Make sure it's sRGB, since we are dealing with RGB colors. This ensures we don't
	// accidentally work in CMYK or something even weirder instead.
	result, err := vipsColourspace(result, InterpretationSRGB)
	if err != nil {
		return nil, fmt.Errorf("cannot adjust colourspace: %w", err)
	}

	return result, nil
}

func vipsEmbed(input *vipsImage, left, top, width, height int, extend Extend, background RGBAProvider) (*vipsImage, error) {
	var image *C.VipsImage

	// Max extend value, see: https://libvips.github.io/libvips/API/current/libvips-conversion.html#VipsExtend
	if extend > 5 {
		extend = ExtendBackground
	}

	if extend == ExtendBackground && background == nil {
		return nil, errors.New("cannot use ExtendBackground without specifying a background")
	}

	// If it's not ExtendBackground, the values are not really used anyway. Therefore just use black.
	var vipsBackground *C.VipsArrayDouble
	if background == nil {
		vipsBackground = nil
	} else {
		if image, err := vipsAdjustColourspaceToColour(input, background); err != nil {
			return nil, err
		} else {
			input = image
		}

		r, g, b, a := background.RGBA()
		channels := input.c.Bands
		hasAlpha := vipsHasAlpha(input)

		if hasAlpha || a < 0xFF {
			if channels == 2 {
				bgArray := [2]C.double{C.double(r), C.double(a)}
				vipsBackground = C.vips_array_double_new(&bgArray[0], 2)
			} else {
				bgArray := [4]C.double{C.double(r), C.double(g), C.double(b), C.double(a)}
				vipsBackground = C.vips_array_double_new(&bgArray[0], 4)
			}

		} else {
			if channels == 1 {
				bgArray := [1]C.double{C.double(r)}
				vipsBackground = C.vips_array_double_new(&bgArray[0], 1)
			} else {
				bgArray := [3]C.double{C.double(r), C.double(g), C.double(b)}
				vipsBackground = C.vips_array_double_new(&bgArray[0], 3)
			}
		}
	}

	err := C.vips_embed_bridge(input.c, &image, C.int(left), C.int(top), C.int(width),
		C.int(height), C.int(extend), vipsBackground)
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(image), nil
}

func vipsAffine(input *vipsImage, residualx, residualy float64, i Interpolator, extend Extend) (*vipsImage, error) {
	if extend > 5 {
		extend = ExtendBackground
	}

	var image *C.VipsImage
	cstring := C.CString(i.String())
	interpolator := C.vips_interpolate_new(cstring)

	defer C.free(unsafe.Pointer(cstring))
	defer C.g_object_unref(C.gpointer(interpolator))

	err := C.vips_affine_interpolator(input.c, &image, C.double(residualx), 0, 0, C.double(residualy), interpolator, C.int(extend))
	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(image), nil
}

var magickBufferMatch = regexp.MustCompile(`Magick(\d+)?Buffer$`)

func vipsImageType(buf []byte) ImageType {
	if len(buf) < 12 {
		return UNKNOWN
	}
	if bytes.HasPrefix(buf, []byte{0xFF, 0xD8, 0xFF}) {
		return JPEG
	}
	if bytes.HasPrefix(buf, []byte("GIF")) {
		return GIF
	}
	if bytes.HasPrefix(buf, []byte{0x89, 'P', 'N', 'G'}) {
		return PNG
	}
	if IsTypeSupported(TIFF) &&
		(bytes.HasPrefix(buf, []byte{0x49, 0x49, 0x2A, 0x0}) ||
			bytes.HasPrefix(buf, []byte{0x4D, 0x4D, 0x0, 0x2A})) {
		return TIFF
	}
	if IsTypeSupported(PDF) && bytes.HasPrefix(buf, []byte("%PDF")) {
		return PDF
	}
	if IsTypeSupported(WEBP) && string(buf[8:12]) == "WEBP" {
		return WEBP
	}
	if IsTypeSupported(SVG) && IsSVGImage(buf) {
		return SVG
	}
	// NOTE: libheif currently only supports heic sub types; see:
	//   https://github.com/strukturag/libheif/issues/83#issuecomment-421427091
	if IsTypeSupported(HEIF) && string(buf[4:8]) == "ftyp" {
		subType := string(buf[8:12])
		switch subType {
		case "heic", "mif1", "msf1", "heis", "hevc":
			return HEIF
		}
	}
	if IsTypeSupported(JP2K) && (bytes.HasPrefix(buf, []byte{0x0, 0x0, 0x0, 0xC, 0x6A, 0x50, 0x20, 0x20, 0xD, 0xA, 0x87, 0xA}) ||
		bytes.HasPrefix(buf, []byte{0xFF, 0x4F, 0xFF, 0x51})) {
		return JP2K
	}

	// If nothing matched directly, try to fallback to imagemagick (if available).
	if IsTypeSupported(MAGICK) && magickBufferMatch.MatchString(readImageType(buf)) {
		return MAGICK
	}
	if IsTypeSupported(HEIF) && buf[4] == 0x66 && buf[5] == 0x74 && buf[6] == 0x79 && buf[7] == 0x70 &&
		buf[8] == 0x61 && buf[9] == 0x76 && buf[10] == 0x69 && buf[11] == 0x66 {
		return AVIF
	}

	return UNKNOWN
}

func readImageType(buf []byte) string {
	length := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])
	load := C.vips_foreign_find_load_buffer(imageBuf, length)
	return C.GoString(load)
}

func catchVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func vipsGaussianBlur(image *vipsImage, o GaussianBlur) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_gaussblur_bridge(image.c, &out, C.double(o.Sigma), C.double(o.MinAmpl))
	if err != 0 {
		return nil, catchVipsError()
	}
	return wrapVipsImage(out), nil
}

func vipsSharpen(image *vipsImage, o Sharpen) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_sharpen_bridge(image.c, &out, C.int(o.Radius), C.double(o.X1), C.double(o.Y2), C.double(o.Y3), C.double(o.M1), C.double(o.M2))
	if err != 0 {
		return nil, catchVipsError()
	}
	return wrapVipsImage(out), nil
}

func max(x int) int {
	return int(math.Max(float64(x), 0))
}

type drawWatermarkOptions struct {
	Left    int
	Top     int
	Image   *vipsImage
	Opacity float32
}

func vipsDrawWatermark(image *vipsImage, o drawWatermarkOptions) (*vipsImage, error) {
	var out *C.VipsImage

	opts := vipsWatermarkImageOptions{C.int(o.Left), C.int(o.Top), C.float(o.Opacity)}

	err := C.vips_watermark_image(image.c, o.Image.c, &out, (*C.WatermarkImageOptions)(unsafe.Pointer(&opts)))

	if err != 0 {
		return nil, catchVipsError()
	}

	return wrapVipsImage(out), nil
}

func vipsGamma(image *vipsImage, Gamma float64) (*vipsImage, error) {
	var out *C.VipsImage

	err := C.vips_gamma_bridge(image.c, &out, C.double(Gamma))
	if err != 0 {
		return nil, catchVipsError()
	}
	return wrapVipsImage(out), nil
}
