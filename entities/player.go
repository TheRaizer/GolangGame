package entities

import (
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	objs.BaseGameObject

	rect  *sdl.Rect
	pixel uint32
	dir   util.Vec2[int8]
	speed uint8
}

var colour = sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple

func NewPlayer(name string, initPos util.Vec2[float32]) Player {
	return Player{
		BaseGameObject: objs.NewBaseGameObject(name, initPos),
	}
}

func (player *Player) move(dt uint64) {
	if player.dir.X != 0 || player.dir.Y != 0 {
		player.UpdatePos(float32(player.dir.X)*float32(dt)*0.1, float32(player.dir.Y)*float32(dt)*0.1)
		player.rect.X = int32(player.Pos.X)
		player.rect.Y = int32(player.Pos.Y)
	}
}

func (player *Player) OnInit(surface *sdl.Surface) {
	player.rect = &sdl.Rect{X: int32(player.Pos.X), Y: int32(player.Pos.Y), W: 32, H: 32}

	player.pixel = sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(player.rect, player.pixel)
}

func (player *Player) OnUpdate(dt uint64, surface *sdl.Surface) {
	surface.FillRect(player.rect, 0)
	player.move(dt)
	surface.FillRect(player.rect, player.pixel)
}

func (player *Player) OnInput(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		if t.State == sdl.PRESSED {
			if t.Keysym.Sym == sdl.K_LEFT {
				player.dir.X = -1
			} else if t.Keysym.Sym == sdl.K_RIGHT {
				player.dir.X = 1
			}
			if t.Keysym.Sym == sdl.K_UP {
				player.dir.Y = -1
			} else if t.Keysym.Sym == sdl.K_DOWN {
				player.dir.Y = 1
			}
		} else if t.State == sdl.RELEASED {
			if (t.Keysym.Sym == sdl.K_LEFT && player.dir.X == -1) || (t.Keysym.Sym == sdl.K_RIGHT && player.dir.X == 1) {
				player.dir.X = 0
			}
			if (t.Keysym.Sym == sdl.K_UP && player.dir.Y == -1) || (t.Keysym.Sym == sdl.K_DOWN && player.dir.Y == 1) {
				player.dir.Y = 0
			}
		}
		break
	}
}
