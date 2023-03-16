package main

import (
	"fmt"
	"os"

	"github.com/nestor-sk/vimgo"
)

func main() {
	buffer, err := os.ReadFile("./testdata/test.jpg")
	check(err)

	resize("./bin/resize.jpg", buffer)
	rotate("./bin/rotate.jpg", buffer)
	convert("./bin/convert.png", buffer)
	forceResize("./bin/force_resize.jpg", buffer)
	customColorSpace("./bin/custom_color_space.jpg", buffer)
	customOptions("./bin/custom_options.jpg", buffer)
	waterMark("./bin/watermark.jpg", buffer)
	fuentInterface("./bin/fuent_interface.jpg", buffer)

}

func write(name string, buffer []byte) {
	os.WriteFile(name, buffer, 0664)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func resize(name string, buffer []byte) {
	newImage, err := vimgo.NewImage(buffer).Resize(800, 600)
	check(err)

	size, err := vimgo.NewImage(newImage).Size()
	check(err)
	if size.Width == 800 && size.Height == 600 {
		fmt.Println("The image size is valid")
	}
	write(name, newImage)
}

func convert(name string, buffer []byte) {
	newImage, err := vimgo.NewImage(buffer).Convert(vimgo.PNG)
	check(err)

	if vimgo.NewImage(newImage).Type() == "png" {
		fmt.Fprintln(os.Stderr, "The image was converted into png")
	}
	write(name, newImage)
}

func rotate(name string, buffer []byte) {
	newImage, err := vimgo.NewImage(buffer).Rotate(90)
	check(err)
	write(name, newImage)
}

func forceResize(name string, buffer []byte) {
	newImage, err := vimgo.NewImage(buffer).ForceResize(1000, 500)
	check(err)

	size, err := vimgo.Size(newImage)
	check(err)
	if size.Width != 1000 || size.Height != 500 {
		fmt.Fprintln(os.Stderr, "Incorrect image size")
	}
	write(name, newImage)
}

func customColorSpace(name string, buffer []byte) {
	newImage, err := vimgo.NewImage(buffer).Colourspace(vimgo.InterpretationBW)
	check(err)

	colourSpace, err := vimgo.ImageInterpretation(newImage)
	check(err)
	if colourSpace != vimgo.InterpretationBW {
		fmt.Fprintln(os.Stderr, "Invalid colour space")
	}
	write(name, newImage)
}

func customOptions(name string, buffer []byte) {
	//Dee Options struct in options.go to discover all the available fields

	options := vimgo.Options{
		Width:     800,
		Height:    600,
		Crop:      true,
		Quality:   95,
		Rotate:    180,
		Interlace: true,
	}

	newImage, err := vimgo.NewImage(buffer).Process(options)
	check(err)
	write(name, newImage)
}

func waterMark(name string, buffer []byte) {
	watermark := vimgo.Watermark{
		Text:       "Chuck Norris (c) 2315",
		Opacity:    0.25,
		Width:      200,
		DPI:        100,
		Margin:     150,
		Font:       "sans bold 12",
		Background: vimgo.Color{R: 255, G: 255, B: 255},
	}

	newImage, err := vimgo.NewImage(buffer).Watermark(watermark)
	check(err)
	write(name, newImage)
}

func fuentInterface(name string, buffer []byte) {
	image := vimgo.NewImage(buffer)

	// first crop image
	_, err := image.CropByWidth(300)
	check(err)

	// then flip it
	newImage, err := image.Flip()
	check(err)
	write(name, newImage)
}
