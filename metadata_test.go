package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSize(t *testing.T) {
	files := []struct {
		name   string
		width  int
		height int
	}{
		{"test.jpg", 1680, 1050},
		{"test.png", 400, 300},
		{"test.webp", 550, 368},
	}

	for _, file := range files {
		size, err := Size(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %#v", err)
		}

		if size.Width != file.width || size.Height != file.height {
			t.Fatalf("Unexpected image size: %dx%d", size.Width, size.Height)
		}
	}
}

func TestMetadata(t *testing.T) {
	files := []struct {
		name   string
		format string
	}{
		{"test.jpg", "jpeg"},
		{"test.png", "png"},
		{"test.webp", "webp"},
	}

	for _, file := range files {
		size, err := Metadata(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %#v", err)
		}

		if size.Type != file.format {
			t.Fatalf("Unexpected image format: %s", file.format)
		}
	}
}

func readFile(file string) []byte {
	data, _ := os.Open(path.Join("fixtures", file))
	buf, _ := ioutil.ReadAll(data)
	return buf
}
