# bimg [![Build Status](https://travis-ci.org/h2non/bimg.png)](https://travis-ci.org/h2non/bimg) [![GitHub release](http://img.shields.io/github/tag/h2non/bimg.svg?style=flat-square)](https://github.com/h2non/bimg/releases) [![GoDoc](https://godoc.org/github.com/h2non/bimg?status.svg)](https://godoc.org/github.com/h2non/bimg) [![Coverage Status](https://coveralls.io/repos/h2non/bimg/badge.svg?branch=master)](https://coveralls.io/r/h2non/bimg?branch=master)

Small [Go](http://golang.org) package for fast high-level image processing and transformation using [libvips](https://github.com/jcupitt/libvips) via C bindings, provinding a simple, elegant and fluent [programmatic API](#examples).

bimg was designed to be a small and efficient library providing a generic high-level [image operations](#supported-image-operations) such as crop, resize, rotate, zoom, watermark...
It can read JPEG, PNG, WEBP and TIFF formats and output to JPEG, PNG and WEBP, including conversion between them.

bimg uses internally libvips, a powerful library written in C for image processing which requires a [low memory footprint](http://www.vips.ecs.soton.ac.uk/index.php?title=Speed_and_Memory_Use) 
and it's typically 4x faster than using the quickest ImageMagick and GraphicsMagick settings or Go native `image` package, and in some cases it's even 8x faster processing JPEG images. 

To get started you could take a look to the [examples](#examples) and [API](https://godoc.org/github.com/h2non/bimg) documentation. 
If you're looking for a HTTP based image processing solution, see [imaginary](https://github.com/h2non/imaginary). 
bimg was heavily inspired in [sharp](https://github.com/lovell/sharp), its homologous package built for node.js by [Lovell Fuller](https://github.com/lovell).

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

If you wanna take the advantage of [OpenSlide](http://openslide.org/), simply add `--with-openslide` to enable it:
```bash
curl -s https://raw.githubusercontent.com/lovell/sharp/master/preinstall.sh | sudo bash -s --with-openslide
```

The [install script](https://github.com/lovell/sharp/blob/master/preinstall.sh) requires `curl` and `pkg-config`

For platform specific installations, see  [Mac OS](https://github.com/lovell/sharp/blob/master/README.md#mac-os-tips) tips or [Windows](https://github.com/lovell/sharp/blob/master/README.md#windows) tips

## Supported image operations

- Resize
- Enlarge
- Crop
- Rotate (with auto-rotate based on EXIF orientation)
- Flip (with auto-flip based on EXIF metadata)
- Flop
- Zoom
- Thumbnail
- Extract area
- Watermark (text-based)
- Custom output color space (RGB, grayscale...)
- Format conversion (with additional quality/compression settings)
- EXIF metadata (size, alpha channel, profile, orientation...)

## Performance

libvips is probably the faster open source solution for image processing. 
Here you can see some performance test comparisons for multiple scenarios:

- [libvips speed and memory usage](http://www.vips.ecs.soton.ac.uk/index.php?title=Speed_and_Memory_Use)
- [sharp performance tests](https://github.com/lovell/sharp#the-task) 

#### Benchmarks

Tested using Go 1.4.1 and libvips-7.42.3 in OSX i7 2.7Ghz
```
BenchmarkResizeLargeJpeg  50    43400480 ns/op
BenchmarkResizePng        20    57592174 ns/op
BenchmarkResizeWebP       500    2872295 ns/op
BenchmarkConvertToJpeg    30    41835497 ns/op
BenchmarkConvertToPng     10   153382204 ns/op
BenchmarkConvertToWebp    10000   264542 ns/op
BenchmarkCropJpeg         30    52267699 ns/op
BenchmarkCropPng          30    56477454 ns/op
BenchmarkCropWebP         5000    274302 ns/op
BenchmarkExtractJpeg      50    27827670 ns/op
BenchmarkExtractPng       2000    769761 ns/op
BenchmarkExtractWebp      3000    513954 ns/op
BenchmarkZoomJpeg         10   159272494 ns/op
BenchmarkZoomPng          20    65771476 ns/op
BenchmarkZoomWebp         5000    368327 ns/op
BenchmarkWatermarkJpeg    100   10026033 ns/op
BenchmarkWatermarPng      200    7350821 ns/op
BenchmarkWatermarWebp     200    9014197 ns/op
ok 30.698s
```

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

#### Force resize

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

newImage, err := bimg.NewImage(buffer).ForceResize(1000, 500)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}
  
size := bimg.Size(newImage)
if size.Width != 1000 || size.Height != 500 {
  fmt.Fprintln(os.Stderr, "Incorrect image size")
}
```

#### Custom colour space (black & white)

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

newImage, err := bimg.NewImage(buffer).Colourspace(bimg.INTERPRETATION_B_W)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}
  
colourSpace, _ := bimg.ImageInterpretation(newImage)
if colourSpace != bimg.INTERPRETATION_B_W {
  fmt.Fprintln(os.Stderr, "Invalid colour space")
}
```

#### Custom options

See [Options](https://godoc.org/github.com/h2non/bimg#Options) struct to discover all the available fields

```go
options := bimg.Options{
  Width:        800,
  Height:       600,
  Crop:         true,
  Quality:      95,
  Rotate:       180,
  Interlace:    true,
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

#### Watermark

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

watermark := bimg.Watermark{
  Text:       "Chuck Norris (c) 2315",
  Opacity:    0.25,
  Width:      200,
  DPI:        100,
  Margin:     150,
  Font:       "sans bold 12",
  Background: bimg.Color{255, 255, 255},
}

newImage, err := bimg.NewImage(buffer).Watermark(watermark)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

bimg.Write("new.jpg", newImage)
```

#### Fluent interface

```go
buffer, err := bimg.Read("image.jpg")
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

image := bimg.NewImage(buffer)

// first crop image
_, err := image.CropByWidth(300)
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

// then flip it
newImage, err := image.Flip()
if err != nil {
  fmt.Fprintln(os.Stderr, err)
}

// save the cropped and flipped image
bimg.Write("new.jpg", newImage)
```

#### Debugging

Run the process passing the `DEBUG` environment variable
```
DEBUG=bimg ./app 
```

Enable libvips traces (note that a lot of data will be written in stdout):
```
VIPS_TRACE=1 ./app 
```

### Programmatic API 

#### func  ColourspaceIsSupported

```go
func ColourspaceIsSupported(buf []byte) (bool, error)
```
Check in the image colourspace is supported by libvips

#### func  DetermineImageTypeName

```go
func DetermineImageTypeName(buf []byte) string
```
Determines the image type format by name (jpeg, png, webp or tiff)

#### func  Initialize

```go
func Initialize()
```
Explicit thread-safe start of libvips. Only call this function if you've
previously shutdown libvips

#### func  IsTypeNameSupported

```go
func IsTypeNameSupported(t string) bool
```
Check if a given image type name is supported

#### func  IsTypeSupported

```go
func IsTypeSupported(t ImageType) bool
```
Check if a given image type is supported

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
Thread-safe function to shutdown libvips. You can call this to drop caches as
well. If libvips was already initialized, the function is no-op

#### func  VipsDebugInfo

```go
func VipsDebugInfo()
```
Output to stdout vips collected data. Useful for debugging

#### func  Write

```go
func Write(path string, buf []byte) error
```

#### type Angle

```go
type Angle int
```


```go
const (
  D0   Angle = 0
  D90  Angle = 90
  D180 Angle = 180
  D270 Angle = 270
)
```

#### type Color

```go
type Color struct {
  R, G, B uint8
}
```

Color represents a traditional RGB color scheme

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
Creates a new image

#### func (*Image) Colourspace

```go
func (i *Image) Colourspace(c Interpretation) ([]byte, error)
```
Colour space conversion

#### func (*Image) ColourspaceIsSupported

```go
func (i *Image) ColourspaceIsSupported() (bool, error)
```
Check if the current image has a valid colourspace

#### func (*Image) Convert

```go
func (i *Image) Convert(t ImageType) ([]byte, error)
```
Convert image to another format

#### func (*Image) Crop

```go
func (i *Image) Crop(width, height int, gravity Gravity) ([]byte, error)
```
Crop the image to the exact size specified

#### func (*Image) CropByHeight

```go
func (i *Image) CropByHeight(height int) ([]byte, error)
```
Crop an image by height (auto width)

#### func (*Image) CropByWidth

```go
func (i *Image) CropByWidth(width int) ([]byte, error)
```
Crop an image by width (auto height)

#### func (*Image) Enlarge

```go
func (i *Image) Enlarge(width, height int) ([]byte, error)
```
Enlarge the image by width and height. Aspect ratio is maintained

#### func (*Image) EnlargeAndCrop

```go
func (i *Image) EnlargeAndCrop(width, height int) ([]byte, error)
```
Enlarge the image by width and height with additional crop transformation

#### func (*Image) Extract

```go
func (i *Image) Extract(top, left, width, height int) ([]byte, error)
```
Extract area from the by X/Y axis

#### func (*Image) Flip

```go
func (i *Image) Flip() ([]byte, error)
```
Flip the image about the vertical Y axis

#### func (*Image) Flop

```go
func (i *Image) Flop() ([]byte, error)
```
Flop the image about the horizontal X axis

#### func (*Image) Image

```go
func (i *Image) Image() []byte
```
Get image buffer

#### func (*Image) Interpretation

```go
func (i *Image) Interpretation() (Interpretation, error)
```
Get the image interpretation type See:
http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/VipsImage.html#VipsInterpretation

#### func (*Image) Metadata

```go
func (i *Image) Metadata() (ImageMetadata, error)
```
Get image metadata (size, alpha channel, profile, EXIF rotation)

#### func (*Image) Process

```go
func (i *Image) Process(o Options) ([]byte, error)
```
Transform the image by custom options

#### func (*Image) Resize

```go
func (i *Image) Resize(width, height int) ([]byte, error)
```
Resize the image to fixed width and height

#### func (*Image) ResizeAndCrop

```go
func (i *Image) ResizeAndCrop(width, height int) ([]byte, error)
```
Resize the image to fixed width and height with additional crop transformation

#### func (*Image) Rotate

```go
func (i *Image) Rotate(a Angle) ([]byte, error)
```
Rotate the image by given angle degrees (0, 90, 180 or 270)

#### func (*Image) Size

```go
func (i *Image) Size() (ImageSize, error)
```
Get image size

#### func (*Image) Thumbnail

```go
func (i *Image) Thumbnail(pixels int) ([]byte, error)
```
Thumbnail the image by the a given width by aspect ratio 4:4

#### func (*Image) Type

```go
func (i *Image) Type() string
```
Get image type format (jpeg, png, webp, tiff)

#### func (*Image) Watermark

```go
func (i *Image) Watermark(w Watermark) ([]byte, error)
```
Add text as watermark on the given image

#### func (*Image) Zoom

```go
func (i *Image) Zoom(factor int) ([]byte, error)
```
Zoom the image by the given factor. You should probably call Extract() before

#### type ImageMetadata

```go
type ImageMetadata struct {
  Orientation int
  Channels    int
  Alpha       bool
  Profile     bool
  Type        string
  Space       string
  Colourspace string
  Size        ImageSize
}
```


#### func  Metadata

```go
func Metadata(buf []byte) (ImageMetadata, error)
```
Extract the image metadata (size, type, alpha channel, profile, EXIF
orientation...)

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
Get the image size by width and height pixels

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
Determines the image type format (jpeg, png, webp or tiff)

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

#### type Interpretation

```go
type Interpretation int
```

Image interpretation type See:
http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/VipsImage.html#VipsInterpretation

```go
const (
  INTERPRETATION_ERROR     Interpretation = C.VIPS_INTERPRETATION_ERROR
  INTERPRETATION_MULTIBAND Interpretation = C.VIPS_INTERPRETATION_MULTIBAND
  INTERPRETATION_B_W       Interpretation = C.VIPS_INTERPRETATION_B_W
  INTERPRETATION_CMYK      Interpretation = C.VIPS_INTERPRETATION_CMYK
  INTERPRETATION_RGB       Interpretation = C.VIPS_INTERPRETATION_RGB
  INTERPRETATION_sRGB      Interpretation = C.VIPS_INTERPRETATION_sRGB
  INTERPRETATION_RGB16     Interpretation = C.VIPS_INTERPRETATION_RGB16
  INTERPRETATION_GREY16    Interpretation = C.VIPS_INTERPRETATION_GREY16
  INTERPRETATION_scRGB     Interpretation = C.VIPS_INTERPRETATION_scRGB
)
```

#### func  ImageInterpretation

```go
func ImageInterpretation(buf []byte) (Interpretation, error)
```
Get the image interpretation type See:
http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/VipsImage.html#VipsInterpretation

#### type Options

```go
type Options struct {
  Height         int
  Width          int
  AreaHeight     int
  AreaWidth      int
  Top            int
  Left           int
  Extend         int
  Quality        int
  Compression    int
  Zoom           int
  Crop           bool
  Enlarge        bool
  Embed          bool
  Flip           bool
  Flop           bool
  Force          bool
  NoAutoRotate   bool
  NoProfile      bool
  Interlace      bool
  Rotate         Angle
  Gravity        Gravity
  Watermark      Watermark
  Type           ImageType
  Interpolator   Interpolator
  Interpretation Interpretation
}
```


#### type VipsMemoryInfo

```go
type VipsMemoryInfo struct {
  Memory          int64
  MemoryHighwater int64
  Allocations     int64
}
```


#### func  VipsMemory

```go
func VipsMemory() VipsMemoryInfo
```
Get memory info stats from vips (cache size, memory allocs...)

#### type Watermark

```go
type Watermark struct {
  Width       int
  DPI         int
  Margin      int
  Opacity     float32
  NoReplicate bool
  Text        string
  Font        string
  Background  Color
}
```

## Special Thanks

- [John Cupitt](https://github.com/jcupitt)

## License

MIT - Tomas Aparicio

[![views](https://sourcegraph.com/api/repos/github.com/h2non/bimg/.counters/views.svg)](https://sourcegraph.com/github.com/h2non/bimg)
