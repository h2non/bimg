package bimg

// ImageType represents an image type value.
type ImageType int

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
	// MAGICK represents the libmagick compatible genetic image type.
	MAGICK
)

// ImageTypes stores as pairs of image types supported and its alias names.
var ImageTypes = map[ImageType]string{
	JPEG:   "jpeg",
	PNG:    "png",
	WEBP:   "webp",
	TIFF:   "tiff",
	MAGICK: "magick",
}

// DetermineImageType determines the image type format (jpeg, png, webp or tiff)
func DetermineImageType(buf []byte) ImageType {
	return vipsImageType(buf)
}

// DetermineImageTypeName determines the image type format by name (jpeg, png, webp or tiff)
func DetermineImageTypeName(buf []byte) string {
	return ImageTypeName(vipsImageType(buf))
}

// IsTypeSupported checks if a given image type is supported
func IsTypeSupported(t ImageType) bool {
	return ImageTypes[t] != ""
}

// IsTypeNameSupported checks if a given image type name is supported
func IsTypeNameSupported(t string) bool {
	for _, name := range ImageTypes {
		if name == t {
			return true
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
