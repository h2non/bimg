package bimg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestImageResize(t *testing.T) {
	data, _ := os.Open("fixtures/test.jpg")
	buf, err := ioutil.ReadAll(data)

	image := NewImage(buf)

	buf, err = image.Resize(300, 240)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
}
