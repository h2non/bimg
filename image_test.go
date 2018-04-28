package bimg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
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

	Write("testdata/test_resize_out.jpg", buf)
}

func TestImageGifResize(t *testing.T) {
	_, err := initImage("test.gif").Resize(300, 240)
	if err == nil {
		t.Errorf("GIF shouldn't be saved within VIPS")
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

		Write("testdata/test_gif.jpg", buf)
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

		Write("testdata/test_pdf.jpg", buf)
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

		Write("testdata/test_svg.jpg", buf)
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

	Write("testdata/test_resize_crop_out.jpg", buf)
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

	Write("testdata/test_extract_out.jpg", buf)
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

	Write("testdata/test_extract_zero_out.jpg", buf)
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

	Write("testdata/test_enlarge_out.jpg", buf)
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

	Write("testdata/test_enlarge_crop_out.jpg", buf)
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

	Write("testdata/test_crop_out.jpg", buf)
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

	Write("testdata/test_crop_width_out.jpg", buf)
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

	Write("testdata/test_crop_height_out.jpg", buf)
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

	Write("testdata/test_thumbnail_out.jpg", buf)
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

	Write("testdata/test_watermark_text_out.jpg", buf)
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

	Write("testdata/test_watermark_image_out.jpg", buf)
}

func TestImageWatermarkWithImageGravity(t *testing.T) {
	tests := []struct {
		name         string
		wm           WatermarkImage
		validateFunc func(buf []byte)
	}{
		{"1nw", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravityNorthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 0, 0, 29, 19) }},
		{"2n", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravityNorth}, func(buf []byte) { validateWatermarkArea(t, buf, 45, 0, 74, 19) }},
		{"3ne", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravityNorthEast}, func(buf []byte) { validateWatermarkArea(t, buf, 90, 0, 119, 19) }},
		{"4w", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravityWest}, func(buf []byte) { validateWatermarkArea(t, buf, 0, 40, 29, 59) }},
		{"5c", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravityCentre}, func(buf []byte) { validateWatermarkArea(t, buf, 45, 40, 74, 59) }},
		{"6e", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravityEast}, func(buf []byte) { validateWatermarkArea(t, buf, 90, 40, 119, 59) }},
		{"7sw", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravitySouthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 0, 80, 29, 99) }},
		{"8s", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravitySouth}, func(buf []byte) { validateWatermarkArea(t, buf, 45, 80, 74, 99) }},
		{"9se", WatermarkImage{X: 0, Y: 0, Gravity: WatermarkGravitySouthEast}, func(buf []byte) { validateWatermarkArea(t, buf, 90, 80, 119, 99) }},
		{"10off10nw", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravityNorthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 10, 10, 39, 29) }},
		{"11off10n", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravityNorth}, func(buf []byte) { validateWatermarkArea(t, buf, 55, 10, 84, 29) }},
		{"12off10ne", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravityNorthEast}, func(buf []byte) { validateWatermarkArea(t, buf, 80, 10, 109, 29) }},
		{"13off10w", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravityWest}, func(buf []byte) { validateWatermarkArea(t, buf, 10, 50, 39, 69) }},
		{"14off10c", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravityCentre}, func(buf []byte) { validateWatermarkArea(t, buf, 55, 50, 84, 69) }},
		{"15off10e", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravityEast}, func(buf []byte) { validateWatermarkArea(t, buf, 80, 50, 109, 69) }},
		{"16off10sw", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravitySouthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 10, 70, 39, 89) }},
		{"17off10s", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravitySouth}, func(buf []byte) { validateWatermarkArea(t, buf, 55, 70, 84, 89) }},
		{"18off10se", WatermarkImage{X: 10, Y: 10, Gravity: WatermarkGravitySouthEast}, func(buf []byte) { validateWatermarkArea(t, buf, 80, 70, 109, 89) }},
		{"19off10cminus", WatermarkImage{X: -10, Y: -10, Gravity: WatermarkGravityCentre}, func(buf []byte) { validateWatermarkArea(t, buf, 35, 30, 64, 49) }},
		{"20pnt10nw", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravityNorthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 12, 10, 41, 29) }},
		{"21pnt10n", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravityNorth}, func(buf []byte) { validateWatermarkArea(t, buf, 57, 10, 86, 29) }},
		{"22pnt10ne", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravityNorthEast}, func(buf []byte) { validateWatermarkArea(t, buf, 78, 10, 107, 29) }},
		{"23pnt10w", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravityWest}, func(buf []byte) { validateWatermarkArea(t, buf, 12, 50, 41, 69) }},
		{"24pnt10c", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravityCentre}, func(buf []byte) { validateWatermarkArea(t, buf, 57, 50, 86, 69) }},
		{"25pnt10e", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravityEast}, func(buf []byte) { validateWatermarkArea(t, buf, 79, 50, 107, 69) }},
		{"26pnt10sw", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravitySouthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 12, 70, 41, 89) }},
		{"27pnt10s", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravitySouth}, func(buf []byte) { validateWatermarkArea(t, buf, 57, 70, 86, 89) }},
		{"28pnt10se", WatermarkImage{XRate: 0.1, YRate: 0.1, Gravity: WatermarkGravitySouthEast}, func(buf []byte) { validateWatermarkArea(t, buf, 78, 70, 107, 89) }},
		{"29pnt10cminus", WatermarkImage{XRate: -0.1, YRate: -0.1, Gravity: WatermarkGravityCentre}, func(buf []byte) { validateWatermarkArea(t, buf, 33, 30, 62, 49) }},
		{"50protruded1", WatermarkImage{X: -10000, Y: -10000, Gravity: WatermarkGravityNorthWest}, func(buf []byte) { validateWatermarkArea(t, buf, 0, 0, 29, 19) }},
		{"51protruded2", WatermarkImage{X: 10000, Y: 10000, Gravity: WatermarkGravityCentre}, func(buf []byte) { validateWatermarkArea(t, buf, 90, 80, 119, 99) }},
		{"60oldway", WatermarkImage{Left: 10, Top: 10}, func(buf []byte) { validateWatermarkArea(t, buf, 10, 10, 39, 29) }},
		{"61none", WatermarkImage{}, func(buf []byte) { validateWatermarkArea(t, buf, 0, 0, 29, 19) }},
	}

	for _, test := range tests {
		image := initImage("white_1000x1000.png")
		watermark, _ := imageBuf("black_30x20.png")

		_, err := image.Crop(120, 100, GravityNorth)
		if err != nil {
			t.Errorf("Cannot process the image: %#v", err)
		}

		test.wm.Buf = watermark
		_, err = image.WatermarkImage(test.wm)
		if err != nil {
			t.Error(err)
		}

		buf, err := image.Colourspace(InterpretationBW)
		if err != nil {
			t.Error(err)
		}

		test.validateFunc(buf)

		Write("testdata/test_watermark_image_gravity_"+test.name+"_out.png", buf)
	}
}

