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
	ID() string
	UpdatePos(distX float32, distY float32)
	AddChild(child GameObject)
	RemoveChild(id string)
	Parent() GameObject
	SetParent(parent GameObject)
	Layer() int
}

type BaseGameObject struct {
	layer           int
	parent          GameObject
	Pos             util.Vec2[float32]
	name            string
	children        map[string]GameObject
	GameObjectStore GameObjectStore
}

func NewBaseGameObject(layer int, name string, pos util.Vec2[float32], gameObjectStore GameObjectStore) BaseGameObject {
	return BaseGameObject{
		layer,
		nil,
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
	obj.children[child.ID()] = child
	child.UpdatePos(obj.Pos.X, obj.Pos.Y)
	obj.GameObjectStore.AddGameObject(child)
	child.SetParent(obj)
}

func (obj *BaseGameObject) RemoveChild(id string) {
	child := obj.children[id]
	child.SetParent(nil)
	delete(obj.children, id)
	obj.GameObjectStore.RemoveGameObject(id)
}

func (obj *BaseGameObject) Parent() GameObject {
	return obj.parent
}

func (obj *BaseGameObject) SetParent(parent GameObject) {
	obj.parent = parent
}

func (obj *BaseGameObject) ID() string {
	return obj.name
}

func (obj *BaseGameObject) Layer() int {
	return obj.layer
}

func (obj *BaseGameObject) OnUpdate(dt uint64, surface *sdl.Surface) {}
func (obj *BaseGameObject) OnInit(surface *sdl.Surface)              {}
func (obj *BaseGameObject) OnInput(event sdl.Event)                  {}
