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
		if file.expected == "heif" && VipsMajorVersion <= 8 && VipsMinorVersion < 8 {
			continue
		}
		if file.expected == "avif" && VipsMajorVersion <= 8 && VipsMinorVersion < 9 {
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

func TestIsTypeSupported(t *testing.T) {
	types := []struct {
		name ImageType
	}{
		{JPEG}, {PNG}, {WEBP}, {GIF}, {PDF}, {HEIF}, {AVIF},
	}

	for _, n := range types {
		if n.name == HEIF && VipsMajorVersion <= 8 && VipsMinorVersion < 8 {
			continue
		}
		if n.name == AVIF && VipsMajorVersion <= 8 && VipsMinorVersion < 9 {
			continue
		}
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
		{"avif", true},
	}

	for _, n := range types {
		if n.name == "heif" && VipsMajorVersion <= 8 && VipsMinorVersion < 8 {
			continue
		}
		if n.name == "avif" && VipsMajorVersion <= 8 && VipsMinorVersion < 9 {
			continue
		}
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
	if VipsVersion >= "8.9.0" {
		types = append(types, struct{ name ImageType }{AVIF})
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
		{"avif", VipsVersion >= "8.9.0"},
	}

	for _, n := range types {
		if IsTypeNameSupportedSave(n.name) != n.expected {
			t.Fatalf("Image type %s is not valid", n.name)
		}
	}
}
