package objs

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	core.BaseGameObject

	rect  *sdl.Rect
	pixel uint32
}

func NewSprite(name string, initPos util.Vec2[float32], gameObjectStore core.GameObjectStore, width int32, height int32) Solid {
	return Solid{
		BaseGameObject: core.NewBaseGameObject(core.WALL_LAYER, name, initPos, gameObjectStore),
		rect:           &sdl.Rect{X: int32(initPos.X), Y: int32(initPos.Y), W: width, H: height},
	}
}

func (sprite *Sprite) OnInit(surface *sdl.Surface) {
	// TODO: using decoded RGB data, map to pixels
	// sdl.CreateRGBSurfaceFrom()
}

func (sprite *Sprite) OnUpdate(dt uint64, surface *sdl.Surface) {
}
