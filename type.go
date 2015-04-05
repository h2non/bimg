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

func DetermineImageType(buf []byte) ImageType {
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
		imageType = "tiff"
		break
	case imageCode == MAGICK:
		imageType = "magick"
		break
	}

	return imageType
}

func IsTypeSupported(t ImageType) bool {
	return t == JPEG || t == PNG || t == WEBP
}

func IsTypeNameSupported(t string) bool {
	return t == "jpeg" || t == "jpg" ||
		t == "png" || t == "webp"
}
