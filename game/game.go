package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	camLerp      = 0.1
)

// Game is the top‑level Ebiten implementation that ties the
// player, world and camera together.
type Game struct {
	Player     *Player
	World      *World
	CamX, CamY float64

	// room‑transition animation state
	Transitioning bool
	TransitionT   float64
	TransitionDir int // 1 = closing, –1 = opening
	PendingRoomX  int
	PendingRoomY  int

	// debug
	ShowRays bool

	// full‑screen buffers
	BlackImage *ebiten.Image // for room‑transition bars
}

// NewGame sets up all top‑level singletons and render targets.
func NewGame() *Game {
	black := ebiten.NewImage(screenWidth, screenHeight)
	black.Fill(color.Black)

	return &Game{
		Player:     NewPlayer(),
		World:      NewWorld(),
		BlackImage: black,
	}
}

// Update runs once per tick.
func (g *Game) Update() error {
	// --- input -----------------------------------------------------------
	if inpututil.IsKeyJustPressed(ebiten.KeyR) { // toggle debug rays once per key‑press
		g.ShowRays = !g.ShowRays
	}

	// --- room‑transition animation --------------------------------------
	if g.Transitioning {
		g.TransitionT += 0.05 * float64(g.TransitionDir)
		if g.TransitionT >= 1.0 {
			// switch to the next room when the black bars fully close
			g.World.MoveTo(g.PendingRoomX, g.PendingRoomY)
			g.Player.SetCenter(384, 284) // spawn in the centre of the new room
			g.TransitionDir = -1         // start opening again
		} else if g.TransitionT <= 0 {
			// finished the open phase – back to normal play
			g.TransitionT = 0
			g.Transitioning = false
		}
		return nil // skip the rest of the update while animating
	}

	// --- camera ----------------------------------------------------------
	tx := g.Player.X + playerSize/2 - screenWidth/2
	ty := g.Player.Y + playerSize/2 - screenHeight/2
	g.CamX += (tx - g.CamX) * camLerp
	g.CamY += (ty - g.CamY) * camLerp

	// --- player & possible room change ----------------------------------
	dx, dy := g.Player.Update(g.World.Current.Map) // uses TileMap collision now
	if (dx != 0 || dy != 0) && !g.Transitioning {
		// kick off the transition animation toward the neighbour room
		g.Transitioning = true
		g.TransitionT = 0
		g.TransitionDir = 1
		g.PendingRoomX = g.World.Current.X + dx
		g.PendingRoomY = g.World.Current.Y + dy
	}

	return nil
}

// Draw renders the current frame.
func (g *Game) Draw(screen *ebiten.Image) {
	room := g.World.Current

	// // 1) draw every visible tile -----------------------------------------
	// for ty := 0; ty < room.Map.H; ty++ {
	// 	for tx := 0; tx < room.Map.W; tx++ {
	// 		tile := room.Map.At(tx, ty)
	// 		sprite := TileSprites[tile.Type]
	// 		op := &ebiten.DrawImageOptions{}
	// 		op.GeoM.Translate(float64(tx*cellSize)-g.CamX, float64(ty*cellSize)-g.CamY)
	// 		screen.DrawImage(sprite, op)
	// 	}
	// }

	// 2) dynamic light / FOV / RayCasting ---------------------------------
	px, py := g.Player.Center()
	segs := MapSegments(room.Map) // converts solid tiles to line segments

	visible := VisibleTiles(px, py, room.Map)

	for ty := 0; ty < room.Map.H; ty++ {
		for tx := 0; tx < room.Map.W; tx++ {
			if !visible[ty][tx] { // Only draw visible tiles
				continue
			}
			tile := room.Map.At(tx, ty)
			sprite := TileSprites[tile.Type]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tx*cellSize)-g.CamX,
				float64(ty*cellSize)-g.CamY)
			screen.DrawImage(sprite, op)
		}
	}

	// Debug
	DrawLight(screen, px, py, segs, g.CamX, g.CamY, g.ShowRays)

	// 3) draw the player --------------------------------------------------
	g.Player.Draw(screen, g.CamX, g.CamY)

	// 4) transition bars --------------------------------------------------
	if g.Transitioning {
		t := g.TransitionT
		if t > 1 {
			t = 1
		}
		w := int(float64(screenWidth) * t / 2)

		// left bar
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(w), float64(screenHeight))
		screen.DrawImage(g.BlackImage, op)

		// right bar
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(w), float64(screenHeight))
		op.GeoM.Translate(float64(screenWidth-w), 0)
		screen.DrawImage(g.BlackImage, op)
	}
}

// Layout tells Ebiten the logical screen size.
func (g *Game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }
