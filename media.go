package ogpgen

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/h2non/filetype"
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
