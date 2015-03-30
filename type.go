package bimg

const (
	UNKNOWN = "unknown"
	JPEG    = "jpeg"
	WEBP    = "webp"
	PNG     = "png"
	TIFF    = "tiff"
	MAGICK  = "magick"
)

type Type struct {
	Name string
	Mime string
}

func DetermineType(buf []byte) *Type {
	return &Type{Name: "jpg", Mime: "image/jpg"}
}
