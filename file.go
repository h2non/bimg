package bimg

import (
	"io/ioutil"
	"os"
)

func Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func Write(path string, buf []byte) error {
	return ioutil.WriteFile(path, buf, 0644)
}
