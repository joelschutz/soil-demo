package util

import (
	"image"
	"image/color"
)

var (
	RED   = color.RGBA{255, 0, 0, 255}
	GREEN = color.RGBA{0, 255, 0, 255}
	BLUE  = color.RGBA{0, 0, 255, 255}
)

func AsRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()
	n := image.NewRGBA(rect)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			n.Set(x, y, img.At(x, y))
		}
	}
	return n
}

type LineSquare struct {
	Size   image.Rectangle
	Border uint
	Clr    color.Color
}

// ColorModel returns the Image's color model.
func (ln *LineSquare) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (ln *LineSquare) Bounds() image.Rectangle {
	return ln.Size
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (ln *LineSquare) At(x int, y int) color.Color {
	if x <= int(ln.Border) || x >= ln.Size.Dx()-int(ln.Border) {
		return ln.Clr
	}
	if y <= int(ln.Border) || y >= ln.Size.Dy()-int(ln.Border) {
		return ln.Clr
	}
	return color.Transparent
}
