package entities

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/display"
)

type Player struct {
	core.BaseGameObject
	pos core.Vector
}

func NewPlayer(initPos core.Vector) *Player {
	return &Player{pos: initPos, BaseGameObject: core.BaseGameObject{}}
}

func (player *Player) canMove(p display.Pixel) bool {
	return true
}

func (player *Player) move(dt uint64) {
	player.pos.X += float64(dt) * 0.1
}

func (player *Player) OnInit() {

}

func (player *Player) OnUpdate(dt uint64) {
	player.move(dt)
}
