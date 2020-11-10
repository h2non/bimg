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
	}

	for _, file := range files {
		t.Run(file.name, func(t *testing.T) {
			img, _ := os.Open(path.Join("testdata", file.name))
			buf, _ := ioutil.ReadAll(img)
			defer img.Close()

			if VipsIsTypeSupported(file.expected) {
				value := DetermineImageType(buf)
				if value != file.expected {
					t.Fatalf("Image type is not valid: wanted %s, got: %s", ImageTypes[file.expected], ImageTypes[value])
				}
			}
		})
	}
}

func TestDeterminateImageTypeName(t *testing.T) {
	files := []struct {
		name      string
		expected  string
		condition bool
	}{
		{"test.jpg", "jpeg", true},
		{"test.png", "png", true},
		{"test.webp", "webp", true},
		{"test.gif", "gif", true},
		{"test.pdf", "pdf", true},
		{"test.svg", "svg", true},
		// {"test.jp2", "magick"},
		{"test.heic", "heif", vipsVersionMin(8, 8)},
	}

	for _, file := range files {
		t.Run(file.name, func(t *testing.T) {
			if !file.condition {
				t.Skip("condition not met")
			}

			img, _ := os.Open(path.Join("testdata", file.name))
			buf, _ := ioutil.ReadAll(img)
			defer img.Close()

			value := DetermineImageTypeName(buf)
			if value != file.expected {
				t.Fatalf("Image type is not valid: %s != %s, got: %s", file.name, file.expected, value)
			}
		})

	}
}

func TestIsTypeSupported(t *testing.T) {
	types := []struct {
		name      ImageType
		supported bool
	}{
		{JPEG, true},
		{PNG, true},
		{WEBP, true},
		{GIF, true},
		{PDF, true},
		{HEIF, vipsVersionMin(8, 8)},
	}

	for _, typ := range types {
		t.Run(ImageTypes[typ.name], func(t *testing.T) {
			if IsTypeSupported(typ.name) != typ.supported {
				t.Fatalf("Image type support is not as expected")
			}
		})
	}
}

func TestIsTypeNameSupported(t *testing.T) {
	types := []struct {
		name      string
		expected  bool
		condition bool
	}{
		{"jpeg", true, true},
		{"png", true, true},
		{"webp", true, true},
		{"gif", true, true},
		{"pdf", true, true},
		{"heif", true, vipsVersionMin(8, 8)},
	}

	for _, n := range types {
		t.Run(n.name, func(t *testing.T) {
			if !n.condition {
				t.Skip("condition not met")
			}
			if IsTypeNameSupported(n.name) != n.expected {
				t.Fatalf("Image type %s is not valid", n.name)
			}
		})
	}
}

func TestIsTypeSupportedSave(t *testing.T) {
	types := []struct {
		name ImageType
	}{
		{JPEG}, {PNG}, {WEBP},
	}
	if vipsVersionMin(8, 5) {
		types = append(types, struct{ name ImageType }{TIFF})
	}
	if vipsVersionMin(8, 8) {
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
		{"gif", vipsVersionMin(8, 0)},
		{"pdf", false},
		{"tiff", vipsVersionMin(8, 5)},
		{"heif", vipsVersionMin(8, 8)},
	}

	for _, n := range types {
		t.Run(n.name, func(t *testing.T) {
			if IsTypeNameSupportedSave(n.name) != n.expected {
				t.Fatalf("Image type is not valid (expected = %t)", n.expected)
			}
		})
	}
}
