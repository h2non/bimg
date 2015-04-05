package bimg

type Image struct {
	buffer []byte
}

func (i *Image) Resize(width int, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
	}
	return i.Process(options)
}

func (i *Image) Extract(top int, left int, width int, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
		Top:    top,
		Left:   left,
	}
	return i.Process(options)
}

func (i *Image) Rotate(a Angle) ([]byte, error) {
	options := Options{Rotate: a}
	return i.Process(options)
}

func (i *Image) Flip() ([]byte, error) {
	options := Options{Flip: VERTICAL}
	return i.Process(options)
}

func (i *Image) Flop() ([]byte, error) {
	options := Options{Flip: HORIZONTAL}
	return i.Process(options)
}

func (i *Image) Convert(t ImageType) ([]byte, error) {
	options := Options{Type: t}
	return i.Process(options)
}

func (i *Image) Process(o Options) ([]byte, error) {
	image, err := Resize(i.buffer, o)
	if err != nil {
		return nil, err
	}
	i.buffer = image
	return image, nil
}

func (i *Image) Type() string {
	return DetermineImageTypeName(i.buffer)
}

func NewImage(buf []byte) *Image {
	return &Image{buf}
}
