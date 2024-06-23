package core

import (
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type GameObjectStore interface {
	AddGameObject(gameObject GameObject)
	RemoveGameObject(id string)
	GetGameObject(id string) GameObject
}

type GameObject interface {
	OnInit(surface *sdl.Surface)
	OnUpdate(dt uint64, surface *sdl.Surface)
	OnInput(event sdl.Event)
	GetID() string
	UpdatePos(distX float32, distY float32)
	AddChild(child GameObject)
	RemoveChild(id string)
	GetPos() util.Vec2[float32]
}

type BaseGameObject struct {
	Pos             util.Vec2[float32]
	name            string
	children        map[string]GameObject
	gameObjectStore GameObjectStore
}

func NewBaseGameObject(name string, pos util.Vec2[float32], gameObjectStore GameObjectStore) BaseGameObject {
	return BaseGameObject{
		pos,
		name,
		make(map[string]GameObject),
		gameObjectStore,
	}
}

// Updates the position of the current object by displacing it.
// This will also call UpdatePos on the children objects.
func (obj *BaseGameObject) UpdatePos(distX float32, distY float32) {
	obj.Pos.X += distX
	obj.Pos.Y += distY

	for _, child := range obj.children {
		child.UpdatePos(distX, distY)
	}
}

// Adds a gameobject as a child of the current gameobject.
// The child object will move with the parent, and is positioned relative to the parent.
func (obj *BaseGameObject) AddChild(child GameObject) {
	obj.children[child.GetID()] = child
	child.UpdatePos(obj.Pos.X, obj.Pos.Y)
	obj.gameObjectStore.AddGameObject(child)
}

func (obj *BaseGameObject) RemoveChild(id string) {
	delete(obj.children, id)
	obj.gameObjectStore.RemoveGameObject(id)
}

func (obj *BaseGameObject) GetID() string {
	return obj.name
}

func (obj *BaseGameObject) GetPos() util.Vec2[float32] {
	return obj.Pos
}

func (obj *BaseGameObject) OnUpdate(dt uint64, surface *sdl.Surface) {}
func (obj *BaseGameObject) OnInit(surface *sdl.Surface)              {}
func (obj *BaseGameObject) OnInput(event sdl.Event)                  {}
