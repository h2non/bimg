package bimg

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSize(t *testing.T) {
	files := []struct {
		name   string
		width  int
		height int
	}{
		{"test.jpg", 1680, 1050},
		{"test.png", 400, 300},
		{"test.webp", 550, 368},
	}
	for _, file := range files {
		size, err := Size(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %#v", err)
		}

		if size.Width != file.width || size.Height != file.height {
			t.Fatalf("Unexpected image size: %dx%d", size.Width, size.Height)
		}
	}
}

func TestMetadata(t *testing.T) {
	files := []struct {
		name        string
		format      string
		orientation int
		alpha       bool
		profile     bool
		space       string
	}{
		{"test.jpg", "jpeg", 0, false, false, "srgb"},
		{"test_icc_prophoto.jpg", "jpeg", 0, false, true, "srgb"},
		{"test.png", "png", 0, true, false, "srgb"},
		{"test.webp", "webp", 0, false, false, "srgb"},
		{"test.avif", "avif", 0, false, false, "srgb"},
	}

	for _, file := range files {
		metadata, err := Metadata(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}
		if metadata.Type != file.format {
			t.Fatalf("Unexpected image format: %s", file.format)
		}
		if metadata.Orientation != file.orientation {
			t.Fatalf("Unexpected image orientation: %d != %d", metadata.Orientation, file.orientation)
		}
		if metadata.Alpha != file.alpha {
			t.Fatalf("Unexpected image alpha: %t != %t", metadata.Alpha, file.alpha)
		}
		if metadata.Profile != file.profile {
			t.Fatalf("Unexpected image profile: %t != %t", metadata.Profile, file.profile)
		}
		if metadata.Space != file.space {
			t.Fatalf("Unexpected image profile: %t != %t", metadata.Profile, file.profile)
		}
	}
}

func TestImageInterpretation(t *testing.T) {
	files := []struct {
		name           string
		interpretation Interpretation
	}{
		{"test.jpg", InterpretationSRGB},
		{"test.png", InterpretationSRGB},
		{"test.webp", InterpretationSRGB},
	}

	for _, file := range files {
		interpretation, err := ImageInterpretation(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}
		if interpretation != file.interpretation {
			t.Fatalf("Unexpected image interpretation")
		}
	}
}

