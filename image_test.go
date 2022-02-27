package bimg

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func testfile(fileName string) string {
	return filepath.Join("testdata", fileName)
}

func writeImage(image *Image, fileName string) error {
	b, err := image.Save(SaveOptions{})
	if err != nil {
		return err
	}
	return os.WriteFile(testfile(fileName), b, 0644)
}

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
			expected: ImageSize{300, 188},
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
			imageTrans, err := NewImageFromFile(testfile("test.jpg"))
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
				// TODO comparison
				_ = Write(fmt.Sprintf("testdata/transformation_resize_%s_out.jpeg", name), out)
			}
		})
	}

	t.Run("on a tainted buffer", func(t *testing.T) {
		for _, test := range tests {
			name := fmt.Sprintf("%s_%d_%d", test.mode, test.width, test.height)
			t.Run(name, func(t *testing.T) {
				imageTrans, err := NewImageFromFile(testfile("test.jpg"))
				if err != nil {
					t.Fatalf("cannot load image: %v", err)
				}
				// Apply an operation before the resize tains the buffer, which can (but shouldn't)
				// influence calculations.
				if err := imageTrans.AutoRotate(); err != nil {
					t.Fatalf("cannot autorotate image: %v", err)
				}
				if !imageTrans.bufTainted {
					t.Fatalf("buffer should be tainted now") // otherwise the test is pointless
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
					_ = Write(fmt.Sprintf("testdata/transformation_resize_tainted_%s_out.jpeg", name), out)
				}
			})
		}
	})

	t.Run("upscale", func(t *testing.T) {
		image, err := NewImageFromFile(testfile("test.webp"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if err := image.Resize(ResizeOptions{Width: 2000, Height: 2000, Mode: ResizeModeFit}); err != nil {
			t.Fatalf("cannot resize image: %v", err)
		}
		size := image.Size()
		expected := ImageSize{2000, 1338}
		if size != expected {
			t.Errorf("unexpected size. wanted %#v, got %#v", expected, size)
		}
		if out, err := image.Save(SaveOptions{Type: JPEG}); err != nil {
			t.Errorf("cannot save image: %v", err)
		} else {
			_ = Write("testdata/transformation_resize_upscale_out.jpeg", out)
		}
	})

	t.Run("rounding precision", func(t *testing.T) {
		// see https://github.com/h2non/bimg/issues/99
		img := image.NewGray16(image.Rect(0, 0, 1920, 1080))
		input := &bytes.Buffer{}
		_ = jpeg.Encode(input, img, nil)

		desiredWidth := 300

		imageTrans, err := NewImageFromBuffer(input.Bytes())
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if err := imageTrans.Resize(ResizeOptions{Width: desiredWidth}); err != nil {
			t.Fatalf("cannot resize: %v", err)
		}
		size := imageTrans.Size()

		if size.Width != desiredWidth {
			t.Fatalf("Invalid width: %d", size.Width)
		}
	})

	t.Run("handles corrupt image", func(t *testing.T) {
		image, err := NewImageFromFile(testfile("corrupt.jpg"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if err := image.Resize(ResizeOptions{Width: 800, Height: 600}); err != nil {
			t.Fatalf("cannot resize image: %v", err)
		}
		size := image.Size()

		if size.Width != 800 || size.Height != 600 {
			t.Errorf("invalid image size: %d x %d", size.Width, size.Height)
		}
	})
}

func TestImage_Rotate(t *testing.T) {
	image, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	initialSize := image.Size()
	if err := image.Rotate(270); err != nil {
		t.Fatalf("cannot rotate image: %v", err)
	}
	newSize := image.Size()

	if newSize.Width != initialSize.Height || newSize.Height != initialSize.Width {
		t.Errorf("invalid image size: %d x %d", newSize.Width, newSize.Height)
	}
}

func TestImage_Blur(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.Resize(ResizeOptions{Width: 800, Height: 600}); err != nil {
		t.Fatalf("cannot resize image: %v", err)
	}
	if err := img.Blur(GaussianBlurOptions{Sigma: 5}); err != nil {
		t.Fatalf("cannot apply blur: %v", err)
	}
	_ = writeImage(img, "test_gaussian_out.jpg")
}

func TestImage_Sharpen(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.Resize(ResizeOptions{Width: 800, Height: 600}); err != nil {
		t.Fatalf("cannot resize image: %v", err)
	}
	if err := img.Sharpen(SharpenOptions{
		Radius: 1, X1: 1.5, Y2: 20, Y3: 50, M1: 1, M2: 2,
	}); err != nil {
		t.Fatalf("cannot sharpen image: %v", err)
	}
	_ = writeImage(img, "test_sharpen_out.jpg")
}

func TestImageTransformation_Embed(t *testing.T) {
	t.Run("extend image on white", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_issue.jpg"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if err := img.Resize(ResizeOptions{Width: 400, Height: 600, Mode: ResizeModeFit}); err != nil {
			t.Fatalf("cannot resize image: %v", err)
		}
		if err := img.Embed(EmbedOptions{Width: 400, Height: 600, Extend: ExtendWhite, Background: Color{255, 20, 10}}); err != nil {
			t.Fatalf("cannot embed image: %v", err)
		}

		_ = writeImage(img, "test_extend_white_out.jpg")
	})

	t.Run("extend image on background", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_issue.jpg"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if err := img.Resize(ResizeOptions{Width: 400, Height: 600, Mode: ResizeModeFit}); err != nil {
			t.Fatalf("cannot resize image: %v", err)
		}
		if err := img.Embed(EmbedOptions{Width: 400, Height: 600, Extend: ExtendBackground, Background: Color{255, 20, 10}}); err != nil {
			t.Fatalf("cannot embed image: %v", err)
		}

		_ = writeImage(img, "test_extend_background_out.jpg")
	})

	t.Run("B/W on grey", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_bw.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if img.Metadata().Channels != 1 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := img.Embed(EmbedOptions{
			Width:      200,
			Height:     200,
			Extend:     ExtendBackground,
			Background: Color{R: 50, G: 50, B: 50},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if img.Metadata().Channels != 1 {
			t.Fatalf("image should still have one channel")
		}
		_ = writeImage(img, "transformation_embed_bw_grey_out.png")
	})

	t.Run("B/W with alpha on grey", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if img.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := img.Embed(EmbedOptions{
			Width:      200,
			Height:     200,
			Extend:     ExtendBackground,
			Background: Color{R: 120, G: 120, B: 120},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if img.Metadata().Channels != 2 {
			t.Fatalf("image should still have two channels")
		}
		_ = writeImage(img, "transformation_embed_bwa_grey_out.png")
	})

	t.Run("B/W with alpha on red", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if img.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := img.Embed(EmbedOptions{
			Width:      200,
			Height:     200,
			Extend:     ExtendBackground,
			Background: Color{R: 255, G: 0, B: 0},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if img.Metadata().Channels != 4 {
			t.Fatalf("image should have four channels now")
		}
		_ = writeImage(img, "transformation_embed_bwa_red_out.png")
	})

	t.Run("CMYK", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_cmyk.jpeg"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if img.Metadata().Channels != 4 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if img.Metadata().Interpretation != InterpretationCMYK {
			t.Fatalf("source image has unexpected interpretation (should be CMYK)")
		}
		if err := img.Embed(EmbedOptions{
			Width:      1000,
			Height:     1000,
			Extend:     ExtendBackground,
			Background: ColorWithAlpha{Color: Color{R: 255, G: 0, B: 0}, A: 100},
		}); err != nil {
			t.Fatalf("embed returned unexpected error: %v", err)
		}
		if img.Metadata().Channels != 4 {
			t.Fatalf("image should still have four channels now")
		}
		if img.Metadata().Interpretation != InterpretationSRGB {
			t.Fatalf("image should be sRGB now")
		}
		_ = writeImage(img, "transformation_embed_cmyk_on_alpha_out.jpeg")
	})
}

