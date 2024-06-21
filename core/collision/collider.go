package collision

import (
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/util"
	datastructures "github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
)

type Collider struct {
	objs.BaseGameObject

	Rect              datastructures.Rect
	collisionEvents   []func(els []datastructures.QuadElement)
	collisionMediator CollisionSystemMediator
}

func NewCollider(name string, rect datastructures.Rect, collisionSys *CollisionSystem) *Collider {
	collider := Collider{
		Rect: rect,
		BaseGameObject: objs.NewBaseGameObject(
			name,
			util.Vec2[float32]{
				X: float32(rect.X),
				Y: float32(rect.Y),
			},
		),
		collisionMediator: collisionSys,
	}
	collisionSys.RegisterObject(&collider)

	return &collider
}

func (collider *Collider) UpdatePos(distX float32, distY float32) {
	collider.BaseGameObject.UpdatePos(distX, distY)

	newRect := datastructures.Rect{
		X: int32(collider.Pos.X),
		Y: int32(collider.Pos.Y),
		W: collider.Rect.W,
		H: collider.Rect.H,
	}

	collider.collisionMediator.UpdateCollider(collider.GetID(), collider.Rect, newRect)
	collider.Rect = newRect
}

func (collider *Collider) AddCollisionEvent(event func(els []datastructures.QuadElement)) {
	collider.collisionEvents = append(collider.collisionEvents, event)
}

func (collider *Collider) OnCollision(els []datastructures.QuadElement) {
	for _, event := range collider.collisionEvents {
		event(els)
	}
}
