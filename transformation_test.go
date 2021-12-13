package bimg

import (
	"fmt"
	"testing"
)

func TestImageTransformation_Resize(t *testing.T) {
	tests := []struct {
		mode     ResizeMode
		width    int
		height   int
		expected ImageSize
	}{
		{
			mode:     ResizeModeFit,
			width:    400,
			height:   300,
			expected: ImageSize{400, 250},
		},
		{
			mode:     ResizeModeFit,
			width:    300,
			height:   400,
			expected: ImageSize{300, 187},
		},
		{
			mode:     ResizeModeFit,
			width:    800,
			height:   400,
			expected: ImageSize{640, 400},
		},
		{
			mode:     ResizeModeFitUp,
			width:    400,
			height:   300,
			expected: ImageSize{480, 300},
		},
		{
			mode:     ResizeModeFitUp,
			width:    300,
			height:   400,
			expected: ImageSize{640, 400},
		},
		{
			mode:     ResizeModeFitUp,
			width:    800,
			height:   400,
			expected: ImageSize{800, 500},
		},
		{
			mode:     ResizeModeForce,
			width:    400,
			height:   300,
			expected: ImageSize{400, 300},
		},
		{
			mode:     ResizeModeForce,
			width:    300,
			height:   400,
			expected: ImageSize{300, 400},
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s_%d_%d", test.mode, test.width, test.height)
		t.Run(name, func(t *testing.T) {
			imageTrans, err := NewImageTransformation(readImage("test.jpg"))
			if err != nil {
				t.Fatalf("cannot load image: %v", err)
			}
			if err := imageTrans.Resize(ResizeOptions{Width: test.width, Height: test.height, Mode: test.mode}); err != nil {
				t.Fatalf("cannot resize image: %v", err)
			}
			size := imageTrans.Size()
			if size != test.expected {
				t.Errorf("unexpected size. wanted %#v, got %#v", test.expected, size)
			}
			if out, err := imageTrans.Save(SaveOptions{Type: JPEG}); err != nil {
				t.Errorf("cannot save image: %v", err)
			} else {
				_ = Write(fmt.Sprintf("testdata/transformation_resize_%s_out.jpeg", name), out)
			}
		})
	}

	t.Run("upscale", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test.webp"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if err := imageTrans.Resize(ResizeOptions{Width: 2000, Height: 2000}); err != nil {
			t.Fatalf("cannot resize image: %v", err)
		}
		size := imageTrans.Size()
		expected := ImageSize{2000, 1338}
		if size != expected {
			t.Errorf("unexpected size. wanted %#v, got %#v", expected, size)
		}
		if out, err := imageTrans.Save(SaveOptions{Type: JPEG}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_resize_upscale_out.jpeg", out)
		}
	})
}

func TestImageTransformation_Embed(t *testing.T) {
	t.Run("B/W on grey", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test_bw.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if imageTrans.Metadata().Channels != 1 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := imageTrans.Embed(EmbedOptions{
			Width:      200,
			Height:     200,
			Extend:     ExtendBackground,
			Background: Color{R: 50, G: 50, B: 50},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 1 {
			t.Fatalf("image should still have one channel")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_embed_bw_grey_out.png", out)
		}
	})

	t.Run("B/W with alpha on grey", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if imageTrans.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := imageTrans.Embed(EmbedOptions{
			Width:      200,
			Height:     200,
			Extend:     ExtendBackground,
			Background: Color{R: 120, G: 120, B: 120},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 2 {
			t.Fatalf("image should still have two channels")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_embed_bwa_grey_out.png", out)
		}
	})

	t.Run("B/W with alpha on red", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if imageTrans.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := imageTrans.Embed(EmbedOptions{
			Width:      200,
			Height:     200,
			Extend:     ExtendBackground,
			Background: Color{R: 255, G: 0, B: 0},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 4 {
			t.Fatalf("image should have four channels now")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_embed_bwa_red_out.png", out)
		}
	})

	t.Run("CMYK", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test_cmyk.jpeg"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if imageTrans.Metadata().Channels != 4 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if imageTrans.Metadata().Interpretation != InterpretationCMYK {
			t.Fatalf("source image has unexpected interpretation (should be CMYK)")
		}
		if err := imageTrans.Embed(EmbedOptions{
			Width:      1000,
			Height:     1000,
			Extend:     ExtendBackground,
			Background: ColorWithAlpha{Color: Color{R: 255, G: 0, B: 0}, A: 100},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 4 {
			t.Fatalf("image should still have four channels now")
		}
		if imageTrans.Metadata().Interpretation != InterpretationSRGB {
			t.Fatalf("image should be sRGB now")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_embed_cmyk_on_alpha_out.jpeg", out)
		}
	})
}

func TestImageTransformation_Flatten(t *testing.T) {
	t.Run("B/W with alpha", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if imageTrans.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := imageTrans.Flatten(Color{R: 255, G: 255, B: 255}); err != nil {
			t.Fatalf("flatten returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 1 {
			t.Errorf("image should have just one channel now (no alpha)")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_flatten_bwa_out.png", out)
		}
	})

	t.Run("B/W with alpha on red", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if imageTrans.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := imageTrans.Flatten(Color{R: 255, G: 0, B: 0}); err != nil {
			t.Fatalf("flatten returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 3 {
			t.Errorf("image should have three channels now (RGB without alpha)")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_flatten_bwa_on_red_out.png", out)
		}
	})
}

func TestImageTransformation_Save(t *testing.T) {
	t.Run("save bitmap", func(t *testing.T) {
		imageTrans, err := NewImageTransformation(readImage("test.bmp"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if out, err := imageTrans.Save(SaveOptions{MagickFormat: "bmp"}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_save_bmp_out.bmp", out)
		}
	})
}

func TestFormatSupport(t *testing.T) {
	t.Run("jpeg2000", func(t *testing.T) {
		if !vipsVersionMin(8, 11) {
			t.Skip("requires libvips >= 8.11")
		}

		t.Run("can load", func(t *testing.T) {
			buf, err := Read("testdata/test.jp2")
			if err != nil {
				t.Fatalf("cannot load image: %v", err)
			}
			img, err := NewImageTransformation(buf)
			if err != nil {
				t.Fatalf("cannot load image: %v", err)
			}
			metadata := img.Metadata()
			if metadata.Type != "jp2k" {
				t.Errorf("unexpected image type %q", metadata.Type)
			}
			if metadata.Size.Width != 1680 {
				t.Errorf("unexpected width: %d", metadata.Size.Width)
			}
			if metadata.Size.Height != 1050 {
				t.Errorf("unexpected height: %d", metadata.Size.Height)
			}
		})

		t.Run("can save", func(t *testing.T) {
			buf, err := Read("testdata/test.jpg")
			if err != nil {
				t.Fatalf("cannot load image: %v", err)
			}
			img, err := NewImageTransformation(buf)
			if err != nil {
				t.Fatalf("cannot load image: %v", err)
			}
			err = img.Resize(ResizeOptions{Width: 100})
			if err != nil {
				t.Fatalf("cannot resize image: %v", err)
			}
			outBuf, err := img.Save(SaveOptions{Type: JP2K})
			if err != nil {
				t.Fatalf("cannot save image: %v", err)
			}
			outType := vipsImageType(outBuf)
			if outType != JP2K {
				t.Errorf("output has unexpected type: %d", outType)
			}
		})
	})
}
