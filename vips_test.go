package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestVipsRead(t *testing.T) {
	files := []struct {
		name     string
		expected ImageType
	}{
		{"test.jpg", JPEG},
		{"test.png", PNG},
		{"test.webp", WEBP},
	}

	for _, file := range files {
		img, _ := os.Open(path.Join("fixtures", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		image, imageType, _ := vipsRead(buf)
		if image == nil {
			t.Fatal("Empty image")
		}
		if imageType != file.expected {
			t.Fatal("Empty image")
		}
	}
}

func TestVipsSave(t *testing.T) {
	img, _ := os.Open(path.Join("fixtures", "test.jpg"))
	buf, _ := ioutil.ReadAll(img)
	defer img.Close()

	image, _, _ := vipsRead(buf)
	if image == nil {
		t.Fatal("Empty image")
	}

	options := vipsSaveOptions{Quality: 95, Type: JPEG}

	buf, err := vipsSave(image, options)
	if err != nil {
		t.Fatal("Cannot save the image")
	}
	if len(buf) == 0 {
		t.Fatal("Empty image")
	}
}
