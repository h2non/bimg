package bimg

import "io/ioutil"

func Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func Write(path string, buf []byte) error {
	return ioutil.WriteFile(path, buf, 0644)
}
