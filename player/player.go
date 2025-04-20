package player

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSize  = 16
	playerSpeed = 2
)

type Player struct {
	X, Y  float64
	Image *ebiten.Image
}

// NewPlayer spawns a red square at (100,100).
func NewPlayer() *Player {
	img := ebiten.NewImage(playerSize, playerSize)
	img.Fill(color.RGBA{255, 0, 0, 255})
	return &Player{X: 100, Y: 100, Image: img}
}

// Update moves the player, handles collision against the TileMap and returns
// a direction (−1/0/1) if the player just walked through a border door tile.
func (p *Player) Update(m *TileMap) (doorDirX, doorDirY int) {
	// 1) accumulate intent
	vx, vy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		vx -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		vx += playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		vy -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		vy += playerSpeed
	}

	// 2) candidate new position
	nextX := p.X + vx
	nextY := p.Y + vy

	// 3) collision – sample the tile under the player’s centre
	cx := nextX + playerSize/2
	cy := nextY + playerSize/2
	tx := int(cx / float64(cellSize))
	ty := int(cy / float64(cellSize))

	tile := m.At(tx, ty)

	switch tile.Type {
	case TileWall:
		// blocked – stay where you were
		return 0, 0
	case TileDoor:
		// allow the move, then signal a room transition if this door is on a border
		p.X, p.Y = nextX, nextY
		switch {
		case tx == 0:
			return -1, 0 // left neighbour room
		case tx == m.W-1:
			return 1, 0 // right neighbour
		case ty == 0:
			return 0, -1 // up
		case ty == m.H-1:
			return 0, 1 // down
		}
	default: // TileFloor or anything else
		p.X, p.Y = nextX, nextY
	}

	return 0, 0
}

// Draw renders the player relative to the camera.
func (p *Player) Draw(screen *ebiten.Image, camX, camY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X-camX, p.Y-camY)
	screen.DrawImage(p.Image, op)
}

// Center returns the player’s logical centre coordinates.
func (p *Player) Center() (float64, float64) {
	return p.X + playerSize/2, p.Y + playerSize/2
}

// SetCenter warps the player so that its centre is at (cx,cy).
func (p *Player) SetCenter(cx, cy float64) {
	p.X = cx - playerSize/2
	p.Y = cy - playerSize/2
}
