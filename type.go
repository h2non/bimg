package bimg

const (
	UNKNOWN = iota
	JPEG
	WEBP
	PNG
	TIFF
	MAGICK
)

func DetermineImageType(buf []byte) int {
	return vipsImageType(buf)
}

func DetermineImageTypeName(buf []byte) string {
	imageCode := vipsImageType(buf)
	imageType := "unknown"

	switch {
	case imageCode == JPEG:
		imageType = "jpeg"
		break
	case imageCode == WEBP:
		imageType = "webp"
		break
	case imageCode == PNG:
		imageType = "png"
		break
	case imageCode == TIFF:
		imageType = "png"
		break
	case imageCode == MAGICK:
		imageType = "magick"
		break
	}

	return imageType
}
