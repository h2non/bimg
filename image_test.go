package vimgo

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestImageResize(t *testing.T) {
	buf, err := initImage("test.jpg").Resize(300, 240)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 300, 240)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_resize_out.jpg", buf, 0664)
}

func TestImageGifResize(t *testing.T) {
	if VipsMajorVersion >= 8 && VipsMinorVersion >= 12 {
		buf, err := initImage("test.gif").Resize(300, 240)
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}

		err = assertSize(buf, 300, 240)
		if err != nil {
			t.Error(err)
		}

		os.WriteFile("testdata/test_resize_out.gif", buf, 0664)
	}
}

func TestImagePdfResize(t *testing.T) {
	_, err := initImage("test.pdf").Resize(300, 240)
	if err == nil {
		t.Errorf("PDF cannot be saved within VIPS")
	}
}

func TestImageSvgResize(t *testing.T) {
	_, err := initImage("test.svg").Resize(300, 240)
	if err == nil {
		t.Errorf("SVG cannot be saved within VIPS")
	}
}

func TestImageGifToJpeg(t *testing.T) {
	if VipsMajorVersion >= 8 && VipsMinorVersion > 2 {
		i := initImage("test.gif")
		options := Options{
			Type: JPEG,
		}
		buf, err := i.Process(options)
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}

		os.WriteFile("testdata/test_gif.jpg", buf, 0644)
	}
}

func TestImagePdfToJpeg(t *testing.T) {
	if VipsMajorVersion >= 8 && VipsMinorVersion > 2 {
		i := initImage("test.pdf")
		options := Options{
			Type: JPEG,
		}
		buf, err := i.Process(options)
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}

		os.WriteFile("testdata/test_pdf.jpg", buf, 0644)
	}
}

func TestImageSvgToJpeg(t *testing.T) {
	if VipsMajorVersion >= 8 && VipsMinorVersion > 2 {
		i := initImage("test.svg")
		options := Options{
			Type: JPEG,
		}
		buf, err := i.Process(options)
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}

		os.WriteFile("testdata/test_svg.jpg", buf, 0644)
	}
}

func TestImageResizeAndCrop(t *testing.T) {
	buf, err := initImage("test.jpg").ResizeAndCrop(300, 200)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 300, 200)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_resize_crop_out.jpg", buf, 0644)
}

func TestImageExtract(t *testing.T) {
	buf, err := initImage("test.jpg").Extract(100, 100, 300, 200)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 300, 200)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_extract_out.jpg", buf, 0644)
}

func TestImageExtractZero(t *testing.T) {
	buf, err := initImage("test.jpg").Extract(0, 0, 300, 200)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 300, 200)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_extract_zero_out.jpg", buf, 0644)
}

func TestImageEnlarge(t *testing.T) {
	buf, err := initImage("test.png").Enlarge(500, 375)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 500, 375)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_enlarge_out.jpg", buf, 0644)
}

func TestImageEnlargeAndCrop(t *testing.T) {
	buf, err := initImage("test.png").EnlargeAndCrop(800, 480)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 800, 480)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_enlarge_crop_out.jpg", buf, 0644)
}

func TestImageCrop(t *testing.T) {
	buf, err := initImage("test.jpg").Crop(800, 600, GravityNorth)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 800, 600)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_crop_out.jpg", buf, 0644)
}

func TestImageCropByWidth(t *testing.T) {
	buf, err := initImage("test.jpg").CropByWidth(600)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 600, 1050)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_crop_width_out.jpg", buf, 0644)
}

func TestImageCropByHeight(t *testing.T) {
	buf, err := initImage("test.jpg").CropByHeight(300)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 1680, 300)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_crop_height_out.jpg", buf, 0644)
}

func TestImageThumbnail(t *testing.T) {
	buf, err := initImage("test.jpg").Thumbnail(100)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 100, 100)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_thumbnail_out.jpg", buf, 0644)
}

func TestImageWatermark(t *testing.T) {
	image := initImage("test.jpg")
	_, err := image.Crop(800, 600, GravityNorth)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	buf, err := image.Watermark(Watermark{
		Text:       "Copy me if you can",
		Opacity:    0.5,
		Width:      200,
		DPI:        100,
		Background: Color{255, 255, 255},
	})
	if err != nil {
		t.Error(err)
	}

	err = assertSize(buf, 800, 600)
	if err != nil {
		t.Error(err)
	}

	if DetermineImageType(buf) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	os.WriteFile("testdata/test_watermark_text_out.jpg", buf, 0644)
}

