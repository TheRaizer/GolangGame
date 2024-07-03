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
	Speed        float32
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
		Speed:          initialVel,
		collider:       collider,
		collisionSys:   collisionSys,
		isContinous:    isContinous,
	}
}

func (rb *RigidBody) OnUpdate(dt uint64, surface *sdl.Surface) {
	if rb.Dir.X != 0 || rb.Dir.Y != 0 {
		distX := float32(rb.Dir.X) * float32(dt) * rb.Speed
		distY := float32(rb.Dir.Y) * float32(dt) * rb.Speed

		if rb.Dir.X == rb.Dir.Y {
			rb.detectCollision(distX*0.7071, distY*0.7071) // 0.7071 approx 1/sqrt(2) = magnitude of (1, 1) vector

		} else {
			rb.detectCollision(distX, distY)
		}
	}
}

// TODO: tackle the tunelling problem
func (rb *RigidBody) detectCollision(distX float32, distY float32) {
	// TODO: to tackle the tunelling problem instead of checking collisions with future position we can
	// detect collisions with the parallelogram from the current position to the future position.
	// if there are any collisions then deal with them accordingly
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
		if rb.Dir.X == rb.Dir.Y {
			rb.Parent().UpdatePos(distX, distY)
		} else {
			rb.Parent().UpdatePos(distX, distY)
		}
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
