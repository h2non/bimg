# bimg [![Build Status](https://travis-ci.org/h2non/bimg.png)](https://travis-ci.org/h2non/bimg) [![GitHub release](https://img.shields.io/github/tag/h2non/bimg.svg)](https://github.com/h2non/bimg/releases) [![GoDoc](https://godoc.org/github.com/h2non/bimg?status.png)](https://godoc.org/github.com/h2non/bimg) [![Coverage Status](https://coveralls.io/repos/h2non/bimg/badge.svg?branch=master)](https://coveralls.io/r/h2non/bimg?branch=master)

Small Go library for blazing fast and efficient image processing based on [libvips](https://github.com/jcupitt/libvips) using C bindings.

bimg is designed to be a small and efficient library with a generic and useful set of features. 
It uses internally libvips, which is typically 4x faster than using the quickest ImageMagick and GraphicsMagick settings or Go  native `image` package, and in some cases it's even 8x faster processing JPEG images. 

It can read JPEG, PNG, WEBP, TIFF and Magick formats and it can output to JPEG, PNG and WEBP. It supports common [image transformation](#supported-image-operations) operations such as crop, resize, rotate... and conversion into multiple formats. 

To getting start take a look to the [examples](#examples) and [programmatic API](https://godoc.org/github.com/h2non/bimg) documentation.

bimg was heavily inspired in [sharp](https://github.com/lovell/sharp), a great node.js package for image processing build by [Lovell Fuller](https://github.com/lovell).

**Note**: bimg is still beta. Pull requests and issues are highly appreciated

## Prerequisites

- [libvips](https://github.com/jcupitt/libvips) v7.40.0+ (7.42.0+ recommended)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- Go 1.3+

## Installation

```bash
go get -u gopkg.in/h2non/bimg.v0
```

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
- Rotate
- Flip 
- Thumbnail
- Extract area
- Format conversion
- EXIF metadata (size, alpha channel, profile, orientation...)

## Performance

libvips is probably the faster open source solution for image processing. 
Here you can see some performance test comparisons for multiple scenarios:

- [libvips speed and memory usage](http://www.vips.ecs.soton.ac.uk/index.php?title=Speed_and_Memory_Use)
- [sharp performance tests](https://github.com/lovell/sharp#the-task) 

bimg performance tests coming soon!

## API

### Examples

```go
import (
  "fmt"
  "os"
  "gopkg.in/h2non/bimg.v0"
)
```

#### Resize

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

newImage, err := bimg.NewImage(buffer).Resize(800, 600)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

size, err := bimg.NewImage(newImage).Size()
if size.Width == 400 && size.Height == 300 {
  fmt.Println("The image size is valid")
}

bimg.Write("new.jpg", newImage)
```

#### Rotate

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

newImage, err := bimg.NewImage(buffer).Rotate(90)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

bimg.Write("new.jpg", newImage)
```

#### Convert

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

newImage, err := bimg.NewImage(buffer).Convert(bimg.PNG)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

if bimg.NewImage(newImage).Type() == "png" {
  fmt.Fprintln(os.Stderr, "The image was converted into png")
}
```

#### Custom options

See [Options](https://godoc.org/github.com/h2non/bimg#Options) struct to see all available fields

```go
options := bimg.Options{
  Width:        800,
  Height:       600,
  Crop:         true,
  Quality:      95,
  Rotate:       180,
}

buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

newImage, err := bimg.NewImage(buffer).Process(options)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

bimg.Write("new.jpg", newImage)
```

#### func  DetermineImageTypeName

```go
func DetermineImageTypeName(buf []byte) string
```

#### func  Initialize

```go
func Initialize()
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

#### func  Shutdown

```go
func Shutdown()
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

#### func  DetermineImageType

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

## License

MIT - Tomas Aparicio
