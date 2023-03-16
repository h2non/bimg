package vimgo

import "os"

// Read reads all the content of the given file path
// and returns it as byte buffer.
func Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Write writes the given byte buffer into disk
// to the given file path.
func Write(path string, buf []byte) error {
	return os.WriteFile(path, buf, 0644)
}
