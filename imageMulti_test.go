package bimg

import (
	"testing"
)

func TestArrayJoin(t *testing.T) {
	var images [][]byte

	for _, file := range []string{"test.jpg","test_smart_crop.jpg", "vertical.jpg","test.jpg"} {
		buf := initImage(file).Image()
		images = append(images, buf)
	}

	outputOptions := &Options{}

	opts := ArrayJoin{
		Across: 2,
		VSpacing: 128,
		HSpacing: 128,
	}

	buffer, err := NewImageFrom(images, outputOptions).ArrayJoin(opts)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	err = assertSize(buffer, 256, 256)
	if err != nil {
		t.Error(err)
	}

	Write("testdata/test_arrayjoin_out.jpg", buffer)
}

func TestMosaic(t *testing.T) {
	var images [][]byte

	for _, file := range []string{"test.jpg","test_smart_crop.jpg"} {
		buf := initImage(file).Image()
		images = append(images, buf)
	}

	outputOptions := &Options{}

	opts := Mosaic{
	}

	buffer, err := NewImageFrom(images, outputOptions).Mosaic(opts)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	Write("testdata/test_mosaic_out.jpg", buffer)
}

func TestComposite(t *testing.T) {
	var images [][]byte

	for _, file := range []string{"test.jpg", "test_smart_crop.jpg", "transparent.png"} {
		buf := initImage(file).Image()
		images = append(images, buf)
	}

	outputOptions := &Options{
		Type: PNG,
	}

	opts := Composite{
		Mode: []BlendMode{ColorBurn, Overlay},
	}

	buffer, err := NewImageFrom(images, outputOptions).Composite(opts)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	Write("testdata/test_composite_out.png", buffer)
}

func TestComposite2(t *testing.T) {
	var images [][]byte

	for _, file := range []string{"test.jpg", "transparent.png"} {
		buf := initImage(file).Image()
		images = append(images, buf)
	}

	outputOptions := &Options{
		Type: PNG,
	}

	opts := Composite2{
	}

	buffer, err := NewImageFrom(images, outputOptions).Composite2(opts)
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}

	Write("testdata/test_composite2_out.png", buffer)
}