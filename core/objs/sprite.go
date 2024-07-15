package objs

import (
	"unsafe"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	core.BaseGameObject

	rect   *sdl.Rect
	pixels [][][3]uint8 // a matrix of RGB values
}

func NewSprite(name string, initPos util.Vec2[float32], gameObjectStore core.GameObjectStore, width int32, height int32) Solid {
	return Solid{
		BaseGameObject: core.NewBaseGameObject(core.WALL_LAYER, name, initPos, gameObjectStore),
		rect:           &sdl.Rect{X: int32(initPos.X), Y: int32(initPos.Y), W: width, H: height},
	}
}

func (sprite *Sprite) OnInit(surface *sdl.Surface, renderer *sdl.Renderer) {
	// TODO: using decoded RGB data, map to pixels
	// sdl.CreateRGBSurfaceFrom()
	// texture, err := renderer.CreateTexture(
	// 	uint32(sdl.PIXELFORMAT_RGB888),
	// 	sdl.TEXTUREACCESS_STATIC,
	// 	int32(sprite.rect.W),
	// 	int32(sprite.rect.H),
	// )
	// util.CheckErr(err)
	// texture.Update(sprite.rect, sprite.pixels)
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888,
		sdl.TEXTUREACCESS_STATIC, 800, 600)
	util.CheckErr(err)
	defer texture.Destroy()

	pixels := make([]uint32, 800*600)

	texture.Update(nil, unsafe.Pointer(&pixels[0]), 800*3)
}

func (sprite *Sprite) OnUpdate(dt uint64, surface *sdl.Surface) {
}
