package ogp

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/h2non/filetype"
	"golang.org/x/image/draw"
)

const (
	MediaTypeImage = "image"
)

var SupportedMediaFormats = map[string]string{
	"png":  MediaTypeImage,
	"jpeg": MediaTypeImage,
	"gif":  MediaTypeImage,
}

func validateMediaType(data []byte) (string, string, bool) {
	kind, err := filetype.Match(data)
	if err != nil {
		return "", "", false
	}
	typ, ok := SupportedMediaFormats[kind.MIME.Subtype]
	if !ok {
		return "", "", false
	}
	return typ, kind.Extension, true
}

func anyDecode(b []byte) (image.Image, error) {
	_, ext, ok := validateMediaType(b)
	if !ok {
		return nil, fmt.Errorf("error! file extension is not allowed")
	}

	br := bytes.NewReader(b)
	var img image.Image
	var err error
	switch ext {
	case "jpg":
		img, err = jpeg.Decode(br)
	case "png":
		img, err = png.Decode(br)
	case "gif":
		img, err = gif.Decode(br)
	default:
		return nil, fmt.Errorf("error! file extension %s is not allowed", ext)
	}
	if err != nil {
		return nil, err
	}
	return img, nil
}

func resize(src image.Image, w, h int) *image.RGBA {
	rct := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, rct, draw.Over, nil)
	return dst
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if f == nil {
		return nil, fmt.Errorf("error! Can not get Image by %s", path)
	}
	return ioutil.ReadAll(f)
}
