package bimg

import (
	"regexp"
	"sync"
	"unicode/utf8"
)

const (
	// UNKNOWN represents an unknow image type value.
	UNKNOWN ImageType = iota
	// JPEG represents the JPEG image type.
	JPEG
	// WEBP represents the WEBP image type.
	WEBP
	// PNG represents the PNG image type.
	PNG
	// TIFF represents the TIFF image type.
	TIFF
	// GIF represents the GIF image type.
	GIF
	// PDF represents the PDF type.
	PDF
	// SVG represents the SVG image type.
	SVG
	// MAGICK represents the libmagick compatible genetic image type.
	MAGICK
)

// ImageType represents an image type value.
type ImageType int

var (
	htmlCommentRegex = regexp.MustCompile("(?i)<!--([\\s\\S]*?)-->")
	svgRegex         = regexp.MustCompile(`(?i)^\s*(?:<\?xml[^>]*>\s*)?(?:<!doctype svg[^>]*>\s*)?<svg[^>]*>[^*]*<\/svg>\s*$`)
)

// ImageTypes stores as pairs of image types supported and its alias names.
var ImageTypes = map[ImageType]string{
	JPEG:   "jpeg",
	PNG:    "png",
	WEBP:   "webp",
	TIFF:   "tiff",
	GIF:    "gif",
	PDF:    "pdf",
	SVG:    "svg",
	MAGICK: "magick",
}

// imageMutex is used to provide thread-safe synchronization
// for SupportedImageTypes map.
var imageMutex = &sync.RWMutex{}

// SupportedImageTypes stores the optional image type supported
// by the current libvips compilation.
// Note: lazy evaluation as demand is required due
// to bootstrap runtime limitation with C/libvips world.
var SupportedImageTypes = map[ImageType]bool{}

// discoverSupportedImageTypes is used to fill SupportedImageTypes map.
func discoverSupportedImageTypes() {
	imageMutex.Lock()
	SupportedImageTypes[JPEG] = VipsIsTypeSupported(JPEG)
	SupportedImageTypes[PNG] = VipsIsTypeSupported(PNG)
	SupportedImageTypes[GIF] = VipsIsTypeSupported(GIF)
	SupportedImageTypes[WEBP] = VipsIsTypeSupported(WEBP)
	SupportedImageTypes[SVG] = VipsIsTypeSupported(SVG)
	SupportedImageTypes[TIFF] = VipsIsTypeSupported(TIFF)
	SupportedImageTypes[PDF] = VipsIsTypeSupported(PDF)
	SupportedImageTypes[MAGICK] = VipsIsTypeSupported(MAGICK)
	imageMutex.Unlock()
}

// isBinary checks if the given buffer is a binary file.
func isBinary(buf []byte) bool {
	if len(buf) < 24 {
		return false
	}
	for i := 0; i < 24; i++ {
		charCode, _ := utf8.DecodeRuneInString(string(buf[i]))
		if charCode == 65533 || charCode <= 8 {
			return true
		}
	}
	return false
}

// IsSVGImage returns true if the given buffer is a valid SVG image.
func IsSVGImage(buf []byte) bool {
	return !isBinary(buf) && svgRegex.Match(htmlCommentRegex.ReplaceAll(buf, []byte{}))
}

// DetermineImageType determines the image type format (jpeg, png, webp or tiff)
func DetermineImageType(buf []byte) ImageType {
	return vipsImageType(buf)
}

// DetermineImageTypeName determines the image type format by name (jpeg, png, webp or tiff)
func DetermineImageTypeName(buf []byte) string {
	return ImageTypeName(vipsImageType(buf))
}

// IsImageTypeSupportedByVips returns true if the given image type
// is supported by current libvips compilation.
func IsImageTypeSupportedByVips(t ImageType) bool {
	imageMutex.RLock()

	// Discover supported image types and cache the result
	itShouldDiscover := len(SupportedImageTypes) == 0
	if itShouldDiscover {
		imageMutex.RUnlock()
		discoverSupportedImageTypes()
	}

	// Check if image type is actually supported
	isSupported, ok := SupportedImageTypes[t]
	if !itShouldDiscover {
		imageMutex.RUnlock()
	}

	return ok && isSupported
}

// IsTypeSupported checks if a given image type is supported
func IsTypeSupported(t ImageType) bool {
	_, ok := ImageTypes[t]
	return ok && IsImageTypeSupportedByVips(t)
}

// IsTypeNameSupported checks if a given image type name is supported
func IsTypeNameSupported(t string) bool {
	for imageType, name := range ImageTypes {
		if name == t {
			return IsImageTypeSupportedByVips(imageType)
		}
	}
	return false
}

// ImageTypeName is used to get the human friendly name of an image format.
func ImageTypeName(t ImageType) string {
	imageType := ImageTypes[t]
	if imageType == "" {
		return "unknown"
	}
	return imageType
}
