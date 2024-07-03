package objs

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Wall struct {
	core.BaseGameObject

	rect  *sdl.Rect
	pixel uint32
}

var colour = sdl.Color{R: 128, G: 128, B: 128, A: 255}

func NewWall(name string, initPos util.Vec2[float32], gameObjectStore core.GameObjectStore) Wall {
	return Wall{
		BaseGameObject: core.NewBaseGameObject(core.WALL_LAYER, name, initPos, gameObjectStore),
	}
}

func (wall *Wall) OnInit(surface *sdl.Surface) {
	wall.rect = &sdl.Rect{X: int32(wall.Pos.X), Y: int32(wall.Pos.Y), W: 32, H: 32}

	wall.pixel = sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(wall.rect, wall.pixel)
}

func (wall *Wall) OnUpdate(dt uint64, surface *sdl.Surface) {
	surface.FillRect(wall.rect, wall.pixel)
}
