package bimg

import (
	"testing"
)

func TestRead(t *testing.T) {
	buf, err := Read("fixtures/test.jpg")

	if err != nil {
		t.Errorf("Cannot read the image: %#v", err)
	}

	if len(buf) == 0 {
		t.Fatal("Empty buffer")
	}

	if DetermineImageType(buf) != JPEG {
		t.Fatal("Image is not jpeg")
	}
}
