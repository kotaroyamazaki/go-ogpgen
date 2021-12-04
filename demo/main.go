package main

import (
	"image"
	"image/color"

	"github.com/KotaroYamazaki/go-ogp-generator"
)

func main() {
	g, err := ogp.NewGenerator("base_image.png")
	if err != nil {
		panic(err)
	}
	g.SetQuality(100)

	if err := g.AttachText(&ogp.TextCompositionParams{
		Text: "Generate OGP ",
		TextPoint: image.Point{
			X: 600,
			Y: 150,
		},
		Color:    color.Black,
		FontSize: 56,
	}); err != nil {
		panic(err)
	}
	if err := g.AttachText(&ogp.TextCompositionParams{
		Text: "with any text or image!",
		TextPoint: image.Point{
			X: 600,
			Y: 150 + 64,
		},
		Color:    color.Black,
		FontSize: 56,
	}); err != nil {
		panic(err)
	}

	iconSize := 150
	if err := g.AttachImage(&ogp.ImageCompositionParams{
		ResizeWidth:  iconSize,
		ResizeHeight: iconSize,
		ImagePath:    "identicon.png",
		Mask: &ogp.Mask{
			Point: image.Point{
				X: 64 + iconSize/2,
				Y: 630 - iconSize/2 - 24,
			},
			Radius: 75,
		},
	}); err != nil {
		panic(err)
	}
	if err := g.AttachText(&ogp.TextCompositionParams{
		Text: "KotaroYamazaki",
		TextPoint: image.Point{
			X: 64 + iconSize*2,
			Y: 630 - 24 - 32,
		},
		FontSize: 32,
		Color:    color.White,
	}); err != nil {
		panic(err)
	}

	if err := g.Save("output.jpg"); err != nil {
		panic(err)
	}
}
