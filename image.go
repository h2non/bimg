package bimg

type Image struct {
	buffer []byte
}

func (i *Image) Resize(width int, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
	}
	return Resize(i.buffer, options)
}

func (i *Image) Extract(top int, left int, width int, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
		Top:    top,
		Left:   left,
	}
	return Resize(i.buffer, options)
}

func (i *Image) Rotate(degrees Angle) ([]byte, error) {
	options := Options{Rotate: degrees}
	return Resize(i.buffer, options)
}

func (i *Image) Flip() ([]byte, error) {
	options := Options{Flip: VERTICAL}
	return Resize(i.buffer, options)
}

func (i *Image) Flop() ([]byte, error) {
	options := Options{Flip: HORIZONTAL}
	return Resize(i.buffer, options)
}

func (i *Image) Type() string {
	return DetermineImageTypeName(i.buffer)
}

func NewImage(buf []byte) *Image {
	return &Image{buf}
}
