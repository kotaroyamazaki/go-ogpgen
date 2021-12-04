package ogpgen

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

const (
	ogpMaxWidth    = 1200
	ogpMaxHeight   = 630
	defaultQuality = 70
)

type Generator interface {
	AttachImage(*ImageCompositionParams) error
	AttachText(*TextCompositionParams) error
	Get() ([]byte, error)
	Save(string) error
	SetQuality(int)
	SetSize(int, int)
}

type ogpgen struct {
	img     *image.RGBA
	quality int
	height  int
	width   int
}

func New(path string) (Generator, error) {
	b, err := readImageBytes(path)
	if err != nil {
		return nil, err
	}
	img, err := anyDecode(b)
	if err != nil {
		return nil, err
	}
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{}, draw.Src)
	return &ogpgen{
		img:     out,
		quality: defaultQuality,
		width:   out.Rect.Dx(),
		height:  out.Rect.Dy(),
	}, nil
}

// NewRandomColor returns generator which have a random color base image.
func NewRandomBackground() Generator {
	out := image.NewRGBA(image.Rect(0, 0, ogpMaxWidth, ogpMaxHeight))
	draw.Draw(out, out.Bounds(), generateRandomImageUniform(), image.Point{}, draw.Src)
	return &ogpgen{
		img:     out,
		quality: defaultQuality,
		width:   out.Rect.Dx(),
		height:  out.Rect.Dy(),
	}
}

// Get returns the image as []byte.
func (c *ogpgen) Get() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, c.img, &jpeg.Options{Quality: c.quality})
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Save saves the image to the specified path.
func (c *ogpgen) Save(path string) error {
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, c.img, &jpeg.Options{Quality: c.quality})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, buff.Bytes(), 0644)
}

// SetQuality sets the quality of the generated image. 0 is lowest and 100 is highest.
func (c *ogpgen) SetQuality(q int) {
	if q <= 0 {
		q = defaultQuality
	}
	if q > 100 {
		q = 100
	}
	c.quality = q
}

// SetSize sets the size of the generated image.
func (c *ogpgen) SetSize(w, h int) {
	if w <= 0 || w > ogpMaxWidth {
		c.width = ogpMaxWidth
	} else {
		c.width = w
	}
	if h <= 0 || h > ogpMaxHeight {
		c.height = ogpMaxHeight
	} else {
		c.height = h
	}
}

type ImageCompositionParams struct {
	ResizeWidth  int
	ResizeHeight int
	ImagePath    string
	Mask         *Mask
}

// AttachImage attaches an image to the base image.
func (c *ogpgen) AttachImage(params *ImageCompositionParams) error {
	b, err := readImageBytes(params.ImagePath)
	if err != nil {
		return err
	}
	img, err := anyDecode(b)
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

//go:embed fonts/MPLUSRounded1c-Bold.ttf
var defaultFont []byte

func (p *TextCompositionParams) validate() error {
	if p.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}

// AttachText attaches text to the base image.
func (c *ogpgen) AttachText(params *TextCompositionParams) error {
	if err := params.validate(); err != nil {
		return err
	}
	if params.TextPoint.X == 0 && params.TextPoint.Y == 0 {
		params.TextPoint = image.Point{
			c.width / 2,
			c.height / 2,
		}
	}
	if params.FontSize == 0 {
		params.FontSize = 64
	}
	if params.Color == nil {
		params.Color = color.White
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

func (c *ogpgen) getTextWidth(fontSize float64, text string, fonts *truetype.Font) int {
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

func readImageBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
