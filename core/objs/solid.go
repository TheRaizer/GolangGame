package objs

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Solid struct {
	core.BaseGameObject

	rect  *sdl.Rect
	pixel uint32
}

var colour = sdl.Color{R: 128, G: 128, B: 128, A: 255}

func NewSolid(name string, initPos util.Vec2[float32], gameObjectStore core.GameObjectStore, width int32, height int32) Solid {
	return Solid{
		BaseGameObject: core.NewBaseGameObject(core.WALL_LAYER, name, initPos, gameObjectStore),
		rect:           &sdl.Rect{X: int32(initPos.X), Y: int32(initPos.Y), W: width, H: height},
	}
}

func (wall *Solid) OnInit(surface *sdl.Surface) {
	wall.pixel = sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(wall.rect, wall.pixel)
}

func (wall *Solid) OnUpdate(dt uint64, surface *sdl.Surface) {
	surface.FillRect(wall.rect, wall.pixel)
}
