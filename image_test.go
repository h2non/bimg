package bimg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestImageResize(t *testing.T) {
	image := readImage()
	_, err := image.Resize(300, 240)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
}

func TestImageCrop(t *testing.T) {
	image := readImage()
	_, err := image.Crop(800, 600)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
}

func TestImageRotate(t *testing.T) {
	image := readImage()
	_, err := image.Rotate(D90)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
}

func readImage() *Image {
	data, _ := os.Open("fixtures/test.jpg")
	buf, _ := ioutil.ReadAll(data)
	return NewImage(buf)
}
