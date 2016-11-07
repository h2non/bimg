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
	}

	for _, file := range files {
		img, _ := os.Open(path.Join("fixtures", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		if DetermineImageType(buf) != file.expected {
			t.Fatal("Image type is not valid")
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
	}

	for _, file := range files {
		img, _ := os.Open(path.Join("fixtures", file.name))
		buf, _ := ioutil.ReadAll(img)
		defer img.Close()

		if DetermineImageTypeName(buf) != file.expected {
			t.Fatal("Image type is not valid")
		}
	}
}

func TestIsTypeSupported(t *testing.T) {
	types := []struct {
		name ImageType
	}{
		{JPEG}, {PNG}, {WEBP}, {GIF}, {PDF},
	}

	for _, n := range types {
		if IsTypeSupported(n.name) == false {
			t.Fatal("Image type is not valid")
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
	}

	for _, n := range types {
		if IsTypeNameSupported(n.name) != n.expected {
			t.Fatal("Image type is not valid")
		}
	}
}
