package core

import (
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type GameObject interface {
	OnInit(surface *sdl.Surface)
	OnUpdate(dt uint64, surface *sdl.Surface)
	GetID() string
	UpdatePos(distX float64, distY float64)
	AddChild(child GameObject)
	RemoveChild(id string)
}

type BaseGameObject struct {
	Pos      Vector
	children []GameObject
}

func (obj *BaseGameObject) UpdatePos(distX float64, distY float64) {
	obj.Pos.X += distX
	obj.Pos.Y += distY

	for _, child := range obj.children {
		child.UpdatePos(distX, distY)
	}
}

func (obj *BaseGameObject) AddChild(child GameObject) {
	obj.children = append(obj.children, child)
}

func (obj *BaseGameObject) RemoveChild(id string) {
	for i, child := range obj.children {
		if child.GetID() == id {
			obj.children = util.Slice[GameObject](obj.children).RemoveIdx(i)
			return
		}
	}
}

func (obj *BaseGameObject) OnUpdate(dt uint64, surface *sdl.Surface) {}
func (obj *BaseGameObject) OnInit(surface *sdl.Surface)              {}
