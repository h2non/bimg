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
	return getImageTypeName(vipsImageType(buf))
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

func IsTypeSupported(t ImageType) bool {
	return t == JPEG || t == PNG || t == WEBP
}

func IsTypeNameSupported(t string) bool {
	return t == "jpeg" || t == "jpg" ||
		t == "png" || t == "webp"
}
