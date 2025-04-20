package game

import "fmt"

type World struct {
	Rooms   map[string]*Room
	RoomX   int
	RoomY   int
	Current *Room
}

func NewWorld() *World {
	w := &World{
		Rooms: make(map[string]*Room),
	}
	w.loadRoom(0, 0)
	return w
}

func (w *World) loadRoom(x, y int) {
	key := fmt.Sprintf("%d, %d", x, y)
	if _, ok := w.Rooms[key]; !ok {
		w.Rooms[key] = GenerateRoom(x, y)
	}
	w.RoomX = x
	w.RoomY = y
	w.Current = w.Rooms[key]
}

func (w *World) MoveTo(dx, dy int) {
	w.loadRoom(w.RoomX+dx, w.RoomY+dy)
}
