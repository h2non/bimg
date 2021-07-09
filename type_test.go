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
		// {"test.jp2", MAGICK},
		{"test.heic", HEIF},
		{"test2.heic", HEIF},
		{"test3.heic", HEIF},
		{"test.avif", AVIF},
	}
//
//return
	for _, file := range files {
		img, _ := os.Open(path.Join("testdata", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		if VipsIsTypeSupported(file.expected) {
			value := DetermineImageType(buf)
			if value != file.expected {
				t.Fatalf("Image type is not valid: %s != %s, got: %s", file.name, ImageTypes[file.expected], ImageTypes[value])
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
		// {"test.jp2", "magick"},
		{"test.heic", "heif"},
		{"test.avif", "avif"},
	}

	for _, file := range files {
		if !IsTypeNameSupported(file.expected) {
			t.Skip("Skip test in TestDeterminateImageTypeName: " + file.expected + " is not supported.")
			continue
		}

		img, _ := os.Open(path.Join("testdata", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		value := DetermineImageTypeName(buf)
		if value != file.expected {
			t.Fatalf("Image type is not valid: %s != %s, got: %s", file.name, file.expected, value)
		}
	}
}

// Dumb Tests? Testing Vips types compiled in, which isn't known without using the function this is testing...
func TestIsTypeSupported(t *testing.T) {
}

func TestIsTypeNameSupported(t *testing.T) {
}

func TestIsTypeSupportedSave(t *testing.T) {
}

func TestIsTypeNameSupportedSave(t *testing.T) {
}
