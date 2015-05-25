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

	Write("fixtures/test_out.jpg", newImg)
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

	Write("fixtures/test_rotate_out.jpg", newImg)
}

func TestCorruptedImage(t *testing.T) {
	options := Options{Width: 800, Height: 600}
	buf, _ := Read("fixtures/corrupt.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	Write("fixtures/test_corrupt_out.jpg", newImg)
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

	Write("fixtures/test_invalid_rotate_out.jpg", newImg)
}

func TestConvert(t *testing.T) {
	width, height := 300, 240
	formats := [3]ImageType{PNG, WEBP, JPEG}

	files := []string{
		"test.jpg",
		"test.png",
		"test.webp",
	}

	for _, file := range files {
		img, err := os.Open("fixtures/" + file)
		if err != nil {
			t.Fatal(err)
		}

		buf, err := ioutil.ReadAll(img)
		if err != nil {
			t.Fatal(err)
		}
		img.Close()

		for _, format := range formats {
			options := Options{Width: width, Height: height, Crop: true, Type: format}

			newImg, err := Resize(buf, options)
			if err != nil {
				t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
			}

			if DetermineImageType(newImg) != format {
				t.Fatal("Image is not png")
			}

			size, _ := Size(newImg)
			if size.Height != height || size.Width != width {
				t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
			}
		}
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

	Write("fixtures/transparent_out.png", newImg)
}

func runBenchmarkResize(file string, o Options, b *testing.B) {
	buf, _ := Read(path.Join("fixtures", file))

	for n := 0; n < b.N; n++ {
		Resize(buf, o)
	}
}

func BenchmarkRotateJpeg(b *testing.B) {
	options := Options{Rotate: 180}
	runBenchmarkResize("test.jpg", options, b)
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

func BenchmarkConvertToPng(b *testing.B) {
	options := Options{Type: PNG}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkConvertToWebp(b *testing.B) {
	options := Options{Type: WEBP}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkCropJpeg(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkCropPng(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkCropWebP(b *testing.B) {
	options := Options{
		Width:  800,
		Height: 600,
	}
	runBenchmarkResize("test.webp", options, b)
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

func BenchmarkExtractPng(b *testing.B) {
	options := Options{
		Top:        100,
		Left:       50,
		AreaWidth:  600,
		AreaHeight: 480,
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkExtractWebp(b *testing.B) {
	options := Options{
		Top:        100,
		Left:       50,
		AreaWidth:  600,
		AreaHeight: 480,
	}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkZoomJpeg(b *testing.B) {
	options := Options{Zoom: 1}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkZoomPng(b *testing.B) {
	options := Options{Zoom: 1}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkZoomWebp(b *testing.B) {
	options := Options{Zoom: 1}
	runBenchmarkResize("test.webp", options, b)
}

func BenchmarkWatermarkJpeg(b *testing.B) {
	options := Options{
		Watermark: Watermark{
			Text:       "Chuck Norris (c) 2315",
			Opacity:    0.25,
			Width:      200,
			DPI:        100,
			Margin:     150,
			Font:       "sans bold 12",
			Background: Color{255, 255, 255},
		},
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkWatermarPng(b *testing.B) {
	options := Options{
		Watermark: Watermark{
			Text:       "Chuck Norris (c) 2315",
			Opacity:    0.25,
			Width:      200,
			DPI:        100,
			Margin:     150,
			Font:       "sans bold 12",
			Background: Color{255, 255, 255},
		},
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkWatermarWebp(b *testing.B) {
	options := Options{
		Watermark: Watermark{
			Text:       "Chuck Norris (c) 2315",
			Opacity:    0.25,
			Width:      200,
			DPI:        100,
			Margin:     150,
			Font:       "sans bold 12",
			Background: Color{255, 255, 255},
		},
	}
	runBenchmarkResize("test.webp", options, b)
}
