package collision

import "github.com/TheRaizer/GolangGame/util/datastructures/quadtree"

type CollisionSystemMediator interface {
	UpdateCollider(id string, oldRect quadtree.Rect, newRect quadtree.Rect)
}

type CollisionSystem struct {
	tree      quadtree.QuadTree
	colliders map[string]*Collider
}

func NewCollisionSystem(globalRect quadtree.Rect) CollisionSystem {
	tree := quadtree.NewQuadTree(7, 5, globalRect)
	return CollisionSystem{
		tree:      &tree,
		colliders: make(map[string]*Collider),
	}
}

// Checks for collisions between registered colliders and calls their OnCollision callback
func (collisionSys *CollisionSystem) OnLoop() {
	for _, collider := range collisionSys.colliders {
		els := collisionSys.tree.Query(collider.Rect)
		collider.OnCollision(els)
	}
}

// call when updating a collider position or size
func (collisionSys *CollisionSystem) UpdateCollider(id string, oldRect quadtree.Rect, newRect quadtree.Rect) {
	collisionSys.tree.Remove(quadtree.QuadElement{Rect: oldRect, Id: id})
	collisionSys.tree.Insert(quadtree.QuadElement{Rect: newRect, Id: id})
}

func (collisionSys *CollisionSystem) RegisterObject(collider *Collider) {
	collisionSys.colliders[collider.GetID()] = collider
	collisionSys.tree.Insert(quadtree.QuadElement{Rect: collider.Rect, Id: collider.GetID()})
}

func (collisionSys *CollisionSystem) DeregisterObject(collider *Collider) {
	delete(collisionSys.colliders, collider.GetID())
	collisionSys.tree.Remove(quadtree.QuadElement{Rect: collider.Rect, Id: collider.GetID()})
}
