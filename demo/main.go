package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/KotaroYamazaki/go-ogp-generator"
)

func main() {
	img, err := readFile("base_image.png")
	if err != nil {
		panic(err)
	}
	g, err := ogp.NewGenerator(img)
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

	img2, err := readFile("identicon.png")
	if err != nil {
		panic(err)
	}
	iconSize := 150
	if err := g.AttachImage(&ogp.ImageCompositionParams{
		ResizeWidth:  iconSize,
		ResizeHeight: iconSize,
		Image:        img2,
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

	og, err := g.Get()
	if err != nil {
		panic(err)
	}
	_img, _, err := image.Decode(bytes.NewReader(og))
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create("./img.jpeg")
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 70

	err = jpeg.Encode(out, _img, &opts)
	if err != nil {
		log.Println(err)
	}
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if f == nil {
		return nil, fmt.Errorf("Error! Can not get Image by %s.", path)
	}
	return ioutil.ReadAll(f)
}
