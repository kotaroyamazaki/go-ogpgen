package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"

	"github.com/KotaroYamazaki/go-ogp-generator"
)

func main() {
	baseImg, err := readFile("base_image.png")
	if err != nil {
		panic(err)
	}
	g, err := ogp.NewGenerator(baseImg)
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
		FontSize: 56,
	}); err != nil {
		panic(err)
	}

	embedImg, err := readFile("identicon.png")
	if err != nil {
		panic(err)
	}
	iconSize := 150
	if err := g.AttachImage(&ogp.ImageCompositionParams{
		ResizeWidth:  iconSize,
		ResizeHeight: iconSize,
		Image:        embedImg,
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

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if f == nil {
		return nil, fmt.Errorf("error! Can not get image by %s", path)
	}
	return ioutil.ReadAll(f)
}