func TestImageTransformation_Flatten(t *testing.T) {
	t.Run("B/W with alpha", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if img.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := img.Flatten(Color{R: 255, G: 255, B: 255}); err != nil {
			t.Fatalf("flatten returned unexpected error: %v", err)
		}
		if img.Metadata().Channels != 1 {
			t.Errorf("image should have just one channel now (no alpha)")
		}
		_ = writeImage(img, "transformation_flatten_bwa_out.png")
	})

	t.Run("B/W with alpha on red", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test_bwa.png"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if img.Metadata().Channels != 2 {
			t.Fatalf("source image has unexpected number of channels")
		}
		if err := img.Flatten(Color{R: 255, G: 0, B: 0}); err != nil {
			t.Fatalf("flatten returned unexpected error: %v", err)
		}
		if img.Metadata().Channels != 3 {
			t.Errorf("image should have three channels now (RGB without alpha)")
		}
		_ = writeImage(img, "transformation_flatten_bwa_on_red_out.png")
	})
}

func TestImage_FlipHorizontal(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.FlipHorizontal(); err != nil {
		t.Fatalf("cannot flip image: %v", err)
	}
	_ = writeImage(img, "test_flip_horizontal_out.jpg")
}

func TestImage_FlipVertical(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.FlipVertical(); err != nil {
		t.Fatalf("cannot flip image: %v", err)
	}
	_ = writeImage(img, "test_flip_vertical_out.jpg")
}

