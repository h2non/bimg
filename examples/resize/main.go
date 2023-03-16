package main

import (
	"fmt"
	"os"

	"github.com/nestor-sk/vimgo"
)

func main() {
	buffer, err := vimgo.Read("./testdata/test.jpg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newImage, err := vimgo.NewImage(buffer).Resize(800, 600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	size, err := vimgo.NewImage(newImage).Size()
	if size.Width == 800 && size.Height == 600 {
		fmt.Println("The image size is valid")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	vimgo.Write("./bin/new.jpg", newImage)
}
