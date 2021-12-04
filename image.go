package ogpgen

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"golang.org/x/image/draw"
)

func resize(src image.Image, w, h int) *image.RGBA {
	rct := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, rct, draw.Over, nil)
	return dst
}

func generateRandomImageUniform() *image.Uniform {
	rand.Seed(time.Now().UnixNano())
	return image.NewUniform(color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	})
}
