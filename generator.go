package ogp

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

const (
	OGPMaxWidth  = 1200
	OGPMaxHeight = 630
)

//go:embed fonts/MPLUSRounded1c-Bold.ttf
var defaultFont []byte

type Generator interface {
	AttachImage(*ImageCompositionParams) error
	AttachText(*TextCompositionParams) error
	Get() ([]byte, error)
	SetQuality(int)
	SetSize(int, int)
}

type generator struct {
	img     *image.RGBA
	quality int
	Height  int
	Width   int
}

func NewGenerator(b []byte) (Generator, error) {
	img, err := anyDecode(b)
	if err != nil {
		return nil, err
	}
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{}, draw.Src)
	return &generator{
		img:     out,
		quality: 70,
	}, nil
}

func (c *generator) Get() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, c.img, &jpeg.Options{Quality: c.quality})
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// SetQuality sets the quality of the generated image. 0 is lowest and 100 is highest.
func (c *generator) SetQuality(q int) {
	if q <= 0 || q > 100 {
		q = 70
	}
	c.quality = q
}

// SetSize sets the size of the generated image.
func (c *generator) SetSize(w, h int) {
	if w <= 0 || w > OGPMaxWidth {
		c.Width = OGPMaxWidth
	} else {
		c.Width = w
	}
	if h <= 0 || h > OGPMaxHeight {
		c.Height = OGPMaxHeight
	} else {
		c.Height = h
	}
}

type ImageCompositionParams struct {
	ResizeWidth  int
	ResizeHeight int
	Image        []byte
	Mask         *Mask
}

// AttachImage attaches an image to the base image.
func (c *generator) AttachImage(params *ImageCompositionParams) error {
	img, err := anyDecode(params.Image)
	if err != nil {
		return err
	}
	if params.ResizeWidth != 0 && params.ResizeHeight != 0 {
		img = resize(img, params.ResizeHeight, params.ResizeWidth)
	}
	if params.Mask != nil {
		sp := image.Pt(params.ResizeWidth/2-params.Mask.Point.X, params.ResizeHeight/2-params.Mask.Point.Y)
		draw.DrawMask(c.img, c.img.Bounds(), img, sp, params.Mask, image.Point{}, draw.Over)
	} else {
		sp := image.Point{}
		draw.Draw(c.img, img.Bounds(), img, sp, draw.Over)
	}
	return nil
}

type TextCompositionParams struct {
	Text      string
	TextPoint image.Point
	Color     color.Color
	FontSize  int
	FontPath  string
}

func (p *TextCompositionParams) validate() error {
	if p.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}

// AttachText attaches text to the base image.
func (c *generator) AttachText(params *TextCompositionParams) error {
	if err := params.validate(); err != nil {
		return err
	}
	if params.TextPoint.X == 0 && params.TextPoint.Y == 0 {
		params.TextPoint = image.Point{c.Width / 2,
			c.Height / 2}
	}
	if params.FontSize == 0 {
		params.FontSize = 64
	}
	if params.Color == nil {
		params.Color = color.Black
	}
	textColor := image.NewUniform(params.Color)

	var font []byte
	if params.FontPath == "" {
		font = defaultFont
	} else {
		var err error
		font, err = ioutil.ReadFile(params.FontPath)
		if err != nil {
			return err
		}
	}

	_font, err := truetype.Parse(font)
	if err != nil {
		return err
	}
	f := freetype.NewContext()
	f.SetFont(_font)
	f.SetFontSize(float64(params.FontSize))
	f.SetDst(c.img)
	f.SetClip(c.img.Bounds())
	f.SetSrc(textColor)

	textWidth := c.getTextWidth(float64(params.FontSize), params.Text, _font)
	pt := freetype.Pt(params.TextPoint.X-textWidth/2, params.TextPoint.Y+params.FontSize/2)
	_, err = f.DrawString(params.Text, pt)
	return err
}

// getTextWidth returns the width of the text in the font.
func (c *generator) getTextWidth(fontSize float64, text string, fonts *truetype.Font) int {
	var textWidth int
	var face font.Face
	opts := truetype.Options{}
	opts.Size = fontSize
	face = truetype.NewFace(fonts, &opts)
	for _, x := range text {
		awidth, ok := face.GlyphAdvance(rune(x))
		if !ok {
			return textWidth
		}
		textWidth += int(float64(awidth) / 64)
	}
	return textWidth
}
