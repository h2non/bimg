package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestDeterminateImageType(t *testing.T) {
	files := []struct {
		name     string
		expected ImageType
	}{
		{"test.jpg", JPEG},
		{"test.png", PNG},
		{"test.webp", WEBP},
		{"test.gif", GIF},
		{"test.pdf", PDF},
		{"test.svg", SVG},
		{"test.jp2", MAGICK},
		{"test.heic", HEIF},
		{"test2.heic", HEIF},
	}

	for _, file := range files {
		img, _ := os.Open(path.Join("testdata", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		if VipsIsTypeSupported(file.expected) {
			if DetermineImageType(buf) != file.expected {
				t.Fatalf("Image type is not valid: %s != %s", file.name, ImageTypes[file.expected])
			}
		}
	}
}

func TestDeterminateImageTypeName(t *testing.T) {
	files := []struct {
		name     string
		expected string
	}{
		{"test.jpg", "jpeg"},
		{"test.png", "png"},
		{"test.webp", "webp"},
		{"test.gif", "gif"},
		{"test.pdf", "pdf"},
		{"test.svg", "svg"},
		{"test.jp2", "magick"},
		{"test.heic", "heif"},
	}

	for _, file := range files {
		img, _ := os.Open(path.Join("testdata", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		if DetermineImageTypeName(buf) != file.expected {
			t.Fatalf("Image type is not valid: %s != %s", file.name, file.expected)
		}
	}
}

func TestIsTypeSupported(t *testing.T) {
	types := []struct {
		name ImageType
	}{
		{JPEG}, {PNG}, {WEBP}, {GIF}, {PDF}, {HEIF},
	}

	for _, n := range types {
		if IsTypeSupported(n.name) == false {
			t.Fatalf("Image type %s is not valid", ImageTypes[n.name])
		}
	}
}

func TestIsTypeNameSupported(t *testing.T) {
	types := []struct {
		name     string
		expected bool
	}{
		{"jpeg", true},
		{"png", true},
		{"webp", true},
		{"gif", true},
		{"pdf", true},
		{"heif", true},
	}

	for _, n := range types {
		if IsTypeNameSupported(n.name) != n.expected {
			t.Fatalf("Image type %s is not valid", n.name)
		}
	}
}

func TestIsTypeSupportedSave(t *testing.T) {
	types := []struct {
		name ImageType
	}{
		{JPEG}, {PNG}, {WEBP},
	}
	if VipsVersion >= "8.5.0" {
		types = append(types, struct{ name ImageType }{TIFF})
	}
	if VipsVersion >= "8.8.0" {
		types = append(types, struct{ name ImageType }{HEIF})
	}

	for _, n := range types {
		if IsTypeSupportedSave(n.name) == false {
			t.Fatalf("Image type %s is not valid", ImageTypes[n.name])
		}
	}
}

func TestIsTypeNameSupportedSave(t *testing.T) {
	types := []struct {
		name     string
		expected bool
	}{
		{"jpeg", true},
		{"png", true},
		{"webp", true},
		{"gif", false},
		{"pdf", false},
		{"tiff", VipsVersion >= "8.5.0"},
		{"heif", VipsVersion >= "8.8.0"},
	}

	for _, n := range types {
		if IsTypeNameSupportedSave(n.name) != n.expected {
			t.Fatalf("Image type %s is not valid", n.name)
		}
	}
}
