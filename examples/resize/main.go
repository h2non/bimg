package main

import (
	"fmt"
	"os"

	bimg "github.com/nestor-sk/vimgo"
)

func main() {
	buffer, err := bimg.Read("./testdata/test.jpg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newImage, err := bimg.NewImage(buffer).Resize(800, 600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	size, err := bimg.NewImage(newImage).Size()
	if size.Width == 800 && size.Height == 600 {
		fmt.Println("The image size is valid")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	bimg.Write("./bin/new.jpg", newImage)
}
