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
		{"test.jpg", "jpeg", 0, false, false, "srgb"},
		{"test_icc_prophoto.jpg", "jpeg", 0, false, true, "srgb"},
		{"test.png", "png", 0, true, false, "srgb"},
		{"test.webp", "webp", 0, false, false, "srgb"},
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
			t.Fatalf("Unexpected image alpha: %t != %t", metadata.Alpha, file.alpha)
		}
		if metadata.Profile != file.profile {
			t.Fatalf("Unexpected image profile: %t != %t", metadata.Profile, file.profile)
		}
		if metadata.Space != file.space {
			t.Fatalf("Unexpected image profile: %t != %t", metadata.Profile, file.profile)
		}
	}
}

func TestImageInterpretation(t *testing.T) {
	files := []struct {
		name           string
		interpretation Interpretation
	}{
		{"test.jpg", InterpretationSRGB},
		{"test.png", InterpretationSRGB},
		{"test.webp", InterpretationSRGB},
	}

	for _, file := range files {
		interpretation, err := ImageInterpretation(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}
		if interpretation != file.interpretation {
			t.Fatalf("Unexpected image interpretation")
		}
	}
}

func TestEXIF(t *testing.T) {
	files := []struct {
		name           string
		make string
		model string
		orientation int
		software string
		datetime string
	}{
		{"test.jpg", "", "", 0, "", ""},
		{"exif/Landscape_1.jpg", "", "", 1, "", ""},
		{"test_exif.jpg", "Jolla", "Jolla", 1, "", "2014:09:21 16:00:56"},
		{"test_exif_canon.jpg", "Canon", "Canon EOS 40D", 1, "GIMP 2.4.5", "2008:07:31 10:38:11"},
	}

	for _, file := range files {
		metadata, err := Metadata(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}
		if metadata.EXIF.Make != file.make {
			t.Fatalf("Unexpected image exif make: %s != %s", metadata.EXIF.Make, file.make)
		}
		if metadata.EXIF.Model != file.model {
			t.Fatalf("Unexpected image exif model: %s != %s", metadata.EXIF.Model, file.model)
		}
		if metadata.EXIF.Orientation != file.orientation {
			t.Fatalf("Unexpected image exif orientation: %d != %d", metadata.EXIF.Orientation, file.orientation)
		}
		if metadata.EXIF.Software != file.software {
			t.Fatalf("Unexpected image exif software: %s != %s", metadata.EXIF.Software, file.software)
		}
		if metadata.EXIF.Datetime != file.datetime {
			t.Fatalf("Unexpected image exif datetime: %s != %s", metadata.EXIF.Datetime, file.datetime)
		}
	}
}

func TestColourspaceIsSupported(t *testing.T) {
	files := []struct {
		name string
	}{
		{"test.jpg"},
		{"test.png"},
		{"test.webp"},
	}

	for _, file := range files {
		supported, err := ColourspaceIsSupported(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}
		if supported != true {
			t.Fatalf("Unsupported image colourspace")
		}
	}

	supported, err := initImage("test.jpg").ColourspaceIsSupported()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	if supported != true {
		t.Errorf("Non-supported colourspace")
	}
}

func readFile(file string) []byte {
	data, _ := os.Open(path.Join("testdata", file))
	buf, _ := ioutil.ReadAll(data)
	return buf
}
