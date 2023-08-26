package internal

import (
	"image"
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/soil-demo/internal/boards"
	"github.com/joelschutz/soil-demo/util"
	"github.com/joelschutz/stagehand"
	"github.com/solarlune/ldtkgo"
)

type Board[V any] interface {
	ebiten.Game
	Size() (int, int)
	Reset() error
	Setup(m [][]V)
	Click(btn ebiten.MouseButton)
	Hover(x, y int)
}

type State struct {
	age         uint
	paused      bool
	speed       uint
	sceneNum    uint
	hover       int
	scaleFac    float64
	isPreview   bool
	menuScale   float64
	Levels      []*ldtkgo.Level
	boardStates [][]mgl32.Vec2
	Config      Config
}

type SimulationScene struct {
	Board   *boards.HumidityBoard
	Preview Board[color.Color]
	sm      *stagehand.SceneManager[State]
	state   State
}

func (s *SimulationScene) Update() error {
	if !s.state.paused {
		for i := 0; i < int(s.state.speed+1); i++ {
			s.Board.Update()
		}
	}

	cx, cy := ebiten.CursorPosition()
	if float64(cx) < float64(s.state.Config.btnSize)*s.state.menuScale {
		s.state.hover = cy/int(float64(s.state.Config.btnSize)*s.state.menuScale) + 1
		s.Board.Hover(-1, -1)
	} else {
		s.state.hover = 0
		bx := int((float64(cx) - float64(s.state.Config.btnSize)*s.state.menuScale) / s.state.scaleFac)
		by := int(float64(cy) / s.state.scaleFac)
		s.Board.Hover(bx, by)
		for btn := ebiten.MouseButtonLeft; btn < ebiten.MouseButtonMiddle; btn++ {
			if ebiten.IsMouseButtonPressed(btn) {
				s.Board.Click(btn)
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		switch s.state.hover {
		case 1:
			s.state.paused = !s.state.paused
		case 2:
			s.Board.Reset()
		case 3:
			s.state.speed++
			if s.state.speed > 4 {
				s.state.speed = 0
			}
		case 4:
			s.state.isPreview = !s.state.isPreview
		case 5:
			s.state.sceneNum++
			if s.state.sceneNum > 3 {
				s.state.sceneNum = 0
			}
			s.sm.SwitchTo(&SimulationScene{
				Board:   &boards.HumidityBoard{},
				Preview: &boards.EnumBoard{},
			})
		}
	}
	s.state.age++
	return nil
}

func (s *SimulationScene) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(s.Board.Size())
	if !s.state.isPreview {
		s.Board.Draw(img)
	} else {
		s.Preview.Draw(img)
	}

	op := &ebiten.DrawImageOptions{}
	s.state.scaleFac = float64(screen.Bounds().Dy()) / float64(img.Bounds().Dy())
	op.GeoM.Scale(s.state.scaleFac, s.state.scaleFac)
	op.GeoM.Translate(float64(s.state.Config.btnSize)*s.state.menuScale, 0)
	screen.DrawImage(img, op)

	invertedClr := ebiten.ColorM{}
	invertedClr.Scale(-1, -1, -1, 1)
	invertedClr.Translate(1, 1, 1, 0)

	opPlay := &ebiten.DrawImageOptions{}
	opPlay.GeoM.Scale(s.state.menuScale, s.state.menuScale)

	opReset := &ebiten.DrawImageOptions{}
	opReset.GeoM.Scale(s.state.menuScale, s.state.menuScale)
	opReset.GeoM.Translate(0, float64(s.state.Config.btnSize)*s.state.menuScale)

	opSpeed := &ebiten.DrawImageOptions{}
	opSpeed.GeoM.Scale(s.state.menuScale, s.state.menuScale)
	opSpeed.GeoM.Translate(0, 2*float64(s.state.Config.btnSize)*s.state.menuScale)

	opTree := &ebiten.DrawImageOptions{}
	opTree.GeoM.Scale(s.state.menuScale, s.state.menuScale)
	opTree.GeoM.Translate(0, 3*float64(s.state.Config.btnSize)*s.state.menuScale)
	if s.state.isPreview {
		opTree.ColorM = invertedClr
	}

	opScene := &ebiten.DrawImageOptions{}
	opScene.GeoM.Scale(s.state.menuScale, s.state.menuScale)
	opScene.GeoM.Translate(0, 4*float64(s.state.Config.btnSize)*s.state.menuScale)

	switch s.state.hover {
	case 1:
		opPlay.ColorM = invertedClr
	case 2:
		opReset.ColorM = invertedClr
	case 3:
		opSpeed.ColorM = invertedClr
	case 5:
		opScene.ColorM = invertedClr
	}

	// Draw MENU
	if s.state.paused {
		screen.DrawImage(s.state.Config.pauseBtn, opPlay)
	} else {
		screen.DrawImage(s.state.Config.playBtn, opPlay)
	}
	screen.DrawImage(s.state.Config.resetBtn, opReset)
	screen.DrawImage(s.state.Config.speedBtn.SubImage(image.Rect(int(s.state.speed)*s.state.Config.btnSize, 0, (int(s.state.speed)+1)*s.state.Config.btnSize, s.state.Config.btnSize)).(*ebiten.Image), opSpeed)
	screen.DrawImage(s.state.Config.treeBtn, opTree)
	screen.DrawImage(s.state.Config.sceneBtn.SubImage(image.Rect(int(s.state.sceneNum)*s.state.Config.btnSize, 0, (int(s.state.sceneNum)+1)*s.state.Config.btnSize, s.state.Config.btnSize)).(*ebiten.Image), opScene)

}

func (s *SimulationScene) Load(state State, manager *stagehand.SceneManager[State]) {
	s.state = state
	s.sm = manager

	soilGrid := s.state.Levels[s.state.sceneNum].LayerByIdentifier("SoilType").IntGrid

	hum, rocks, rain := MakeSoilGrid(16, soilGrid)
	s.Board.Rain = rain
	s.Board.Rocks = rocks
	s.Board.Setup(hum)

	s.Preview.Setup(MakeColorGrid(16, soilGrid))
}

func (s *SimulationScene) Unload() State {
	return s.state
}

func (s *SimulationScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	s.state.menuScale = float64(outsideHeight) / 160
	if s.state.menuScale > 4 {
		s.state.menuScale = 4
	}
	return outsideWidth, outsideHeight
}

func MakeSoilGrid(size int, ldtkMap []*ldtkgo.Integer) (hum [][]mgl32.Vec2, rocks, rain [][]bool) {
	// Create Air Grid
	hum = util.MakeMatrix(size, mgl32.Vec2{0, 1})
	rocks = util.MakeMatrixBool(size)
	rain = util.MakeMatrixBool(size)
	for i, row := range hum {
		for j := range row {
			cell := ldtkMap[(j*16)+i].Value
			switch cell {
			case 1: // Air tile
				continue
			case 7: // Cloud Tile
				hum[i][j][0] = 1023
				rain[i][j] = true
			case 6: // Rock Tile
				hum[i][j][1] = math.MaxFloat32
				rocks[i][j] = true
			default:
				hum[i][j][1] = float32(math.Pow(5, float64(cell)))
			}
		}
	}
	return hum, rocks, rain
}

func MakeColorGrid(size int, ldtkMap []*ldtkgo.Integer) [][]color.Color {
	clrs := make([][]color.Color, size)
	for i := range clrs {
		clrs[i] = make([]color.Color, size)
		for j := range clrs[i] {
			clr := color.RGBA{0xff, 0xff, 0xff, 0xff}
			cell := ldtkMap[(j*16)+i].Value
			switch cell {
			case 7: // Cloud Tile
				clr = color.RGBA{0x12, 0x4e, 0x89, 0xff}
			case 6: // Rock Tile
				clr = color.RGBA{0x5A, 0x69, 0x88, 0xff}
			case 5: // Clay Tile
				clr = color.RGBA{0xBE, 0x4A, 0x2F, 0xff}
			case 4: // Sand Tile
				clr = color.RGBA{0xEA, 0xD4, 0xAA, 0xff}
			case 3: // hardSoil Tile
				clr = color.RGBA{0x55, 0x38, 0x29, 0xff}
			case 2: // looseSoil Tile
				clr = color.RGBA{0x27, 0x1F, 0x1E, 0xff}
			}
			clrs[i][j] = clr
		}
	}
	return clrs
}