func validateWatermarkArea(t *testing.T, buf []byte, x0, y0, x1, y1 int) {
	t.Helper()

	imgResult, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		t.Error(err)
	}

	imgWidth := imgResult.Bounds().Max.X
	imgHeight := imgResult.Bounds().Max.Y
	assertPixelColor(t, imgResult, x0, y0, color.Black)
	assertPixelColor(t, imgResult, x1, y0, color.Black)
	assertPixelColor(t, imgResult, x1, y1, color.Black)
	assertPixelColor(t, imgResult, x0, y1, color.Black)

	perimeterNW := image.Pt(x0-1, y0-1)
	if perimeterNW.X >= 0 && perimeterNW.Y >= 0 {
		assertPixelColor(t, imgResult, perimeterNW.X, perimeterNW.Y, color.White)
	}
	perimeterNE := image.Pt(x1+1, y0-1)
	if perimeterNE.X < imgWidth && perimeterNE.Y >= 0 {
		assertPixelColor(t, imgResult, perimeterNE.X, perimeterNE.Y, color.White)
	}
	perimeterSW := image.Pt(x1+1, y1+1)
	if perimeterSW.X < imgWidth && perimeterSW.Y < imgHeight {
		assertPixelColor(t, imgResult, perimeterSW.X, perimeterSW.Y, color.White)
	}
	perimeterSE := image.Pt(x0-1, y1+1)
	if perimeterSE.X >= 0 && perimeterSE.Y < imgHeight {
		assertPixelColor(t, imgResult, perimeterSE.X, perimeterSE.Y, color.White)
	}
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

	Write("testdata/test_watermark_replicate_out.jpg", buf)
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

	Write("testdata/test_zoom_out.jpg", buf)
}

func TestImageFlip(t *testing.T) {
	buf, err := initImage("test.jpg").Flip()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("testdata/test_flip_out.jpg", buf)
}

func TestImageFlop(t *testing.T) {
	buf, err := initImage("test.jpg").Flop()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("testdata/test_flop_out.jpg", buf)
}

func TestImageRotate(t *testing.T) {
	buf, err := initImage("test_flip_out.jpg").Rotate(90)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("testdata/test_image_rotate_out.jpg", buf)
}

func TestImageConvert(t *testing.T) {
	buf, err := initImage("test.jpg").Convert(PNG)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("testdata/test_image_convert_out.png", buf)
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
	Write("testdata/test_transparent_image_convert_out.jpg", buf)
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

		interpretation, err := ImageInterpretation(buf)
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

	Write("testdata/test_image_fluent_out.png", image.Image())
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

	Write("testdata/test_smart_crop.jpg", buf)
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

	Write("testdata/transparent_trim.png", buf)
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

	Write("testdata/parameter_trim.png", buf)
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
	return Read(path.Join("testdata", file))
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

func assertPixelColor(t *testing.T, imgResult image.Image, x, y int, colorExpect color.Color) {
	t.Helper()
	ceR, ceG, ceB, ceA := colorExpect.RGBA()
	caR, caG, caB, caA := imgResult.At(x, y).RGBA()
	if ceR != caR || ceG != caG || ceB != caB || ceA != caA {
		t.Error(fmt.Errorf("Expected pixel color (%02X%02X%02X%02X), but actual (%02X%02X%02X%02X) at (%d, %d)", uint8(ceR), uint8(ceG), uint8(ceB), uint8(ceA), uint8(caR), uint8(caG), uint8(caB), uint8(caA), x, y))
	}
}
