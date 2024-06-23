package collision

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	datastructures "github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
)

type Collider struct {
	core.BaseGameObject

	Rect              datastructures.Rect
	collisionEvents   []func(els []datastructures.QuadElement)
	collisionMediator CollisionSystemMediator
}

func NewCollider(name string,
	rect datastructures.Rect,
	collisionSys *CollisionSystem,
	collisionEvents []func(els []datastructures.QuadElement),
	gameObjectStore core.GameObjectStore,
) *Collider {
	collider := Collider{Rect: rect,
		BaseGameObject: core.NewBaseGameObject(
			name,
			util.Vec2[float32]{
				X: float32(rect.X),
				Y: float32(rect.Y),
			},
			gameObjectStore,
		),
		collisionMediator: collisionSys,
		collisionEvents:   collisionEvents,
	}
	collisionSys.RegisterObject(&collider)

	return &collider
}

// Update the position of the collider in the collision system
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

// executes all the collision events with the given collision elements
func (collider *Collider) OnCollision(els []datastructures.QuadElement) {
	for _, event := range collider.collisionEvents {
		event(els)
	}
}
