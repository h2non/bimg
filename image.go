package bimg

type Image struct {
	buffer []byte
}

func (i *Image) Resize(width, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
	}
	return i.Process(options)
}

func (i *Image) Extract(top, left, width, height int) ([]byte, error) {
	options := Options{
		Top:        top,
		Left:       left,
		AreaWidth:  width,
		AreaHeight: height,
	}
	return i.Process(options)
}

func (i *Image) Crop(width, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
		Crop:   true,
	}
	return i.Process(options)
}

func (i *Image) Thumbnail(pixels int) ([]byte, error) {
	options := Options{
		Width:   pixels,
		Height:  pixels,
		Crop:    true,
		Quality: 95,
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

func (i *Image) Metadata() (ImageMetadata, error) {
	return Metadata(i.buffer)
}

func (i *Image) Size() (ImageSize, error) {
	return Size(i.buffer)
}

func NewImage(buf []byte) *Image {
	return &Image{buf}
}
