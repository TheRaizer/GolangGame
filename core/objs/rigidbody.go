package objs

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
)

type RigidBody struct {
	core.BaseGameObject

	// TODO: should every rigid body have a reference to a collider? should i automatically add an on collison event?
	// TODO: on that collision event should it ComputePosOnCollision and update the parent with the new position?
	// outline requirements and functionality of a rigid body
	dir      util.Vec2[int8]
	velocity float32
}

func NewRigidBody(name string, initialVel float32, gameObjectStore core.GameObjectStore) RigidBody {
	return RigidBody{
		BaseGameObject: core.NewBaseGameObject(name, util.Vec2[float32]{}, gameObjectStore),
		dir:            util.Vec2[int8]{},
	}
}

// compute the position of the rect after being blocked by the otherRect
func (rb RigidBody) ComputePosOnCollision(rect quadtree.Rect, otherRect quadtree.Rect) util.Vec2[float32] {
}