func TestImage_Trim(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.png"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.Trim(TrimOptions{Background: ColorBlack, Threshold: 0}); err != nil {
		t.Fatalf("cannot trim image: %v", err)
	}
	size := img.Size()
	if size.Width != 400 || size.Height != 257 {
		t.Errorf("unexpected size: %d x %d", size.Width, size.Height)
	}
	_ = writeImage(img, "test_trim_out.png")
}

func TestImage_Extract(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}

	if err := img.Extract(ExtractOptions{
		Left:   100,
		Top:    100,
		Width:  200,
		Height: 200,
	}); err != nil {
		t.Fatalf("cannot extract area: %v", err)
	}

	_ = writeImage(img, "test_extract_out.jpg")
}

func TestImage_Crop(t *testing.T) {
	tests := []struct {
		name    string
		gravity Gravity
	}{
		{"center", GravityCentre},
		{"north", GravityNorth},
		{"south", GravitySouth},
		{"west", GravityWest},
		{"east", GravityEast},
		{"smart", GravitySmart},
	}

	sourceImg, err := NewImageFromFile(testfile("test_issue.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := sourceImg.Clone()
			if err := img.Crop(CropOptions{
				Width:   200,
				Height:  200,
				Gravity: tt.gravity,
			}); err != nil {
				t.Fatalf("cannot crop image: %v", err)
			}
			_ = writeImage(img, fmt.Sprintf("test_crop_%s_out.jpg", tt.name))
		})
	}
}

func TestImage_Gamma(t *testing.T) {
	img, err := NewImageFromFile(testfile("northern_cardinal_bird.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.Gamma(0.1); err != nil {
		t.Fatalf("cannot apply gamma: %v", err)
	}
	_ = writeImage(img, "test_gamma_out.jpg")
}

func TestImage_ChangeColorspace(t *testing.T) {
	img, err := NewImageFromFile(testfile("test.jpg"))
	if err != nil {
		t.Fatalf("cannot load image: %v", err)
	}
	if err := img.ChangeColorspace(InterpretationGREY16); err != nil {
		t.Fatalf("cannot apply colorspace: %v", err)
	}
	_ = writeImage(img, "test_grey16_out.jpg")
}

func TestImageTransformation_Save(t *testing.T) {
	t.Run("save bitmap", func(t *testing.T) {
		img, err := NewImageFromFile(testfile("test.bmp"))
		if err != nil {
			t.Fatalf("cannot load image: %v", err)
		}
		if out, err := img.Save(SaveOptions{MagickFormat: "bmp"}); err != nil {
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
			img, err := NewImageFromFile(testfile("test.jp2"))
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
			img, err := NewImageFromFile(testfile("test.jpg"))
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

	t.Run("pdf", func(t *testing.T) {
		t.Run("can load", func(t *testing.T) {
			img, err := NewImageFromFile(testfile("test.pdf"))
			if err != nil {
				t.Fatalf("cannot load the pdf: %v", err)
			}

			size := img.Size()
			if size.Height != 1050 || size.Width != 1680 {
				t.Errorf("unexpected size: %#v", size)
			}
		})

		t.Run("cannot save", func(t *testing.T) {
			img, err := NewImageFromFile(testfile("test.pdf"))
			if err != nil {
				t.Fatalf("cannot load the pdf: %v", err)
			}

			_, err = img.Save(SaveOptions{})
			if err == nil {
				t.Error("saving should not work")
			}
		})
	})

	// TODO add a table test for all expected formats
}
