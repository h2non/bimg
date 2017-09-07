// +build !go17

package bimg

// Resize is used to transform a given image as byte buffer
// with the passed options.
// Used as proxy to resizer()
func Resize(buf []byte, o Options) ([]byte, error) {
  return resizer(buf, o)
}
