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
		name        string
		format      string
		orientation int
		alpha       bool
		profile     bool
		space       string
	}{
		{"test.jpg", "jpeg", 0, false, false, "bicubic"},
		{"test.png", "png", 0, true, false, "bicubic"},
		{"test.webp", "webp", 0, false, false, "bicubic"},
	}

	for _, file := range files {
		metadata, err := Metadata(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}

		if metadata.Type != file.format {
			t.Fatalf("Unexpected image format: %s", file.format)
		}
		if metadata.Orientation != file.orientation {
			t.Fatalf("Unexpected image orientation: %d != %d", metadata.Orientation, file.orientation)
		}
		if metadata.Alpha != file.alpha {
			t.Fatalf("Unexpected image alpha: %s != ", metadata.Alpha, file.alpha)
		}
		if metadata.Profile != file.profile {
			t.Fatalf("Unexpected image profile: %s != %s", metadata.Profile, file.profile)
		}
	}
}

func readFile(file string) []byte {
	data, _ := os.Open(path.Join("fixtures", file))
	buf, _ := ioutil.ReadAll(data)
	return buf
}
