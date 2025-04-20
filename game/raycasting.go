package game

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type line struct{ X1, Y1, X2, Y2 float64 }

func (l line) angle() float64 { return math.Atan2(l.Y2-l.Y1, l.X2-l.X1) }

func newRay(x, y, length, angle float64) line {
	return line{x, y, x + length*math.Cos(angle), y + length*math.Sin(angle)}
}

func intersection(a, b line) (x, y float64, ok bool) {
	d := (a.X1-a.X2)*(b.Y1-b.Y2) - (a.Y1-a.Y2)*(b.X1-b.X2)
	if d == 0 {
		return
	}
	t := ((a.X1-b.X1)*(b.Y1-b.Y2) - (a.Y1-b.Y1)*(b.X1-b.X2)) / d
	u := -((a.X1-a.X2)*(a.Y1-b.Y1) - (a.Y1-a.Y2)*(a.X1-b.X1)) / d
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return
	}
	x = a.X1 + t*(a.X2-a.X1)
	y = a.Y1 + t*(a.Y2-a.Y1)
	ok = true
	return
}

/* ---------- convert solid tiles to segment list ---------- */

func tileRectToLines(tx, ty int) []line {
	x := float64(tx * cellSize)
	y := float64(ty * cellSize)
	w := float64(cellSize)
	h := float64(cellSize)
	return []line{
		{x, y, x + w, y},
		{x + w, y, x + w, y + h},
		{x + w, y + h, x, y + h},
		{x, y + h, x, y},
	}
}

// MapSegments returns all wall edges as line segments.
func MapSegments(m *TileMap) []line {
	out := make([]line, 0)
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if m.At(x, y).Type == TileWall {
				out = append(out, tileRectToLines(x, y)...)
			}
		}
	}
	return out
}

/* ---------- ray casting ---------- */

func castRays(px, py float64, segs []line) []line {
	const rayLen = 2000
	rays := make([]line, 0, len(segs)*2)

	for _, s := range segs {
		for _, p := range [][2]float64{{s.X1, s.Y1}, {s.X2, s.Y2}} {
			base := math.Atan2(p[1]-py, p[0]-px)
			for _, off := range []float64{-0.0005, 0.0005} {
				r := newRay(px, py, rayLen, base+off)

				hx, hy, best := 0.0, 0.0, math.Inf(1)
				for _, seg := range segs {
					if ix, iy, ok := intersection(r, seg); ok {
						if d := (ix-px)*(ix-px) + (iy-py)*(iy-py); d < best {
							hx, hy, best = ix, iy, d
						}
					}
				}
				rays = append(rays, line{px, py, hx, hy})
			}
		}
	}
	sort.Slice(rays, func(i, j int) bool { return rays[i].angle() < rays[j].angle() })
	return rays
}

/* ---------- drawing helpers ---------- */

var (
	shadowImg   = ebiten.NewImage(screenWidth, screenHeight)
	triangleImg = func() *ebiten.Image { img := ebiten.NewImage(1, 1); img.Fill(color.White); return img }()
)

func triangleVerts(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x2), DstY: float32(y2), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

// DrawLight darkens the screen and carves a visibility polygon.
func DrawLight(screen *ebiten.Image, px, py float64, segs []line,
	camX, camY float64, debug bool) {

	shadowImg.Clear()
	shadowImg.Fill(color.Black)

	rays := castRays(px, py, segs)
	opt := &ebiten.DrawTrianglesOptions{Blend: ebiten.BlendDestinationOut}

	for i, r := range rays {
		n := rays[(i+1)%len(rays)]
		v := triangleVerts(px-camX, py-camY, n.X2-camX, n.Y2-camY, r.X2-camX, r.Y2-camY)
		shadowImg.DrawTriangles(v, []uint16{0, 1, 2}, triangleImg, opt)
	}

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(shadowImg, op)

	if debug {
		for _, r := range rays {
			vector.StrokeLine(screen,
				float32(r.X1-camX), float32(r.Y1-camY),
				float32(r.X2-camX), float32(r.Y2-camY),
				1, color.RGBA{255, 255, 0, 160}, true)
		}
	}
}

const viewRadius = 1000.0

var viewRadius2 = viewRadius * viewRadius

// hasLOS returns true when the straight line from (px,py) to the centre
// of tile (tx,ty) reaches that tile without passing through a wall tile.
func hasLOS(px, py float64, tx, ty int, m *TileMap) bool {
	// Bresenham on tile coordinates
	x0, y0 := int(px)/cellSize, int(py)/cellSize
	x1, y1 := tx, ty
	dx, dy := abs(x1-x0), abs(y1-y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		if x0 == x1 && y0 == y1 {
			return true
		}

		// skip the player’s own tile; block on any other wall
		if !(x0 == int(px)/cellSize && y0 == int(py)/cellSize) &&
			m.At(x0, y0).Type == TileWall {
			return false
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// VisibleTiles builds the visibility mask using the LOS test above.
func VisibleTiles(px, py float64, m *TileMap) [][]bool {
	vis := make([][]bool, m.H)
	for y := range vis {
		vis[y] = make([]bool, m.W)
	}

	pcx, pcy := int(px)/cellSize, int(py)/cellSize
	max := int(math.Floor(viewRadius/float64(cellSize))) + 1

	for dy := -max; dy <= max; dy++ {
		for dx := -max; dx <= max; dx++ {
			tx, ty := pcx+dx, pcy+dy
			if !m.InBounds(tx, ty) {
				continue
			}

			cx := float64(tx*cellSize + cellSize/2)
			cy := float64(ty*cellSize + cellSize/2)
			if (cx-px)*(cx-px)+(cy-py)*(cy-py) > viewRadius2 {
				continue
			}

			if hasLOS(px, py, tx, ty, m) {
				vis[ty][tx] = true
			}
		}
	}
	return vis
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
