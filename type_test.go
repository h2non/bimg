package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestDeterminateImageType(t *testing.T) {
	files := []struct {
		name       string
		expected   ImageType
		shouldTest bool
	}{
		{"test.jpg", JPEG, true},
		{"test.png", PNG, true},
		{"test.webp", WEBP, true},
		{"test.gif", GIF, true},
		{"test.pdf", PDF, true},
		{"test.svg", SVG, true},
		{"test.jp2", JP2K, vipsVersionMin(8, 11)},
		{"test.heic", HEIF, true},
		{"test2.heic", HEIF, true},
		{"test3.heic", HEIF, true},
		{"test.avif", AVIF, true},
		{"test.bmp", MAGICK, true},
	}

	for _, file := range files {
		if !file.shouldTest {
			t.Skip("condition not met")
		}
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
		{"test.jp2", "jp2k", vipsVersionMin(8, 11)},
		{"test.heic", "heif", vipsVersionMin(8, 8)},
		{"test.avif", "avif", vipsVersionMin(8, 9)},
		{"test.bmp", "magick", true},
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
		{AVIF, vipsVersionMin(8, 9)},
		{JP2K, vipsVersionMin(8, 11)},
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
		{"avif", true, vipsVersionMin(8, 9)},
		{"jp2k", true, vipsVersionMin(8, 11)},
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
	types := []ImageType{
		JPEG, PNG, WEBP,
	}
	if vipsVersionMin(8, 5) {
		types = append(types, TIFF)
	}
	if vipsVersionMin(8, 8) {
		types = append(types, HEIF)
	}
	if vipsVersionMin(8, 9) {
		types = append(types, AVIF)
	}
	if vipsVersionMin(8, 11) {
		types = append(types, JP2K)
	}

	for _, tt := range types {
		if IsTypeSupportedSave(tt) == false {
			t.Fatalf("Image type %s is not valid", ImageTypes[tt])
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
		{"avif", vipsVersionMin(8, 9)},
		{"jp2k", vipsVersionMin(8, 11)},
	}

	for _, n := range types {
		t.Run(n.name, func(t *testing.T) {
			if IsTypeNameSupportedSave(n.name) != n.expected {
				t.Fatalf("Image type is not valid (expected = %t)", n.expected)
			}
		})
	}
}
