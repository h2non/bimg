package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestResize(t *testing.T) {
	options := Options{Width: 800, Height: 600}
	buf, _ := Read("fixtures/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	err = Write("fixtures/test_out.jpg", newImg)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}

func TestRotate(t *testing.T) {
	options := Options{Width: 800, Height: 600, Rotate: 270}
	buf, _ := Read("fixtures/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	err = Write("fixtures/test_rotate_out.jpg", newImg)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}

func TestInvalidRotate(t *testing.T) {
	options := Options{Width: 800, Height: 600, Rotate: 111}
	buf, _ := Read("fixtures/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	err = Write("fixtures/test_invalid_rotate_out.jpg", newImg)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}

func TestConvert(t *testing.T) {
	width, height := 640, 480

	options := Options{Width: width, Height: height, Crop: true, Type: PNG}
	img, err := os.Open("fixtures/test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer img.Close()

	buf, err := ioutil.ReadAll(img)
	if err != nil {
		t.Fatal(err)
	}

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != PNG {
		t.Fatal("Image is not png")
	}

	size, _ := Size(newImg)
	if size.Height != height || size.Width != width {
		t.Fatal("Invalid image size")
	}

	err = Write("fixtures/test_out.png", newImg)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}

func TestResizePngWithTransparency(t *testing.T) {
	width, height := 300, 240

	options := Options{Width: width, Height: height, Crop: true}
	img, err := os.Open("fixtures/transparent.png")
	if err != nil {
		t.Fatal(err)
	}
	defer img.Close()

	buf, err := ioutil.ReadAll(img)
	if err != nil {
		t.Fatal(err)
	}

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != PNG {
		t.Fatal("Image is not png")
	}

	size, _ := Size(newImg)
	if size.Height != height || size.Width != width {
		t.Fatal("Invalid image size")
	}

	err = Write("fixtures/transparent_out.png", newImg)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}

func runBenchmarkResize(file string, o Options, b *testing.B) {
	buf, _ := Read(path.Join("fixtures", file))

	for n := 0; n < b.N; n++ {
		Resize(buf, o)
	}
}

func BenchmarkResizeLargeJpeg(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkResizePng(b *testing.B) {
	options := Options{
		Width:  200,
		Height: 200,
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkResizeWebP(b *testing.B) {
	options := Options{
		Width:  200,
		Height: 200,
	}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkConvertToJpeg(b *testing.B) {
	options := Options{Type: JPEG}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkCrop(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkExtractJpeg(b *testing.B) {
	options := Options{
		Top:        100,
		Left:       50,
		AreaWidth:  600,
		AreaHeight: 480,
	}
	runBenchmarkResize("test.jpg", options, b)
}
