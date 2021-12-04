package main

import (
	"image"
	"image/color"
	"ogpgen"
)

func main() {
	g, err := ogpgen.New("base_image.png")
	if err != nil {
		panic(err)
	}
	g.SetQuality(100)

	if err := g.AttachText(&ogpgen.TextCompositionParams{
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
	if err := g.AttachText(&ogpgen.TextCompositionParams{
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
	if err := g.AttachImage(&ogpgen.ImageCompositionParams{
		ResizeWidth:  iconSize,
		ResizeHeight: iconSize,
		ImagePath:    "identicon.png",
		Mask: &ogpgen.Mask{
			Point: image.Point{
				X: 64 + iconSize/2,
				Y: 630 - iconSize/2 - 24,
			},
			Radius: 75,
		},
	}); err != nil {
		panic(err)
	}
	if err := g.AttachText(&ogpgen.TextCompositionParams{
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
