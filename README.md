# bimg [![Build Status](https://travis-ci.org/h2non/bimg.png)](https://travis-ci.org/h2non/bimg) [![GitHub release](https://img.shields.io/github/tag/h2non/bimg.svg)]() [![GoDoc](https://godoc.org/github.com/h2non/bimg?status.png)](https://godoc.org/github.com/h2non/bimg)

Go library for blazing fast image processing based on [libvips](https://github.com/jcupitt/libvips) using C bindings. 

libvips is typically 4x faster than using the quickest ImageMagick and GraphicsMagick settings, and even 8x faster with JPEG format. It support JPEG, PNG, WEBP, TIFF and Magick image formats.

bimg was heavily inspired in [sharp](https://github.com/lovell/sharp), a great node.js package for image processing build by [Lovell Fuller](https://github.com/lovell).

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

The [install script](https://github.com/lovell/sharp/blob/master/preinstall.sh) requires `curl` and `pkg-config`

## Supported image operations

- Resize
- Enlarge
- Crop
- Zoom
- Rotate
- Flip/Flop
- Extract area
- Extract image metadata (size, format, profile, orientation...)
- Image conversion to multiple formats

## API

### Example

```go
import (
  "fmt"
  "os"
  "gopkg.in/h2non/bimg.v0"
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

#### func  DetermineImageTypeName

```go
func DetermineImageTypeName(buf []byte) string
```

#### func  IsTypeNameSupported

```go
func IsTypeNameSupported(t string) bool
```

#### func  IsTypeSupported

```go
func IsTypeSupported(t ImageType) bool
```

#### func  Read

```go
func Read(path string) ([]byte, error)
```

#### func  Resize

```go
func Resize(buf []byte, o Options) ([]byte, error)
```

#### type Angle

```go
type Angle int
```


```go
const (
  D0   Angle = C.VIPS_ANGLE_D0
  D90  Angle = C.VIPS_ANGLE_D90
  D180 Angle = C.VIPS_ANGLE_D180
  D270 Angle = C.VIPS_ANGLE_D270
)
```

#### type Direction

```go
type Direction int
```


```go
const (
  HORIZONTAL Direction = C.VIPS_DIRECTION_HORIZONTAL
  VERTICAL   Direction = C.VIPS_DIRECTION_VERTICAL
)
```

#### type Gravity

```go
type Gravity int
```


```go
const (
  CENTRE Gravity = iota
  NORTH
  EAST
  SOUTH
  WEST
)
```

#### type Image

```go
type Image struct {
}
```


#### func  NewImage

```go
func NewImage(buf []byte) *Image
```

#### func (*Image) Convert

```go
func (i *Image) Convert(t ImageType) ([]byte, error)
```

#### func (*Image) Crop

```go
func (i *Image) Crop(width int, height int) ([]byte, error)
```

#### func (*Image) Extract

```go
func (i *Image) Extract(top int, left int, width int, height int) ([]byte, error)
```

#### func (*Image) Flip

```go
func (i *Image) Flip() ([]byte, error)
```

#### func (*Image) Flop

```go
func (i *Image) Flop() ([]byte, error)
```

#### func (*Image) Metadata

```go
func (i *Image) Metadata() (ImageMetadata, error)
```

#### func (*Image) Process

```go
func (i *Image) Process(o Options) ([]byte, error)
```

#### func (*Image) Resize

```go
func (i *Image) Resize(width int, height int) ([]byte, error)
```

#### func (*Image) Rotate

```go
func (i *Image) Rotate(a Angle) ([]byte, error)
```

#### func (*Image) Size

```go
func (i *Image) Size() (ImageSize, error)
```

#### func (*Image) Type

```go
func (i *Image) Type() string
```

#### type ImageMetadata

```go
type ImageMetadata struct {
  Orientation int
  Alpha       bool
  Profile     bool
  Space       int
  Type        string
  Size        ImageSize
}
```


#### func  Metadata

```go
func Metadata(buf []byte) (ImageMetadata, error)
```

#### type ImageSize

```go
type ImageSize struct {
  Width  int
  Height int
}
```


#### func  Size

```go
func Size(buf []byte) (ImageSize, error)
```

#### type ImageType

```go
type ImageType int
```


```go
const (
  UNKNOWN ImageType = iota
  JPEG
  WEBP
  PNG
  TIFF
  MAGICK
)
```

#### func DetermineImageType

```go
func DetermineImageType(buf []byte) ImageType
```

#### type Interpolator

```go
type Interpolator int
```


```go
const (
  BICUBIC Interpolator = iota
  BILINEAR
  NOHALO
)
```

#### func (Interpolator) String

```go
func (i Interpolator) String() string
```

#### type Options

```go
type Options struct {
  Height       int
  Width        int
  Top          int
  Left         int
  Crop         bool
  Enlarge      bool
  Extend       int
  Embed        bool
  Quality      int
  Compression  int
  Type         ImageType
  Rotate       Angle
  Flip         Direction
  Gravity      Gravity
  Interpolator Interpolator
}
```


#### type Vips

```go
type Vips struct {
}
```

## License

MIT - Tomas Aparicio
