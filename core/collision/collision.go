package collision

import datastructures "github.com/TheRaizer/GolangGame/util/datastructures/quadtree"

type CollisionSystemMediator interface {
	UpdateCollider(id string, oldRect datastructures.Rect, newRect datastructures.Rect)
}

type CollisionSystem struct {
	tree      datastructures.QuadTree
	colliders map[string]*Collider
}

func NewCollisionSystem(globalRect datastructures.Rect) CollisionSystem {
	return CollisionSystem{
		tree:      datastructures.NewQuadTree(7, 5, globalRect),
		colliders: make(map[string]*Collider),
	}
}

func (collisionSys *CollisionSystem) OnLoop() {
	for _, collider := range collisionSys.colliders {
		els := collisionSys.tree.Query(collider.Rect)
		collider.OnCollision(els)
	}
}

// call when updating a collider position or size
func (collisionSys *CollisionSystem) UpdateCollider(id string, oldRect datastructures.Rect, newRect datastructures.Rect) {
	collisionSys.tree.Remove(datastructures.QuadElement{Rect: oldRect, Id: id})
	collisionSys.tree.Insert(datastructures.QuadElement{Rect: newRect, Id: id})
}

func (collisionSys *CollisionSystem) RegisterObject(collider *Collider) {
	collisionSys.colliders[collider.GetID()] = collider
	collisionSys.tree.Insert(datastructures.QuadElement{Rect: collider.Rect, Id: collider.GetID()})
}

func (collisionSys *CollisionSystem) DeregisterObject(collider *Collider) {
	delete(collisionSys.colliders, collider.GetID())
	collisionSys.tree.Remove(datastructures.QuadElement{Rect: collider.Rect, Id: collider.GetID()})
}
