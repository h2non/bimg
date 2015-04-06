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

	err = ioutil.WriteFile("fixtures/test_out.jpg", newImg, 0644)
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

	err = ioutil.WriteFile("fixtures/test_out.png", newImg, 0644)
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

	err = ioutil.WriteFile("fixtures/transparent_out.png", newImg, 0644)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}
