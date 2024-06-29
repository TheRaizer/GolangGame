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

	Dir          util.Vec2[int8]
	Velocity     float32
	collider     *collision.Collider
	collisionSys collision.CollisionSystemMediator
	restriction  uint8 // 0 = horizontal, 1 = vertical, 2 = none
	isContinous  bool
}

func NewRigidBody(
	layer int,
	name string,
	initialVel float32,
	gameObjectStore core.GameObjectStore,
	collider *collision.Collider,
	collisionSys collision.CollisionSystemMediator,
	isContinous bool,
) RigidBody {
	return RigidBody{
		BaseGameObject: core.NewBaseGameObject(layer, name, util.Vec2[float32]{}, gameObjectStore),
		Dir:            util.Vec2[int8]{},
		Velocity:       initialVel,
		collider:       collider,
		collisionSys:   collisionSys,
		isContinous:    isContinous,
	}
}

func (rb *RigidBody) OnInit(surface *sdl.Surface) {
	if !rb.isContinous {
		// add collision event that restricts parent on the frame update (discrete)
		rb.collider.AddCollisionEvent(rb.restrictParent)
	}
}

func (rb *RigidBody) OnUpdate(dt uint64, surface *sdl.Surface) {
	if rb.Dir.X != 0 || rb.Dir.Y != 0 {
		distX := float32(rb.Dir.X) * float32(dt) * rb.Velocity
		distY := float32(rb.Dir.Y) * float32(dt) * rb.Velocity

		if rb.isContinous {
			rb.continousCollisionDetect(distX, distY)
		} else {
			// move on the restriction event registered OnInit
			rb.moveOnRestriction(distX, distY)
		}
	}
}

func (rb *RigidBody) continousCollisionDetect(distX float32, distY float32) {
	newRect := quadtree.Rect{
		X: int32(rb.collider.Pos.X + distX),
		Y: int32(rb.collider.Pos.Y + distY),
		W: rb.collider.Rect.W,
		H: rb.collider.Rect.H,
	}
	// precompute possible collision for dynamic movement
	els := rb.collisionSys.DetectCollisions(newRect)
	rb.restrictParent(els)
	rb.moveOnRestriction(distX, distY)
}

func (rb *RigidBody) moveOnRestriction(distX float32, distY float32) {
	if rb.restriction == 0 {
		rb.Parent().UpdatePos(0, distY)
	}
	if rb.restriction == 1 {
		rb.Parent().UpdatePos(distX, 0)
	}

	if rb.restriction == 2 {
		rb.Parent().UpdatePos(distX, distY)
	}
}

// restricts the parent object according to colliding elements and updates rb.restriction with corresponding restriction
func (rb *RigidBody) restrictParent(els []quadtree.QuadElement) {
	for _, el := range els {
		obj := rb.GameObjectStore.GetGameObject(el.Id)
		if obj.Layer() != rb.Layer() {
			overlapLeft := rb.collider.Rect.Right() - el.Rect.X
			overlapRight := el.Rect.Right() - rb.collider.Rect.X
			overlapTop := rb.collider.Rect.Bottom() - el.Rect.Y
			overlapBottom := el.Rect.Bottom() - rb.collider.Rect.Y

			if min(overlapLeft, overlapRight) < min(overlapTop, overlapBottom) {
				if overlapLeft < overlapRight {
					rb.Parent().UpdatePos(-float32(overlapLeft), 0)
				} else {
					rb.Parent().UpdatePos(float32(overlapRight), 0)
				}
				rb.restriction = 0
				return
			} else {
				if overlapTop < overlapBottom {
					rb.Parent().UpdatePos(0, -float32(overlapTop))
				} else {
					rb.Parent().UpdatePos(0, float32(overlapBottom))
				}
				rb.restriction = 1
				return
			}
		}
	}
	rb.restriction = 2
}
