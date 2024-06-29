package entities

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	core.BaseGameObject

	rect    *sdl.Rect
	pixel   uint32
	rb      *objs.RigidBody
	surface *sdl.Surface
}

var colour = sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple

func NewPlayer(name string, initPos util.Vec2[float32], gameObjectStore core.GameObjectStore, rb *objs.RigidBody) Player {
	return Player{
		BaseGameObject: core.NewBaseGameObject(core.PLAYER_LAYER, name, initPos, gameObjectStore),
		rb:             rb,
	}
}

func (player *Player) OnInit(surface *sdl.Surface) {
	player.surface = surface
	player.rect = &sdl.Rect{X: int32(player.Pos.X), Y: int32(player.Pos.Y), W: 32, H: 32}

	player.pixel = sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(player.rect, player.pixel)
}

func (player *Player) UpdatePos(distX float32, distY float32) {
	player.surface.FillRect(player.rect, 0)

	player.BaseGameObject.UpdatePos(distX, distY)
	player.rect.X = int32(player.Pos.X)
	player.rect.Y = int32(player.Pos.Y)

	player.surface.FillRect(player.rect, player.pixel)

}

func (player *Player) OnInput(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		if t.State == sdl.PRESSED {
			if t.Keysym.Sym == sdl.K_LEFT {
				player.rb.Dir.X = -1
			} else if t.Keysym.Sym == sdl.K_RIGHT {
				player.rb.Dir.X = 1
			}
			if t.Keysym.Sym == sdl.K_UP {
				player.rb.Dir.Y = -1
			} else if t.Keysym.Sym == sdl.K_DOWN {
				player.rb.Dir.Y = 1
			}
		} else if t.State == sdl.RELEASED {
			if (t.Keysym.Sym == sdl.K_LEFT && player.rb.Dir.X == -1) || (t.Keysym.Sym == sdl.K_RIGHT && player.rb.Dir.X == 1) {
				player.rb.Dir.X = 0
			}
			if (t.Keysym.Sym == sdl.K_UP && player.rb.Dir.Y == -1) || (t.Keysym.Sym == sdl.K_DOWN && player.rb.Dir.Y == 1) {
				player.rb.Dir.Y = 0
			}
		}
		break
	}
}

func (player *Player) AddChild(child core.GameObject) {
	player.BaseGameObject.AddChild(child)
	child.SetParent(player)
}
