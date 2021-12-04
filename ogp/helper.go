package ogp

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/h2non/filetype"
	"golang.org/x/image/draw"
)

var SupportedMediaFormats = map[string]bool{
	"png":  true,
	"jpeg": true,
	"gif":  true,
}

func validateMediaType(data []byte) (string, bool) {
	kind, err := filetype.Match(data)
	if err != nil {
		return "", false
	}
	if !SupportedMediaFormats[kind.Extension] {
		return "", false
	}
	return kind.Extension, true
}

func anyDecode(b []byte) (image.Image, error) {
	ext, ok := validateMediaType(b)
	if !ok {
		return nil, fmt.Errorf("file extension '%s' is not allowed", ext)
	}

	br := bytes.NewReader(b)
	switch ext {
	case "jpeg":
		return jpeg.Decode(br)
	case "jpg":
		return jpeg.Decode(br)
	case "png":
		return png.Decode(br)
	case "gif":
		return gif.Decode(br)
	default:
		return nil, fmt.Errorf("file extension '%s' is not allowed", ext)
	}
}

func resize(src image.Image, w, h int) *image.RGBA {
	rct := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, rct, draw.Over, nil)
	return dst
}
