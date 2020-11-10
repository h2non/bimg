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
