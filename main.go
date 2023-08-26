package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/soil-demo/internal"
	"github.com/joelschutz/soil-demo/internal/boards"
	"github.com/joelschutz/stagehand"
	"github.com/solarlune/ldtkgo"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	ldtkProject *ldtkgo.Project
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Soil Demo")
	ebiten.SetWindowResizable(true)

	// Load Map
	// Load the LDtk Project
	ldtkProject, err := ldtkgo.Open("soil-demo.ldtk")
	if err != nil {
		log.Fatalf("Map Loading Fail: %s", err)
	}

	// Setup Simulation

	sm := stagehand.NewSceneManager[internal.State](&internal.SimulationScene{
		Board:   &boards.HumidityBoard{},
		Preview: &boards.EnumBoard{},
	}, internal.State{
		Levels: ldtkProject.Levels,
		Config: internal.NewConf(),
	})

	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}
