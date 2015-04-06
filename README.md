# bimg [![Build Status](https://travis-ci.org/h2non/bimg.png)](https://travis-ci.org/h2non/bimg) [![GitHub release](https://img.shields.io/github/tag/h2non/bimg.svg)]() [![GoDoc](https://godoc.org/github.com/h2non/bimg?status.png)](https://godoc.org/github.com/h2non/bimg)

Go library for blazing fast image processing based on [libvips](https://github.com/jcupitt/libvips) using C bindings. 

**bimg** was focused on performance, resizing an image with libvips is typically 4x faster than using the quickest ImageMagick and GraphicsMagick settings.

**bimg** was heavily inspired in [sharp](https://github.com/lovell/sharp), a great node.js package for image processing build by [Lovell Fuller](https://github.com/lovell).

`Work in progress`

## Prerequisites

- [libvips](https://github.com/jcupitt/libvips) v7.40.0+ (7.42.0+ recommended)
- C++11 compatible compiler such as gcc 4.6+ or clang 3.0+

## Installation

```bash
go get gopkg.in/h2non/bimg.v0
```
Requires Go 1.3+

### libvips

Run the following script as `sudo` (supports OSX, Debian/Ubuntu, Redhat, Fedora, Amazon Linux):
```bash
curl -s https://raw.githubusercontent.com/lovell/sharp/master/preinstall.sh | sudo bash -
```

The [install script](https://github.com/lovell/sharp/blob/master/preinstall.sh) requires `curl` and `pkg-config`.

## API

```go
import (
  "fmt"
  "os"
  "gopkg.in/h2non/bimg"
)

options := bimg.Options{
    Width:        800,
    Height:       600,
    Crop:         true,
    Quality:      95,
}

newImage, err := bimg.Resize(image, options)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}
```

## License

MIT - Tomas Aparicio