func TestEXIF(t *testing.T) {
	if VipsMajorVersion <= 8 && VipsMinorVersion < 10 {
		t.Skip("Skip test in libvips < 8.10")
		return
	}

	files := map[string]EXIF{
		"test.jpg": {},
		"exif/Landscape_1.jpg": {
			Orientation:      1,
			XResolution:      "72/1",
			YResolution:      "72/1",
			ResolutionUnit:   2,
			YCbCrPositioning: 1,
			ExifVersion:      "Exif Version 2.1",
			ColorSpace:       65535,
		},
		"test_exif.jpg": {
			Make:              "Jolla",
			Model:             "Jolla",
			XResolution:       "72/1",
			YResolution:       "72/1",
			ResolutionUnit:    2,
			Orientation:       1,
			Datetime:          "2014:09:21 16:00:56",
			ExposureTime:      "1/25",
			FNumber:           "12/5",
			ISOSpeedRatings:   320,
			ExifVersion:       "Exif Version 2.3",
			DateTimeOriginal:  "2014:09:21 16:00:56",
			ShutterSpeedValue: "205447286/44240665",
			ApertureValue:     "334328577/132351334",
			ExposureBiasValue: "0/1",
			MeteringMode:      1,
			Flash:             0,
			FocalLength:       "4/1",
			WhiteBalance:      1,
			ColorSpace:        65535,
		},
		"test_exif_canon.jpg": {
			Make:                    "Canon",
			Model:                   "Canon EOS 40D",
			Orientation:             1,
			XResolution:             "72/1",
			YResolution:             "72/1",
			ResolutionUnit:          2,
			Software:                "GIMP 2.4.5",
			Datetime:                "2008:07:31 10:38:11",
			YCbCrPositioning:        2,
			Compression:             6,
			ExposureTime:            "1/160",
			FNumber:                 "71/10",
			ExposureProgram:         1,
			ISOSpeedRatings:         100,
			ExifVersion:             "Exif Version 2.21",
			DateTimeOriginal:        "2008:05:30 15:56:01",
			DateTimeDigitized:       "2008:05:30 15:56:01",
			ComponentsConfiguration: "Y Cb Cr -",
			ShutterSpeedValue:       "483328/65536",
			ApertureValue:           "368640/65536",
			ExposureBiasValue:       "0/1",
			MeteringMode:            5,
			Flash:                   9,
			FocalLength:             "135/1",
			SubSecTimeOriginal:      "00",
			SubSecTimeDigitized:     "00",
			ColorSpace:              1,
			PixelXDimension:         100,
			PixelYDimension:         68,
			ExposureMode:            1,
			WhiteBalance:            0,
			SceneCaptureType:        0,
		},
		"test_exif_full.jpg": {
			Make:                    "Apple",
			Model:                   "iPhone XS",
			Orientation:             6,
			XResolution:             "72/1",
			YResolution:             "72/1",
			ResolutionUnit:          2,
			Software:                "13.3.1",
			Datetime:                "2020:07:28 19:18:49",
			YCbCrPositioning:        1,
			Compression:             6,
			ExposureTime:            "1/835",
			FNumber:                 "9/5",
			ExposureProgram:         2,
			ISOSpeedRatings:         25,
			ExifVersion:             "Unknown Exif Version",
			DateTimeOriginal:        "2020:07:28 19:18:49",
			DateTimeDigitized:       "2020:07:28 19:18:49",
			ComponentsConfiguration: "Y Cb Cr -",
			ShutterSpeedValue:       "77515/7986",
			ApertureValue:           "54823/32325",
			BrightnessValue:         "77160/8623",
			ExposureBiasValue:       "0/1",
			MeteringMode:            5,
			Flash:                   16,
			FocalLength:             "17/4",
			SubjectArea:             "2013 1511 2217 1330",
			MakerNote:               "1110 bytes undefined data",
			SubSecTimeOriginal:      "777",
			SubSecTimeDigitized:     "777",
			ColorSpace:              65535,
			PixelXDimension:         4032,
			PixelYDimension:         3024,
			SensingMethod:           2,
			SceneType:               "Directly photographed",
			ExposureMode:            0,
			WhiteBalance:            0,
			FocalLengthIn35mmFilm:   26,
			SceneCaptureType:        0,
			GPSLatitudeRef:          "N",
			GPSLatitude:             "55/1 43/1 5287/100",
			GPSLongitudeRef:         "E",
			GPSLongitude:            "37/1 35/1 5571/100",
			GPSAltitudeRef:          "Sea level",
			GPSAltitude:             "90514/693",
			GPSSpeedRef:             "K",
			GPSSpeed:                "114272/41081",
			GPSImgDirectionRef:      "M",
			GPSImgDirection:         "192127/921",
			GPSDestBearingRef:       "M",
			GPSDestBearing:          "192127/921",
			GPSDateStamp:            "2020:07:28",
		},
	}

	for name, file := range files {
		metadata, err := Metadata(readFile(name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", name, err)
		}
		if metadata.EXIF.Make != file.Make {
			t.Fatalf("Unexpected image exif Make: %s != %s", metadata.EXIF.Make, file.Make)
		}
		if metadata.EXIF.Model != file.Model {
			t.Fatalf("Unexpected image exif Model: %s != %s", metadata.EXIF.Model, file.Model)
		}
		if metadata.EXIF.Orientation != file.Orientation {
			t.Fatalf("Unexpected image exif Orientation: %d != %d", metadata.EXIF.Orientation, file.Orientation)
		}
		if metadata.EXIF.XResolution != file.XResolution {
			t.Fatalf("Unexpected image exif XResolution: %s != %s", metadata.EXIF.XResolution, file.XResolution)
		}
		if metadata.EXIF.YResolution != file.YResolution {
			t.Fatalf("Unexpected image exif YResolution: %s != %s", metadata.EXIF.YResolution, file.YResolution)
		}
		if metadata.EXIF.ResolutionUnit != file.ResolutionUnit {
			t.Fatalf("Unexpected image exif ResolutionUnit: %d != %d", metadata.EXIF.ResolutionUnit, file.ResolutionUnit)
		}
		if metadata.EXIF.Software != file.Software {
			t.Fatalf("Unexpected image exif Software: %s != %s", metadata.EXIF.Software, file.Software)
		}
		if metadata.EXIF.Datetime != file.Datetime {
			t.Fatalf("Unexpected image exif Datetime: %s != %s", metadata.EXIF.Datetime, file.Datetime)
		}
		if metadata.EXIF.YCbCrPositioning != file.YCbCrPositioning {
			t.Fatalf("Unexpected image exif YCbCrPositioning: %d != %d", metadata.EXIF.YCbCrPositioning, file.YCbCrPositioning)
		}
		if metadata.EXIF.Compression != file.Compression {
			t.Fatalf("Unexpected image exif Compression: %d != %d", metadata.EXIF.Compression, file.Compression)
		}
		if metadata.EXIF.ExposureTime != file.ExposureTime {
			t.Fatalf("Unexpected image exif ExposureTime: %s != %s", metadata.EXIF.ExposureTime, file.ExposureTime)
		}
		if metadata.EXIF.FNumber != file.FNumber {
			t.Fatalf("Unexpected image exif FNumber: %s != %s", metadata.EXIF.FNumber, file.FNumber)
		}
		if metadata.EXIF.ExposureProgram != file.ExposureProgram {
			t.Fatalf("Unexpected image exif ExposureProgram: %d != %d", metadata.EXIF.ExposureProgram, file.ExposureProgram)
		}
		if metadata.EXIF.ISOSpeedRatings != file.ISOSpeedRatings {
			t.Fatalf("Unexpected image exif ISOSpeedRatings: %d != %d", metadata.EXIF.ISOSpeedRatings, file.ISOSpeedRatings)
		}
		if metadata.EXIF.ExifVersion != file.ExifVersion {
			t.Fatalf("Unexpected image exif ExifVersion: %s != %s", metadata.EXIF.ExifVersion, file.ExifVersion)
		}
		if metadata.EXIF.DateTimeOriginal != file.DateTimeOriginal {
			t.Fatalf("Unexpected image exif DateTimeOriginal: %s != %s", metadata.EXIF.DateTimeOriginal, file.DateTimeOriginal)
		}
		if metadata.EXIF.DateTimeDigitized != file.DateTimeDigitized {
			t.Fatalf("Unexpected image exif DateTimeDigitized: %s != %s", metadata.EXIF.DateTimeDigitized, file.DateTimeDigitized)
		}
		if metadata.EXIF.ComponentsConfiguration != file.ComponentsConfiguration {
			t.Fatalf("Unexpected image exif ComponentsConfiguration: %s != %s", metadata.EXIF.ComponentsConfiguration, file.ComponentsConfiguration)
		}
		if metadata.EXIF.ShutterSpeedValue != file.ShutterSpeedValue {
			t.Fatalf("Unexpected image exif ShutterSpeedValue: %s != %s", metadata.EXIF.ShutterSpeedValue, file.ShutterSpeedValue)
		}
		if metadata.EXIF.ApertureValue != file.ApertureValue {
			t.Fatalf("Unexpected image exif ApertureValue: %s != %s", metadata.EXIF.ApertureValue, file.ApertureValue)
		}
		if metadata.EXIF.BrightnessValue != file.BrightnessValue {
			t.Fatalf("Unexpected image exif BrightnessValue: %s != %s", metadata.EXIF.BrightnessValue, file.BrightnessValue)
		}
		if metadata.EXIF.ExposureBiasValue != file.ExposureBiasValue {
			t.Fatalf("Unexpected image exif ExposureBiasValue: %s != %s", metadata.EXIF.ExposureBiasValue, file.ExposureBiasValue)
		}
		if metadata.EXIF.MeteringMode != file.MeteringMode {
			t.Fatalf("Unexpected image exif MeteringMode: %d != %d", metadata.EXIF.MeteringMode, file.MeteringMode)
		}
		if metadata.EXIF.Flash != file.Flash {
			t.Fatalf("Unexpected image exif Flash: %d != %d", metadata.EXIF.Flash, file.Flash)
		}
		if metadata.EXIF.FocalLength != file.FocalLength {
			t.Fatalf("Unexpected image exif FocalLength: %s != %s", metadata.EXIF.FocalLength, file.FocalLength)
		}
		if metadata.EXIF.SubjectArea != file.SubjectArea {
			t.Fatalf("Unexpected image exif SubjectArea: %s != %s", metadata.EXIF.SubjectArea, file.SubjectArea)
		}
		if metadata.EXIF.MakerNote != file.MakerNote {
			t.Fatalf("Unexpected image exif MakerNote: %s != %s", metadata.EXIF.MakerNote, file.MakerNote)
		}
		if metadata.EXIF.SubSecTimeOriginal != file.SubSecTimeOriginal {
			t.Fatalf("Unexpected image exif SubSecTimeOriginal: %s != %s", metadata.EXIF.SubSecTimeOriginal, file.SubSecTimeOriginal)
		}
		if metadata.EXIF.SubSecTimeDigitized != file.SubSecTimeDigitized {
			t.Fatalf("Unexpected image exif SubSecTimeDigitized: %s != %s", metadata.EXIF.SubSecTimeDigitized, file.SubSecTimeDigitized)
		}
		if metadata.EXIF.ColorSpace != file.ColorSpace {
			t.Fatalf("Unexpected image exif ColorSpace: %d != %d", metadata.EXIF.ColorSpace, file.ColorSpace)
		}
		if metadata.EXIF.PixelXDimension != file.PixelXDimension {
			t.Fatalf("Unexpected image exif PixelXDimension: %d != %d", metadata.EXIF.PixelXDimension, file.PixelXDimension)
		}
		if metadata.EXIF.PixelYDimension != file.PixelYDimension {
			t.Fatalf("Unexpected image exif PixelYDimension: %d != %d", metadata.EXIF.PixelYDimension, file.PixelYDimension)
		}
		if metadata.EXIF.SensingMethod != file.SensingMethod {
			t.Fatalf("Unexpected image exif SensingMethod: %d != %d", metadata.EXIF.SensingMethod, file.SensingMethod)
		}
		if metadata.EXIF.SceneType != file.SceneType {
			t.Fatalf("Unexpected image exif SceneType: %s != %s", metadata.EXIF.SceneType, file.SceneType)
		}
		if metadata.EXIF.ExposureMode != file.ExposureMode {
			t.Fatalf("Unexpected image exif ExposureMode: %d != %d", metadata.EXIF.ExposureMode, file.ExposureMode)
		}
		if metadata.EXIF.WhiteBalance != file.WhiteBalance {
			t.Fatalf("Unexpected image exif WhiteBalance: %d != %d", metadata.EXIF.WhiteBalance, file.WhiteBalance)
		}
		if metadata.EXIF.FocalLengthIn35mmFilm != file.FocalLengthIn35mmFilm {
			t.Fatalf("Unexpected image exif FocalLengthIn35mmFilm: %d != %d", metadata.EXIF.FocalLengthIn35mmFilm, file.FocalLengthIn35mmFilm)
		}
		if metadata.EXIF.SceneCaptureType != file.SceneCaptureType {
			t.Fatalf("Unexpected image exif SceneCaptureType: %d != %d", metadata.EXIF.SceneCaptureType, file.SceneCaptureType)
		}
		if metadata.EXIF.GPSLongitudeRef != file.GPSLongitudeRef {
			t.Fatalf("Unexpected image exif GPSLongitudeRef: %s != %s", metadata.EXIF.GPSLongitudeRef, file.GPSLongitudeRef)
		}
		if metadata.EXIF.GPSLongitude != file.GPSLongitude {
			t.Fatalf("Unexpected image exif GPSLongitude: %s != %s", metadata.EXIF.GPSLongitude, file.GPSLongitude)
		}
		if metadata.EXIF.GPSAltitudeRef != file.GPSAltitudeRef {
			t.Fatalf("Unexpected image exif GPSAltitudeRef: %s != %s", metadata.EXIF.GPSAltitudeRef, file.GPSAltitudeRef)
		}
		if metadata.EXIF.GPSAltitude != file.GPSAltitude {
			t.Fatalf("Unexpected image exif GPSAltitude: %s != %s", metadata.EXIF.GPSAltitude, file.GPSAltitude)
		}
		if metadata.EXIF.GPSSpeedRef != file.GPSSpeedRef {
			t.Fatalf("Unexpected image exif GPSSpeedRef: %s != %s", metadata.EXIF.GPSSpeedRef, file.GPSSpeedRef)
		}
		if metadata.EXIF.GPSSpeed != file.GPSSpeed {
			t.Fatalf("Unexpected image exif GPSSpeed: %s != %s", metadata.EXIF.GPSSpeed, file.GPSSpeed)
		}
		if metadata.EXIF.GPSImgDirectionRef != file.GPSImgDirectionRef {
			t.Fatalf("Unexpected image exif GPSImgDirectionRef: %s != %s", metadata.EXIF.GPSImgDirectionRef, file.GPSImgDirectionRef)
		}
		if metadata.EXIF.GPSImgDirection != file.GPSImgDirection {
			t.Fatalf("Unexpected image exif GPSImgDirection: %s != %s", metadata.EXIF.GPSImgDirection, file.GPSImgDirection)
		}
		if metadata.EXIF.GPSDestBearingRef != file.GPSDestBearingRef {
			t.Fatalf("Unexpected image exif GPSDestBearingRef: %s != %s", metadata.EXIF.GPSDestBearingRef, file.GPSDestBearingRef)
		}
		if metadata.EXIF.GPSDestBearing != file.GPSDestBearing {
			t.Fatalf("Unexpected image exif GPSDestBearing: %s != %s", metadata.EXIF.GPSDestBearing, file.GPSDestBearing)
		}
		if metadata.EXIF.GPSDateStamp != file.GPSDateStamp {
			t.Fatalf("Unexpected image exif GPSDateStamp: %s != %s", metadata.EXIF.GPSDateStamp, file.GPSDateStamp)
		}
	}
}

func TestColourspaceIsSupported(t *testing.T) {
	files := []struct {
		name string
	}{
		{"test.jpg"},
		{"test.png"},
		{"test.webp"},
	}

	for _, file := range files {
		supported, err := ColourspaceIsSupported(readFile(file.name))
		if err != nil {
			t.Fatalf("Cannot read the image: %s -> %s", file.name, err)
		}
		if supported != true {
			t.Fatalf("Unsupported image colourspace")
		}
	}

	supported, err := initImage("test.jpg").ColourspaceIsSupported()
	if err != nil {
		t.Errorf("Cannot process the image: %#v", err)
	}
	if supported != true {
		t.Errorf("Non-supported colourspace")
	}
}

func readFile(file string) []byte {
	data, _ := os.Open(path.Join("testdata", file))
	buf, _ := ioutil.ReadAll(data)
	return buf
}
