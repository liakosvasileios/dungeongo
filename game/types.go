package game

type Rect struct {
	X, Y, W, H float64
}

type WallType int

const (
	WallSolid WallType = iota
	WallDoor
)

func RectsCollide(a, b Rect) bool {
	return a.X < b.X+b.W &&
		a.X+a.W > b.X &&
		a.Y < b.Y+b.H &&
		a.Y+a.H > b.Y
}
