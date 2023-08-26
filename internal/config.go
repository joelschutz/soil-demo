package internal

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed assets
	assets embed.FS
)

type Config struct {
	btnSize  int
	playBtn  *ebiten.Image
	pauseBtn *ebiten.Image
	resetBtn *ebiten.Image
	speedBtn *ebiten.Image
	treeBtn  *ebiten.Image
	sceneBtn *ebiten.Image
}

func NewConf() Config {
	Conf := Config{btnSize: 16}
	// Load Play Button
	buf, err := assets.ReadFile("assets/play-btn.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}
	Conf.playBtn = ebiten.NewImageFromImage(img)

	// Load Pause Button
	buf, err = assets.ReadFile("assets/pause-btn.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err = image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}
	Conf.pauseBtn = ebiten.NewImageFromImage(img)

	// Load Reset Button
	buf, err = assets.ReadFile("assets/reset-btn.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err = image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}
	Conf.resetBtn = ebiten.NewImageFromImage(img)

	// Load Speed Button
	buf, err = assets.ReadFile("assets/speed-btn.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err = image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	Conf.speedBtn = ebiten.NewImageFromImage(img)

	// Load Tree Button
	buf, err = assets.ReadFile("assets/tree-btn.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err = image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	Conf.treeBtn = ebiten.NewImageFromImage(img)

	// Load Scene Button
	buf, err = assets.ReadFile("assets/scene-btn.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err = image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	Conf.sceneBtn = ebiten.NewImageFromImage(img)

	return Conf
}
