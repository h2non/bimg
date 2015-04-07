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

// Determines the image type format (jpeg, png, webp or tiff)
func DetermineImageType(buf []byte) ImageType {
	return vipsImageType(buf)
}

// Determines the image type format by name (jpeg, png, webp or tiff)
func DetermineImageTypeName(buf []byte) string {
	return getImageTypeName(vipsImageType(buf))
}

// Check if a given image type is supported
func IsTypeSupported(t ImageType) bool {
	return t == JPEG || t == PNG || t == WEBP
}

// Check if a given image type name is supported
func IsTypeNameSupported(t string) bool {
	return t == "jpeg" || t == "jpg" ||
		t == "png" || t == "webp"
}

func getImageTypeName(code ImageType) string {
	imageType := "unknown"

	switch {
	case code == JPEG:
		imageType = "jpeg"
		break
	case code == WEBP:
		imageType = "webp"
		break
	case code == PNG:
		imageType = "png"
		break
	case code == TIFF:
		imageType = "tiff"
		break
	case code == MAGICK:
		imageType = "magick"
		break
	}

	return imageType
}