func TestImageWatermarkWithImage(t *testing.T) {
	image := initImage("test.jpg")
	watermark, _ := imageBuf("transparent.png")

	_, err := image.Crop(800, 600, GravityNorth)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	buf, err := image.WatermarkImage(WatermarkImage{Left: 100, Top: 100, Buf: watermark})

	if err != nil {
		t.Error(err)
	}

	err = assertSize(buf, 800, 600)
	if err != nil {
		t.Error(err)
	}

	if DetermineImageType(buf) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	os.WriteFile("testdata/test_watermark_image_out.jpg", buf, 0644)
}

func TestImageWatermarkNoReplicate(t *testing.T) {
	image := initImage("test.jpg")
	_, err := image.Crop(800, 600, GravityNorth)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	buf, err := image.Watermark(Watermark{
		Text:        "Copy me if you can",
		Opacity:     0.5,
		Width:       200,
		DPI:         100,
		NoReplicate: true,
		Background:  Color{255, 255, 255},
	})
	if err != nil {
		t.Error(err)
	}

	err = assertSize(buf, 800, 600)
	if err != nil {
		t.Error(err)
	}

	if DetermineImageType(buf) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	os.WriteFile("testdata/test_watermark_replicate_out.jpg", buf, 0644)
}

func TestImageZoom(t *testing.T) {
	image := initImage("test.jpg")

	_, err := image.Extract(100, 100, 400, 300)
	if err != nil {
		t.Errorf("Cannot extract the image: %s", err)
	}

	buf, err := image.Zoom(1)
	if err != nil {
		t.Errorf("Cannot process the image: %s", err)
	}

	err = assertSize(buf, 800, 600)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_zoom_out.jpg", buf, 0644)
}

func TestImageFlip(t *testing.T) {
	buf, err := initImage("test.jpg").Flip()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	os.WriteFile("testdata/test_flip_out.jpg", buf, 0644)
}

func TestImageFlop(t *testing.T) {
	buf, err := initImage("test.jpg").Flop()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	os.WriteFile("testdata/test_flop_out.jpg", buf, 0644)
}

func TestImageRotate(t *testing.T) {
	buf, err := initImage("test_flip_out.jpg").Rotate(90)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	os.WriteFile("testdata/test_image_rotate_out.jpg", buf, 0644)
}

func TestImageAutoRotate(t *testing.T) {
	if VipsMajorVersion <= 8 && VipsMinorVersion < 10 {
		t.Skip("Skip test in libvips < 8.10")
		return
	}

	tests := []struct {
		file        string
		orientation int
	}{
		{"exif/Landscape_1.jpg", 1},
		{"exif/Landscape_2.jpg", 1},
		{"exif/Landscape_3.jpg", 1},
		{"exif/Landscape_4.jpg", 1},
		{"exif/Landscape_5.jpg", 1},
		{"exif/Landscape_6.jpg", 1},
		{"exif/Landscape_7.jpg", 1},
	}

	for index, test := range tests {
		img := initImage(test.file)
		buf, err := img.AutoRotate()
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}
		os.WriteFile(fmt.Sprintf("testdata/test_autorotate_%d_out.jpg", index), buf, 0644)

		meta, err := img.Metadata()
		if err != nil {
			t.Errorf("Cannot read image metadata: %#v", err)
		}
		if meta.Orientation != test.orientation {
			t.Errorf("Invalid image orientation for %s: %d != %d", test.file, meta.Orientation, test.orientation)
		}
	}
}

func TestImageConvert(t *testing.T) {
	buf, err := initImage("test.jpg").Convert(PNG)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	os.WriteFile("testdata/test_image_convert_out.png", buf, 0644)
}

func TestTransparentImageConvert(t *testing.T) {
	image := initImage("transparent.png")
	options := Options{
		Type:       JPEG,
		Background: Color{255, 255, 255},
	}
	buf, err := image.Process(options)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	os.WriteFile("testdata/test_transparent_image_convert_out.jpg", buf, 0644)
}

