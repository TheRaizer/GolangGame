package collision

import (
	"testing"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockQueryResult = []quadtree.QuadElement{
	{
		Id:   "mockId1",
		Rect: quadtree.Rect{X: 5, Y: 2, W: 10, H: 11},
	},
	{
		Id:   "mockId2",
		Rect: quadtree.Rect{X: 5, Y: 12, W: 13, H: 11},
	},
}

type QuadTreeMock struct {
	quadtree.BaseQuadTree
	mock.Mock
}

func (tree *QuadTreeMock) Insert(el quadtree.QuadElement) {
	tree.Called(el)
}

func (tree *QuadTreeMock) Remove(el quadtree.QuadElement) {
	tree.Called(el)
}

func (tree *QuadTreeMock) Query(hitbox quadtree.Rect) []quadtree.QuadElement {
	tree.Called(hitbox)
	return mockQueryResult
}

func TestRegisterObjectShouldStoreObject(t *testing.T) {
	collisionSys := NewCollisionSystem(quadtree.Rect{X: 0, Y: 0, W: 50, H: 50})

	expectedId := "idhere"
	collider := NewCollider(
		0,
		expectedId,
		quadtree.Rect{X: 10, Y: 10, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		[]func(els []quadtree.QuadElement){},
		nil,
	)

	collisionSys.RegisterObject(collider)

	for id, obj := range collisionSys.colliders {
		require.Equal(t, expectedId, id)
		require.Equal(t, *obj, *collider)
	}
}

func TestRegisterObjectShouldInsertIntoQuadTree(t *testing.T) {
	mockTree := &QuadTreeMock{}
	collisionSys := &CollisionSystem{
		tree:      mockTree,
		colliders: make(map[string]*Collider),
	}

	expectedId := "idhere"
	collider := &Collider{
		BaseGameObject: core.NewBaseGameObject(
			0,
			expectedId,
			util.Vec2[float32]{},
			nil,
		),
		Rect: quadtree.Rect{
			X: 1,
			Y: 10,
			W: 11,
			H: 3,
		},
	}

	mockTree.On("Insert", quadtree.QuadElement{Id: expectedId, Rect: collider.Rect})
	collisionSys.RegisterObject(collider)
	mockTree.AssertExpectations(t)
}

func TestDeregisterObjectShouldRemoveObjectFromQuadTree(t *testing.T) {
	mockTree := &QuadTreeMock{}
	expectedId := "id"
	colliderToRemove := &Collider{
		BaseGameObject: core.NewBaseGameObject(
			0,
			expectedId,
			util.Vec2[float32]{},
			nil,
		),
		Rect: quadtree.Rect{
			X: 1,
			Y: 10,
			W: 11,
			H: 3},
	}
	colliders := map[string]*Collider{
		expectedId: colliderToRemove,
	}
	collisionSys := &CollisionSystem{
		tree:      mockTree,
		colliders: colliders,
	}

	mockTree.On("Remove", quadtree.QuadElement{Id: expectedId, Rect: colliderToRemove.Rect})
	collisionSys.DeregisterObject(colliderToRemove)
	mockTree.AssertExpectations(t)
}

func TestDeregisterObjectShouldRemoveObjFromMap(t *testing.T) {
	collisionSys := NewCollisionSystem(quadtree.Rect{X: 0, Y: 0, W: 50, H: 50})
	expectedId := "idhere"
	collider := NewCollider(
		0,
		expectedId,
		quadtree.Rect{X: 10, Y: 10, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		[]func(els []quadtree.QuadElement){},
		nil,
	)

	collisionSys.colliders = map[string]*Collider{
		expectedId: collider,
	}

	collisionSys.DeregisterObject(collider)
	require.Nil(t, collisionSys.colliders[expectedId])
}

func TestUpdateColliderShouldReinsertColliderIntoQuadTree(t *testing.T) {
	mockTree := &QuadTreeMock{}
	expectedId := "id"
	collisionSys := &CollisionSystem{
		tree:      mockTree,
		colliders: make(map[string]*Collider),
	}
	oldRect := quadtree.Rect{X: 0, Y: 0, W: 32, H: 32}
	newRect := quadtree.Rect{X: 10, Y: 10, W: 32, H: 32}

	mockTree.On("Remove", quadtree.QuadElement{Id: expectedId, Rect: oldRect})
	mockTree.On("Insert", quadtree.QuadElement{Id: expectedId, Rect: newRect})

	collisionSys.UpdateCollider(expectedId, oldRect, newRect)
	mockTree.AssertExpectations(t)

}

func TestOnLoop(t *testing.T) {
	mockTree := &QuadTreeMock{}
	collider1 := &Collider{
		BaseGameObject: core.NewBaseGameObject(
			0,
			"id1",
			util.Vec2[float32]{},
			nil,
		),
		Rect: quadtree.Rect{
			X: 1,
			Y: 10,
			W: 11,
			H: 3,
		},
		collisionEvents: []func(els []quadtree.QuadElement){
			func(els []quadtree.QuadElement) {
				// expect els to be returned from Query call
				require.ElementsMatch(t, mockQueryResult, els)
			},
		},
	}
	collider2 := &Collider{
		BaseGameObject: core.NewBaseGameObject(
			0,
			"id1",
			util.Vec2[float32]{},
			nil,
		),
		Rect: quadtree.Rect{
			X: 3,
			Y: 33,
			W: 7,
			H: 4,
		},
		collisionEvents: []func(els []quadtree.QuadElement){
			func(els []quadtree.QuadElement) {
				require.ElementsMatch(t, mockQueryResult, els)
			},
		},
	}

	collisionSys := &CollisionSystem{
		tree: mockTree,
		colliders: map[string]*Collider{
			"id1": collider1,
			"id2": collider2,
		},
	}

	// should query on each collider
	mockTree.On("Query", collider1.Rect)
	mockTree.On("Query", collider2.Rect)

	collisionSys.OnLoop()
	mockTree.AssertExpectations(t)
}
