package bimg

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestResize(t *testing.T) {
	options := Options{Width: 800, Height: 600, Crop: false, Rotate: 270}
	img, err := os.Open("fixtures/space.jpg")
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

	debug("Image %s", http.DetectContentType(newImg))

	if http.DetectContentType(newImg) != "image/jpeg" {
		t.Fatal("Image is not jpeg")
	}

	err = ioutil.WriteFile("fixtures/test.jpg", newImg, 0644)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
}
