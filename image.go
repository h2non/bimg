package bimg

type Image struct {
	buffer []byte
}

// Resize the image to fixed width and height
func (i *Image) Resize(width, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
		Embed:  true,
	}
	return i.Process(options)
}

// Resize the image to fixed width and height with additional crop transformation
func (i *Image) ResizeAndCrop(width, height int) ([]byte, error) {
	options := Options{
		Width:  width,
		Height: height,
		Embed:  true,
		Crop:   true,
	}
	return i.Process(options)
}

// Extract area from the by X/Y axis
func (i *Image) Extract(top, left, width, height int) ([]byte, error) {
	options := Options{
		Top:        top,
		Left:       left,
		AreaWidth:  width,
		AreaHeight: height,
	}
	return i.Process(options)
}

// Enlarge the image by width and height. Aspect ratio is maintained
func (i *Image) Enlarge(width, height int) ([]byte, error) {
	options := Options{
		Width:   width,
		Height:  height,
		Enlarge: true,
	}
	return i.Process(options)
}

// Enlarge the image by width and height with additional crop transformation
func (i *Image) EnlargeAndCrop(width, height int) ([]byte, error) {
	options := Options{
		Width:   width,
		Height:  height,
		Enlarge: true,
		Crop:    true,
	}
	return i.Process(options)
}

// Crop the image to the exact size specified
func (i *Image) Crop(width, height int, gravity Gravity) ([]byte, error) {
	options := Options{
		Width:   width,
		Height:  height,
		Gravity: gravity,
		Crop:    true,
	}
	return i.Process(options)
}

// Crop an image by width (auto height)
func (i *Image) CropByWidth(width int) ([]byte, error) {
	options := Options{
		Width: width,
		Crop:  true,
	}
	return i.Process(options)
}

// Crop an image by height (auto width)
func (i *Image) CropByHeight(height int) ([]byte, error) {
	options := Options{
		Height: height,
		Crop:   true,
	}
	return i.Process(options)
}

// Thumbnail the image by the a given width by aspect ratio 4:4
func (i *Image) Thumbnail(pixels int) ([]byte, error) {
	options := Options{
		Width:   pixels,
		Height:  pixels,
		Crop:    true,
		Quality: 95,
	}
	return i.Process(options)
}

// Add text as watermark on the given image
func (i *Image) Watermark(w Watermark) ([]byte, error) {
	options := Options{Watermark: w}
	return i.Process(options)
}

// Zoom the image by the given factor.
// You should probably call Extract() before
func (i *Image) Zoom(factor int) ([]byte, error) {
	options := Options{Zoom: factor}
	return i.Process(options)
}

// Rotate the image by given angle degrees (0, 90, 180 or 270)
func (i *Image) Rotate(a Angle) ([]byte, error) {
	options := Options{Rotate: a}
	return i.Process(options)
}

// Flip the image about the vertical Y axis
func (i *Image) Flip() ([]byte, error) {
	options := Options{Flip: true}
	return i.Process(options)
}

// Flop the image about the horizontal X axis
func (i *Image) Flop() ([]byte, error) {
	options := Options{Flop: true}
	return i.Process(options)
}

// Convert image to another format
func (i *Image) Convert(t ImageType) ([]byte, error) {
	options := Options{Type: t}
	return i.Process(options)
}

// Transform the image by custom options
func (i *Image) Process(o Options) ([]byte, error) {
	image, err := Resize(i.buffer, o)
	if err != nil {
		return nil, err
	}
	i.buffer = image
	return image, nil
}

// Get image metadata (size, alpha channel, profile, EXIF rotation)
func (i *Image) Metadata() (ImageMetadata, error) {
	return Metadata(i.buffer)
}

// Get image type format (jpeg, png, webp, tiff)
func (i *Image) Type() string {
	return DetermineImageTypeName(i.buffer)
}

// Get image size
func (i *Image) Size() (ImageSize, error) {
	return Size(i.buffer)
}

// Get image buffer
func (i *Image) Image() []byte {
	return i.buffer
}

// Creates a new image
func NewImage(buf []byte) *Image {
	return &Image{buf}
}
