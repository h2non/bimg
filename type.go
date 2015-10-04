package bimg

type ImageType int

const (
	UNKNOWN ImageType = iota
	JPEG
	WEBP
	PNG
	TIFF
	MAGICK
)

// Pairs of image type and its name
var ImageTypes = map[ImageType]string{
	JPEG:   "jpeg",
	PNG:    "png",
	WEBP:   "webp",
	TIFF:   "tiff",
	MAGICK: "magick",
}

// Determines the image type format (jpeg, png, webp or tiff)
func DetermineImageType(buf []byte) ImageType {
	return vipsImageType(buf)
}

// Determines the image type format by name (jpeg, png, webp or tiff)
func DetermineImageTypeName(buf []byte) string {
	return ImageTypeName(vipsImageType(buf))
}

// Check if a given image type is supported
func IsTypeSupported(t ImageType) bool {
	return ImageTypes[t] != ""
}

// Check if a given image type name is supported
func IsTypeNameSupported(t string) bool {
	for _, name := range ImageTypes {
		if name == t {
			return true
		}
	}
	return false
}

func ImageTypeName(t ImageType) string {
	imageType := ImageTypes[t]
	if imageType == "" {
		return "unknown"
	}
	return imageType
}
