//go:build ignore

package bimg

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func formatOptions(o Options) string {
	var attributes []string
	if o.Crop {
		attributes = append(attributes, "Crop")
	}
	if o.Enlarge {
		attributes = append(attributes, "Enlarge")
	}
	if o.Force {
		attributes = append(attributes, "Force")
	}
	if o.Width > 0 {
		attributes = append(attributes, fmt.Sprintf("Width=%d", o.Width))
	}
	if o.Height > 0 {
		attributes = append(attributes, fmt.Sprintf("Height=%d", o.Height))
	}

	return fmt.Sprintf("Options{%s}", strings.Join(attributes, ","))
}

// TODO: port this test
func TestNoColorProfile(t *testing.T) {
	options := Options{Width: 800, Height: 600, NoProfile: true}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	metadata, err := Metadata(newImg)
	if err != nil {
		t.Fatalf("cannot get metadata: %v", err)
	}
	if metadata.Profile == true {
		t.Fatal("Invalid profile data")
	}

	size, _ := Size(newImg)
	if size.Height != options.Height || size.Width != options.Width {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}
}

func TestExtractWithDefaultAxis(t *testing.T) {
	options := Options{AreaWidth: 200, AreaHeight: 200}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.AreaHeight || size.Width != options.AreaWidth {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	_ = Write("testdata/test_extract_defaults_out.jpg", newImg)
}

func TestExtractCustomAxis(t *testing.T) {
	options := Options{Top: 100, Left: 100, AreaWidth: 200, AreaHeight: 200}
	buf, _ := Read("testdata/test.jpg")

	newImg, err := Resize(buf, options)
	if err != nil {
		t.Errorf("Resize(imgData, %#v) error: %#v", options, err)
	}

	size, _ := Size(newImg)
	if size.Height != options.AreaHeight || size.Width != options.AreaWidth {
		t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
	}

	_ = Write("testdata/test_extract_custom_axis_out.jpg", newImg)
}

func TestExtractOrEmbedImage(t *testing.T) {
	buf, _ := Read("testdata/test.jpg")
	transform, err := NewImageFromBuffer(buf)
	if err != nil {
		t.Fatalf("Unable to load image %s", err)
	}

	o := Options{
		Top:    10,
		Left:   10,
		Width:  100,
		Height: 200,

		// Fields to test
		AreaHeight: 0,
		AreaWidth:  0,
	}

	if err := extractOrEmbedImage(transform, o); err != nil {
		if errors.Unwrap(err) == ErrExtractAreaParamsRequired {
			t.Fatalf("Expecting AreaWidth and AreaHeight to have been defined")
		}

		t.Fatalf("Unknown error occurred %s", err)
	}

	image, err := transform.Save(SaveOptions{Quality: 100})
	if err != nil {
		t.Fatalf("Failed saving image %s", err)
	}

	test, err := Size(image)
	if err != nil {
		t.Fatalf("Failed fetching the size %s", err)
	}

	if test.Height != o.Height {
		t.Errorf("Extract failed, resulting Height %d doesn't match %d", test.Height, o.Height)
	}

	if test.Width != o.Width {
		t.Errorf("Extract failed, resulting Width %d doesn't match %d", test.Width, o.Width)
	}
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
		img, err := os.Open("testdata/" + file)
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
	img, err := os.Open("testdata/transparent.png")
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

	_ = Write("testdata/transparent_out.png", newImg)
}

func TestRotationAndFlip(t *testing.T) {
	files := []struct {
		Name  string
		Angle Angle
		Flip  bool
	}{
		{"Landscape_1", 0, false},
		{"Landscape_2", 0, true},
		{"Landscape_3", D180, false},
		{"Landscape_4", D180, true},
		{"Landscape_5", D90, true},
		{"Landscape_6", D90, false},
		{"Landscape_7", D270, true},
		{"Landscape_8", D270, false},
		{"Portrait_1", 0, false},
		{"Portrait_2", 0, true},
		{"Portrait_3", D180, false},
		{"Portrait_4", D180, true},
		{"Portrait_5", D90, true},
		{"Portrait_6", D90, false},
		{"Portrait_7", D270, true},
		{"Portrait_8", D270, false},
	}

	for _, file := range files {
		t.Run(file.Name, func(t *testing.T) {
			img, err := os.Open(fmt.Sprintf("testdata/exif/%s.jpg", file.Name))
			if err != nil {
				t.Fatal(err)
			}

			buf, err := ioutil.ReadAll(img)
			if err != nil {
				t.Fatal(err)
			}
			img.Close()

			image, _, err := loadImage(buf)
			if err != nil {
				t.Fatal(err)
			}

			angle, flip := calculateRotationAndFlip(image, D0)
			if angle != file.Angle {
				t.Errorf("Rotation for %v expected to be %v. got %v", file.Name, file.Angle, angle)
			}
			if flip != file.Flip {
				t.Errorf("Flip for %v expected to be %v. got %v", file.Name, file.Flip, flip)
			}

			// Visual debugging.
			newImg, err := Resize(buf, Options{})
			if err != nil {
				t.Fatal(err)
			}

			_ = Write(fmt.Sprintf("testdata/exif/%s_out.jpg", file.Name), newImg)
		})
	}
}

func TestIfBothSmartCropOptionsAreIdentical(t *testing.T) {
	benchmarkOptions := Options{Width: 100, Height: 100, Crop: true}
	smartCropOptions := Options{Width: 100, Height: 100, Crop: true, SmartCrop: true}
	gravityOptions := Options{Width: 100, Height: 100, Crop: true, Gravity: GravitySmart}

	testImg, err := os.Open("testdata/northern_cardinal_bird.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	testImgByte, err := ioutil.ReadAll(testImg)
	if err != nil {
		t.Fatal(err)
	}

	scImg, err := Resize(testImgByte, smartCropOptions)
	if err != nil {
		t.Fatal(err)
	}

	gImg, err := Resize(testImgByte, gravityOptions)
	if err != nil {
		t.Fatal(err)
	}

	benchmarkImg, err := Resize(testImgByte, benchmarkOptions)
	if err != nil {
		t.Fatal(err)
	}

	sch, gh, bh := md5.Sum(scImg), md5.Sum(gImg), md5.Sum(benchmarkImg)
	if gh == bh || sch == bh {
		t.Error("Expected both options produce a different result from a standard crop.")
	}

	if sch != gh {
		t.Errorf("Expected both options to result in the same output, %x != %x", sch, gh)
	}
}

func TestSkipCropIfTooSmall(t *testing.T) {
	testCases := []struct {
		name    string
		options Options
	}{
		{
			name: "smart crop",
			options: Options{
				Width:   140,
				Height:  140,
				Crop:    true,
				Gravity: GravitySmart,
			},
		},
		{
			name: "centre crop",
			options: Options{
				Width:   140,
				Height:  140,
				Crop:    true,
				Gravity: GravityCentre,
			},
		},
		{
			name: "embed",
			options: Options{
				Width:  140,
				Height: 140,
				Embed:  true,
			},
		},
		{
			name: "extract",
			options: Options{
				Top:        0,
				Left:       0,
				AreaWidth:  140,
				AreaHeight: 140,
			},
		},
	}

	testImg, err := os.Open("testdata/test_bad_extract_area.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer testImg.Close()

	testImgByte, err := ioutil.ReadAll(testImg)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			croppedImage, err := Resize(testImgByte, tc.options)
			if err != nil {
				t.Fatal(err)
			}

			size, _ := Size(croppedImage)
			if tc.options.Height-size.Height > 1 || tc.options.Width-size.Width > 1 {
				t.Fatalf("Invalid image size: %dx%d", size.Width, size.Height)
			}
			t.Logf("size for %s is %dx%d", tc.name, size.Width, size.Height)
		})
	}
}