func TestImageMetadata(t *testing.T) {
	data, err := initImage("test.png").Metadata()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	if data.Alpha != true {
		t.Fatal("Invalid alpha channel")
	}
	if data.Size.Width != 400 {
		t.Fatal("Invalid width size")
	}
	if data.Type != "png" {
		t.Fatal("Invalid image type")
	}
}

func TestInterpretation(t *testing.T) {
	interpretation, err := initImage("test.jpg").Interpretation()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	if interpretation != InterpretationSRGB {
		t.Errorf("Invalid interpretation: %d", interpretation)
	}
}

func TestImageColourspace(t *testing.T) {
	tests := []struct {
		file           string
		interpretation Interpretation
	}{
		{"test.jpg", InterpretationSRGB},
		{"test.jpg", InterpretationBW},
	}

	for _, test := range tests {
		buf, err := initImage(test.file).Colourspace(test.interpretation)
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}

		interpretation, _ := ImageInterpretation(buf)
		if interpretation != test.interpretation {
			t.Errorf("Invalid colourspace")
		}
	}
}

func TestImageColourspaceIsSupported(t *testing.T) {
	supported, err := initImage("test.jpg").ColourspaceIsSupported()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	if supported != true {
		t.Errorf("Non-supported colourspace")
	}
}

func TestFluentInterface(t *testing.T) {
	image := initImage("test.jpg")
	_, err := image.CropByWidth(300)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	_, err = image.Flip()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	_, err = image.Convert(PNG)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	data, _ := image.Metadata()
	if data.Alpha != false {
		t.Fatal("Invalid alpha channel")
	}
	if data.Size.Width != 300 {
		t.Fatal("Invalid width size")
	}
	if data.Type != "png" {
		t.Fatal("Invalid image type")
	}

	os.WriteFile("testdata/test_image_fluent_out.png", image.Image(), 0644)
}

func TestImageSmartCrop(t *testing.T) {

	if !(VipsMajorVersion >= 8 && VipsMinorVersion >= 5) {
		t.Skipf("Skipping this test, libvips doesn't meet version requirement %s >= 8.5", VipsVersion)
	}

	i := initImage("northern_cardinal_bird.jpg")
	buf, err := i.SmartCrop(300, 300)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 300, 300)
	if err != nil {
		t.Error(err)
	}

	os.WriteFile("testdata/test_smart_crop.jpg", buf, 0644)
}

func TestImageTrim(t *testing.T) {

	if !(VipsMajorVersion >= 8 && VipsMinorVersion >= 6) {
		t.Skipf("Skipping this test, libvips doesn't meet version requirement %s >= 8.6", VipsVersion)
	}

	i := initImage("transparent.png")
	buf, err := i.Trim()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 250, 208)
	if err != nil {
		t.Errorf("The image wasn't trimmed.")
	}

	os.WriteFile("testdata/transparent_trim.png", buf, 0644)
}

func TestImageTrimParameters(t *testing.T) {

	if !(VipsMajorVersion >= 8 && VipsMinorVersion >= 6) {
		t.Skipf("Skipping this test, libvips doesn't meet version requirement %s >= 8.6", VipsVersion)
	}

	i := initImage("test.png")
	options := Options{
		Trim:       true,
		Background: Color{0.0, 0.0, 0.0},
		Threshold:  10.0,
	}
	buf, err := i.Process(options)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buf, 400, 257)
	if err != nil {
		t.Errorf("The image wasn't trimmed.")
	}

	os.WriteFile("testdata/parameter_trim.png", buf, 0644)
}

func TestImageLength(t *testing.T) {
	i := initImage("test.jpg")

	actual := i.Length()
	expected := 53653

	if expected != actual {
		t.Errorf("Size in Bytes of the image doesn't correspond. %d != %d", expected, actual)
	}
}

func initImage(file string) *Image {
	buf, _ := imageBuf(file)
	return NewImage(buf)
}

func imageBuf(file string) ([]byte, error) {
	return os.ReadFile(path.Join("testdata", file))
}

func assertSize(buf []byte, width, height int) error {
	size, err := NewImage(buf).Size()
	if err != nil {
		return err
	}
	if size.Width != width || size.Height != height {
		return fmt.Errorf("Invalid image size: %dx%d", size.Width, size.Height)
	}
	return nil
}
