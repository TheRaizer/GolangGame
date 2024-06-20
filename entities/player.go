package entities

import (
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	objs.BaseGameObject
}

func NewPlayer(name string, initPos util.Vec2[float32]) Player {
	return Player{
		BaseGameObject: objs.NewBaseGameObject(name, initPos),
	}
}

func (player *Player) move(dt uint64) {
	player.UpdatePos(float32(dt)*0.1, float32(dt)*0.1)
}

func (player *Player) OnInit(surface *sdl.Surface) {
	// rect := sdl.Rect{X: int32(player.pos.X), Y: int32(player.pos.Y), W: 10, H: 10}
	// colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	// pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	// surface.FillRect(&rect, pixel)
}

func (player *Player) OnUpdate(dt uint64, surface *sdl.Surface) {
	player.move(dt)
}
