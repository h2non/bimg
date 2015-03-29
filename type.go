package bimg

type Type struct {
	Name string
	Mime string
}

func DetermineType(buf []byte) *Type {
	return &Type{Name: "jpg", Mime: "image/jpg"}
}
