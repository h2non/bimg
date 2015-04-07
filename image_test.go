package bimg

import (
	"path"
	"testing"
)

func TestImageResize(t *testing.T) {
	buf, err := initImage("test.jpg").Resize(300, 240)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 300, 240) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_resize_out.jpg", buf)
}

func TestImageExtract(t *testing.T) {
	buf, err := initImage("test.jpg").Extract(100, 100, 300, 300)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 300, 300) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_extract_out.jpg", buf)
}

func TestImageEnlarge(t *testing.T) {
	buf, err := initImage("test.png").Enlarge(500, 380)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 500, 380) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_enlarge_out.jpg", buf)
}

func TestImageCrop(t *testing.T) {
	buf, err := initImage("test.jpg").Crop(800, 600)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 800, 600) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_crop_out.jpg", buf)
}

func TestImageCropByWidth(t *testing.T) {
	buf, err := initImage("test.jpg").CropByWidth(600)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 600, 375) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_crop_width_out.jpg", buf)
}

func TestImageCropByHeight(t *testing.T) {
	buf, err := initImage("test.jpg").CropByHeight(300)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 800, 480) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_crop_height_out.jpg", buf)
}

func TestImageThumbnail(t *testing.T) {
	buf, err := initImage("test.jpg").Thumbnail(100)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	if assertSize(buf, 100, 100) {
		t.Error("Invalid image size")
	}

	Write("fixtures/test_thumbnail_out.jpg", buf)
}

func TestImageFlip(t *testing.T) {
	buf, err := initImage("test.jpg").Flip()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("fixtures/test_flip_out.jpg", buf)
}

func TestImageRotate(t *testing.T) {
	buf, err := initImage("test_flip_out.jpg").Rotate(90)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("fixtures/test_image_rotate_out.jpg", buf)
}

func TestImageConvert(t *testing.T) {
	buf, err := initImage("test.jpg").Convert(PNG)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	Write("fixtures/test_image_convert_out.png", buf)
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

func initImage(file string) *Image {
	buf, _ := Read(path.Join("fixtures", file))
	return NewImage(buf)
}

func assertSize(buf []byte, width, height int) bool {
	size, err := NewImage(buf).Size()
	if err != nil {
		return false
	}
	if size.Width != 220 || size.Height != 300 {
		return false
	}
	return true
}
