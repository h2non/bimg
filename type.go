package bimg

const (
	UNKNOWN = "unknown"
	JPEG    = "jpeg"
	WEBP    = "webp"
	PNG     = "png"
	TIFF    = "tiff"
	MAGICK  = "magick"
)

func DetermineType(buf []byte) string {
	return vipsImageType(buf)
}
