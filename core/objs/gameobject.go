package objs

import (
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type GameObject interface {
	OnInit(surface *sdl.Surface)
	OnUpdate(dt uint64, surface *sdl.Surface)
	OnInput(event sdl.Event)
	GetID() string
	UpdatePos(distX float32, distY float32)
	AddChild(child GameObject)
	RemoveChild(id string)
}

type BaseGameObject struct {
	Pos      util.Vec2[float32]
	name     string
	children map[string]GameObject
}

func NewBaseGameObject(name string, pos util.Vec2[float32]) BaseGameObject {
	return BaseGameObject{
		pos,
		name,
		make(map[string]GameObject),
	}
}

func (obj *BaseGameObject) UpdatePos(distX float32, distY float32) {
	obj.Pos.X += distX
	obj.Pos.Y += distY

	for _, child := range obj.children {
		child.UpdatePos(distX, distY)
	}
}

func (obj *BaseGameObject) AddChild(child GameObject) {
	obj.children[child.GetID()] = child
	child.UpdatePos(obj.Pos.X, obj.Pos.Y)
}

func (obj *BaseGameObject) RemoveChild(id string) {
	delete(obj.children, id)
}

func (obj *BaseGameObject) GetID() string {
	return obj.name
}

func (obj *BaseGameObject) OnUpdate(dt uint64, surface *sdl.Surface) {}
func (obj *BaseGameObject) OnInit(surface *sdl.Surface)              {}
func (obj *BaseGameObject) OnInput(event sdl.Event)                  {}
