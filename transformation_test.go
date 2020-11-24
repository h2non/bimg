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
				Write(fmt.Sprintf("testdata/transformation_resize_%s_out.jpeg", name), out)
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
			Write("testdata/transformation_resize_upscale_out.jpeg", out)
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
			Background: Color{R: 255, G: 0, B: 255},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 1 {
			t.Fatalf("image should still have one channel")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			Write("testdata/transformation_embed_bw_grey_out.png", out)
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
			Background: Color{R: 255, G: 0, B: 0},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 2 {
			t.Fatalf("image should still have two channels")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			Write("testdata/transformation_embed_bwa_grey_out.png", out)
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
		if err := imageTrans.Flatten(Color{R: 255, G: 0, B: 0}); err != nil {
			t.Fatalf("flatten returned unexpected error: %v", err)
		}
		if imageTrans.Metadata().Channels != 1 {
			t.Fatalf("image should still have just one channel")
		}
		if out, err := imageTrans.Save(SaveOptions{}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			Write("testdata/transformation_flatten_bwa_out.png", out)
		}
	})
}
