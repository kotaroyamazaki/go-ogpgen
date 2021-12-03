package ogp

import (
	"image"
	"image/color"
)

type Mask struct {
	Point  image.Point
	Radius int
}

func NewMask(x, y, r int) *Mask {
	return &Mask{
		Point:  image.Pt(x, y),
		Radius: r,
	}
}

func (m *Mask) ColorModel() color.Model {
	return color.AlphaModel
}

func (m *Mask) Bounds() image.Rectangle {
	return image.Rect(m.Point.X-m.Radius, m.Point.Y-m.Radius, m.Point.X+m.Radius, m.Point.Y+m.Radius)
}

func (m *Mask) At(x, y int) color.Color {
	xx, yy, rr := float64(x-m.Point.X)+0.5, float64(y-m.Point.Y)+0.5, float64(m.Radius)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
