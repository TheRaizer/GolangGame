package display

const (
	FLOOR = iota
	WALL
)

var Shade = [2]string{" ", "█"}

// map each shade to a tile
var TileMap = map[string]uint8{Shade[0]: WALL, Shade[1]: FLOOR}
