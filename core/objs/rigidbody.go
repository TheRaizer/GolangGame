package objs

import (
	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
	"github.com/veandco/go-sdl2/sdl"
)

type RigidBody struct {
	core.BaseGameObject

	Dir      util.Vec2[int8]
	Velocity float32
	collider *collision.Collider
}

func NewRigidBody(layer int, name string, initialVel float32, gameObjectStore core.GameObjectStore, collider *collision.Collider) RigidBody {
	return RigidBody{
		BaseGameObject: core.NewBaseGameObject(layer, name, util.Vec2[float32]{}, gameObjectStore),
		Dir:            util.Vec2[int8]{},
		Velocity:       0,
		collider:       collider,
	}
}

func (rb *RigidBody) OnInit(surface *sdl.Surface) {
	rb.collider.AddCollisionEvent(rb.restrictParent)
}

func (rb *RigidBody) restrictParent(els []quadtree.QuadElement) {
	for _, el := range els {
		obj := rb.GameObjectStore.GetGameObject(el.Id)
		if obj.Layer() != rb.Layer() {
			if rb.collider.Rect.Right() > el.Rect.X {
				rb.Parent().UpdatePos(float32(el.Rect.X-rb.collider.Rect.Right()), 0)
			}
		}
	}
}