func runBenchmarkResize(file string, o Options, b *testing.B) {
	buf, _ := Read(path.Join("testdata", file))

	for n := 0; n < b.N; n++ {
		_, _ = Resize(buf, o)
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

func BenchmarkResizeWebp(b *testing.B) {
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

func BenchmarkCropWebp(b *testing.B) {
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
		WatermarkOptions: WatermarkOptions{
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

func BenchmarkWatermarkPng(b *testing.B) {
	options := Options{
		WatermarkOptions: WatermarkOptions{
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

func BenchmarkWatermarkWebp(b *testing.B) {
	options := Options{
		WatermarkOptions: WatermarkOptions{
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

func BenchmarkWatermarkImageJpeg(b *testing.B) {
	watermark := readFile("transparent.png")
	options := Options{
		WatermarkImage: WatermarkImage{
			Buf:     watermark,
			Opacity: 0.25,
			Left:    100,
			Top:     100,
		},
	}
	runBenchmarkResize("test.jpg", options, b)
}

func BenchmarkWatermarkImagePng(b *testing.B) {
	watermark := readFile("transparent.png")
	options := Options{
		WatermarkImage: WatermarkImage{
			Buf:     watermark,
			Opacity: 0.25,
			Left:    100,
			Top:     100,
		},
	}
	runBenchmarkResize("test.png", options, b)
}

func BenchmarkWatermarkImageWebp(b *testing.B) {
	watermark := readFile("transparent.png")
	options := Options{
		WatermarkImage: WatermarkImage{
			Buf:     watermark,
			Opacity: 0.25,
			Left:    100,
			Top:     100,
		},
	}
	runBenchmarkResize("test.webp", options, b)
}
