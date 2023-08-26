package boards

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnumBoard struct {
	values     [][]color.Color
	hvrX, hvrY int
}

func (ba *EnumBoard) Size() (int, int) {
	return len(ba.values), len(ba.values[0])
}

func (ba *EnumBoard) Update() error {
	return nil
}

func (ba *EnumBoard) Draw(screen *ebiten.Image) {
	for x, row := range ba.values {
		for y, v0 := range row {
			screen.Set(x, y, v0)
		}
	}
}

func (ba *EnumBoard) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 0, 0
}

func (ba *EnumBoard) Setup(init [][]color.Color) {
	ba.values = init
}

func (ba *EnumBoard) Reset() error {
	return nil
}

func (ba *EnumBoard) GetState() [][]color.Color {
	return ba.values
}

func (ba *EnumBoard) Click(btn ebiten.MouseButton) {
	return
}

func (ba *EnumBoard) Hover(x, y int) {
	ba.hvrX = x
	ba.hvrY = y
}
