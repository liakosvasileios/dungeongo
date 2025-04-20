package game

type TileType uint8

const (
	TileFloor TileType = iota
	TileWall
	TileDoor
)

type Tile struct {
	Type TileType
}

type TileMap struct {
	W, H int      // Dimensions in tiles
	Data [][]Tile // Data[y][x]
}

func NewTileMap(w, h int) *TileMap {
	d := make([][]Tile, h)
	for y := range d {
		d[y] = make([]Tile, w)
		for x := range d[y] {
			d[y][x] = Tile{Type: TileFloor}
		}
	}
	return &TileMap{W: w, H: h, Data: d}
}

func (m *TileMap) InBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < m.W && y < m.H
}

func (m *TileMap) At(x, y int) Tile {
	if !m.InBounds(x, y) {
		return Tile{Type: TileWall} // We make out of bounds solid
	}
	return m.Data[y][x]
}

func (m *TileMap) Set(x, y int, t TileType) {
	if m.InBounds(x, y) {
		m.Data[y][x].Type = t
	}
}
