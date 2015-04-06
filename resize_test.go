package bimg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestResize(t *testing.T) {
	options := Options{Width: 800, Height: 600, Crop: false, Rotate: 270}
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

	if DetermineImageType(newImg) != JPEG {
		t.Fatal("Image is not jpeg")
	}

	err = ioutil.WriteFile("result.jpg", newImg, 0644)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}

func TestConvert(t *testing.T) {
	options := Options{Width: 640, Height: 480, Crop: true, Type: PNG}
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

	err = ioutil.WriteFile("result.png", newImg, 0644)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}
