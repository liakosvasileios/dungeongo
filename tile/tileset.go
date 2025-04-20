package tile

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var TileSprites map[TileType]*ebiten.Image

func init() {
	TileSprites = map[TileType]*ebiten.Image{
		TileFloor: makeColoredSprite(color.RGBA{40, 40, 40, 255}),
		TileWall:  makeColoredSprite(color.RGBA{0, 0, 255, 255}),
		TileDoor:  makeColoredSprite(color.RGBA{0, 255, 0, 255}),
	}
}

func makeColoredSprite(fill color.Color) *ebiten.Image {
	img := ebiten.NewImage(cellSize, cellSize)
	img.Fill(fill)

	out := color.RGBA{0, 0, 0, 255}
	w := cellSize - 1

	// horizontal edges
	for x := 0; x < cellSize; x++ {
		img.Set(x, 0, out) // top
		img.Set(x, w, out) // bottom
	}
	// vertical edges
	for y := 0; y < cellSize; y++ {
		img.Set(0, y, out) // left
		img.Set(w, y, out) // right
	}
	return img
}
