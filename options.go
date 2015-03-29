package bimg

const QUALITY = 80

type Gravity int

type Interpolator int

var interpolations = map[Interpolator]string{
	BICUBIC:  "bicubic",
	BILINEAR: "bilinear",
	NOHALO:   "nohalo",
}

func (i Interpolator) String() string {
	return interpolations[i]
}

type Options struct {
	Height       int
	Width        int
	Crop         bool
	Enlarge      bool
	Extend       int
	Embed        bool
	Interpolator Interpolator
	Gravity      Gravity
	Quality      int
}
