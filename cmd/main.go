package main

import (
	"dungeoncrawler/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dungeon Crawler")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
