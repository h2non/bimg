package bimg

import (
	"io/ioutil"
	"os"
)

func Read(path string) ([]byte, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
