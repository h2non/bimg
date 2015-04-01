# bimg

Go library for blazing fast image processing based on [libvips](https://github.com/jcupitt/libvips) using C bindings

`Work in progress`

## Installation

```bash
go get gopkg.in/h2non/bimg.v0
```

Or get the lastest development version
```bash
go get github.com/h2non/bimg
```

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
