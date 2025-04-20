package game

import (
	"math/rand"
	"time"
)

type Room struct {
	X, Y int
	Map  *TileMap
}

const (
	cellSize    = 32
	gridWidth   = 800 / cellSize
	gridHeight  = 600 / cellSize
	roomCount   = 5 // number of random rooms
	minRoomSize = 4 // min cells per room dimension
	maxRoomSize = 8 // max cells per room dimension
)

// GenerateRoom creates a room with random rectangular subrooms (hollow) as obstacles
// while keeping the exits open.
func GenerateRoom(x, y int) *Room {
	const (
		wTiles = screenWidth / cellSize  // Width in tiles 25 for 800px / 32px-cells
		hTiles = screenHeight / cellSize // Height in tiles 18 for 600px / 32px-cells
	)

	m := NewTileMap(wTiles, hTiles)

	// 1) Border Walls
	for tx := 0; tx < wTiles; tx++ {
		m.Set(tx, 0, TileWall)
		m.Set(tx, hTiles-1, TileWall)
	}
	for ty := 0; ty < hTiles; ty++ {
		m.Set(0, ty, TileWall)
		m.Set(wTiles-1, ty, TileWall)
	}

	// 2) Door gaps (example: four‑way doors in center)
	m.Set(wTiles/2, 0, TileDoor)
	m.Set(wTiles/2, hTiles-1, TileDoor)
	m.Set(0, hTiles/2, TileDoor)
	m.Set(wTiles-1, hTiles/2, TileDoor)

	// Create a local RNG for obstacle placement
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Add random hollow rectangular rooms with entrances
	addRandomRooms(rng, m)

	return &Room{
		X:   x,
		Y:   y,
		Map: m,
	}
}

// addRandomRooms places hollow rectangular wall blocks in random interior areas
// each with a random entrance on one of the four sides, not at a corner.
func addRandomRooms(rng *rand.Rand, m *TileMap) {
	taken := make([][4]int, 0)
	for i := 0; i < roomCount; i++ {
		// Random room size in cells
		rw := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		rh := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize

		// Try to place without blocking corridors or overlapping
		for attempts := 0; attempts < 10; attempts++ {
			x0 := rng.Intn(gridWidth-2-rw) + 1
			y0 := rng.Intn(gridHeight-2-rh) + 1

			// Avoid central corridors
			if x0 <= gridWidth/2 && x0+rw > gridWidth/2 {
				continue
			}
			if y0 <= gridHeight/2 && y0+rh > gridHeight/2 {
				continue
			}

			// Check overlap
			overlap := false
			for _, t := range taken {
				if x0 < t[0]+t[2] && x0+rw > t[0] && y0 < t[1]+t[3] && y0+rh > t[1] {
					overlap = true
					break
				}
			}
			if overlap {
				continue
			}

			// Record this room
			taken = append(taken, [4]int{x0, y0, rw, rh})

			// Determine random entrance on one side, not at corners
			side := rng.Intn(4)
			var ex, ey int
			switch side {
			case 0: // top
				if rw > 2 {
					ex = x0 + 1 + rng.Intn(rw-2)
				} else {
					ex = x0 + rw/2
				}
				ey = y0
			case 1: // bottom
				if rw > 2 {
					ex = x0 + 1 + rng.Intn(rw-2)
				} else {
					ex = x0 + rw/2
				}
				ey = y0 + rh - 1
			case 2: // left
				ex = x0
				if rh > 2 {
					ey = y0 + 1 + rng.Intn(rh-2)
				} else {
					ey = y0 + rh/2
				}
			case 3: // right
				ex = x0 + rw - 1
				if rh > 2 {
					ey = y0 + 1 + rng.Intn(rh-2)
				} else {
					ey = y0 + rh/2
				}
			}

			// Draw perimeter walls, skipping entrance
			// Top and bottom
			for cx := 0; cx < rw; cx++ {
				x, y := x0+cx, y0
				if x != ex || y != ey {
					m.Set(x, y, TileWall)
				}
				x, y = x0+cx, y0+rh-1
				if x != ex || y != ey {
					m.Set(x, y, TileWall)
				}
			}
			// Left and right
			for cy := 1; cy < rh-1; cy++ {
				x, y := x0, y0+cy
				if x != ex || y != ey {
					m.Set(x, y, TileWall)
				}
				x, y = x0+rw-1, y0+cy
				if x != ex || y != ey {
					m.Set(x, y, TileWall)
				}
			}

			break
		}
	}
}
