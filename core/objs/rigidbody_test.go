package objs

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCollisionSystem struct {
	mock.Mock
	elementsToDetect []quadtree.QuadElement
}

func (collisionSys *MockCollisionSystem) DetectCollisions(rect quadtree.Rect) []quadtree.QuadElement {
	collisionSys.Called(rect)
	return collisionSys.elementsToDetect
}

func (collisionSys *MockCollisionSystem) UpdateCollider(id string, oldRect quadtree.Rect, newRect quadtree.Rect) {
}

func (collisionSys *MockCollisionSystem) RegisterObject(obj *collision.Collider) {}

func (collisionSys *MockCollisionSystem) DeregisterObject(obj *collision.Collider) {}

func (collisionSys *MockCollisionSystem) OnLoop() {}

type MockGameObjectStore struct {
	mock.Mock
}

func (store *MockGameObjectStore) AddGameObject(gameObject core.GameObject) {
	store.Called(gameObject)
}

func (store *MockGameObjectStore) RemoveGameObject(id string) {
	store.Called(id)
}

// any object that will restrict movement will have layer 0
// (so any rb that does not have layer 0 will be blocked by this object)
func (store *MockGameObjectStore) GetGameObject(id string) core.GameObject {
	store.Called(id)
	mockObj := core.NewBaseGameObject(0, id, util.Vec2[float32]{}, store)
	return &mockObj
}

func TestOnUpdateShouldMoveParentWhenNotRestricted(t *testing.T) {
	type TestInput struct {
		dt            uint64
		dir           util.Vec2[int8]
		velocity      float32
		collisionRect quadtree.Rect
	}

	type TestExpected struct {
		// the new position of the parent element
		pos util.Vec2[float32]
		// the rectangle used to detect collisions in the future position
		rect quadtree.Rect
	}

	const NAME string = "should update position with %+v"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input)
	}

	cases := []util.TestCase[TestInput, TestExpected]{
		{
			Name: getName,
			Input: TestInput{
				dt:            10,
				dir:           util.Vec2[int8]{X: 1, Y: 1},
				velocity:      1,
				collisionRect: quadtree.Rect{X: 0, Y: 0, W: 5, H: 5},
			},
			Expected: TestExpected{
				// account for diagonal speed reduction
				pos: util.Vec2[float32]{X: float32(10) * 0.7071, Y: float32(10) * 0.7071},
				rect: quadtree.Rect{
					X: 7, Y: 7, W: 5, H: 5, // X = Y = 10*0.7071 truncated to int
				},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				dt:            2,
				dir:           util.Vec2[int8]{X: 1, Y: 0},
				velocity:      2,
				collisionRect: quadtree.Rect{X: 5, Y: 5, W: 5, H: 5},
			},
			Expected: TestExpected{
				pos: util.Vec2[float32]{X: 5 + float32(2)*2, Y: 5},
				rect: quadtree.Rect{
					X: 9, Y: 5, W: 5, H: 5,
				},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				dt:            1,
				dir:           util.Vec2[int8]{X: 0, Y: 1},
				velocity:      1,
				collisionRect: quadtree.Rect{X: 3, Y: 4, W: 15, H: 15},
			},
			Expected: TestExpected{
				pos: util.Vec2[float32]{X: 3, Y: 5},
				rect: quadtree.Rect{
					X: 3, Y: 5, W: 15, H: 15,
				},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				dt:            1,
				dir:           util.Vec2[int8]{X: 0, Y: 1},
				velocity:      1,
				collisionRect: quadtree.Rect{X: 3, Y: 4, W: 15, H: 15},
			},
			Expected: TestExpected{
				pos: util.Vec2[float32]{X: 3, Y: 5},
				rect: quadtree.Rect{
					X: 3, Y: 5, W: 15, H: 15,
				},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				dt:            16,
				dir:           util.Vec2[int8]{X: -1, Y: 0},
				velocity:      2,
				collisionRect: quadtree.Rect{X: 50, Y: 5, W: 15, H: 15},
			},
			Expected: TestExpected{
				pos: util.Vec2[float32]{X: 50 - float32(16)*2, Y: 5},
				rect: quadtree.Rect{
					X: 18, Y: 5, W: 15, H: 15,
				},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, TestExpected]) {
		parent := core.NewBaseGameObject(
			0,
			"parent",
			util.Vec2[float32]{
				X: float32(testCase.Input.collisionRect.X),
				Y: float32(testCase.Input.collisionRect.Y),
			},
			nil,
		)
		// detect no elements so no restrictions
		collisionSys := MockCollisionSystem{elementsToDetect: []quadtree.QuadElement{}}
		collider := collision.NewCollider(
			0,
			"rb_collider",
			testCase.Input.collisionRect,
			&collisionSys,
			&collisionSys,
			[]func(els []quadtree.QuadElement){},
			nil,
		)
		rb := NewRigidBody(0, "rigidbody", testCase.Input.velocity, nil, collider, &collisionSys, false)

		rb.SetParent(&parent)
		collider.SetParent(&parent)
		rb.Pos = parent.Pos

		rb.Dir = testCase.Input.dir

		collisionSys.Mock.On("DetectCollisions", testCase.Expected.rect)

		rb.OnUpdate(testCase.Input.dt, nil)

		require.Equal(t, testCase.Expected.pos, parent.Pos)
		collisionSys.AssertExpectations(t)
	})
}

func TestOnUpdateShouldRestrictMovement(t *testing.T) {
	type TestInput struct {
		dt               uint64
		dir              util.Vec2[int8]
		velocity         float32
		collisionRect    quadtree.Rect
		elementsToDetect []quadtree.QuadElement
	}

	type TestExpected struct {
		// the new position of the parent element
		pos util.Vec2[float32]
	}

	const NAME string = "should update position from %+v with restrictions"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.collisionRect)
	}

	cases := []util.TestCase[TestInput, TestExpected]{
		{
			Name: getName,
			Input: TestInput{
				dt:            10,
				dir:           util.Vec2[int8]{X: 1, Y: 1},
				velocity:      1,
				collisionRect: quadtree.Rect{X: 0, Y: 0, W: 5, H: 5},
				elementsToDetect: []quadtree.QuadElement{
					{Id: "id", Rect: quadtree.Rect{X: 5, Y: 0, W: 5, H: 5}},
				},
			},
			Expected: TestExpected{
				// account for diagonal speed reduction
				pos: util.Vec2[float32]{X: 0, Y: float32(10) * 0.7071},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, TestExpected]) {
		store := MockGameObjectStore{}

		// set parent and children with layer 1 since mock restrictions will have layer 0
		// and we only restrict when layers differ
		parent := core.NewBaseGameObject(
			1,
			"parent",
			util.Vec2[float32]{
				X: float32(testCase.Input.collisionRect.X),
				Y: float32(testCase.Input.collisionRect.Y),
			},
			&store,
		)
		// detect elements that may cause restrictions
		collisionSys := MockCollisionSystem{elementsToDetect: testCase.Input.elementsToDetect}
		collider := collision.NewCollider(
			1,
			"rb_collider",
			testCase.Input.collisionRect,
			&collisionSys,
			&collisionSys,
			[]func(els []quadtree.QuadElement){},
			&store,
		)
		rb := NewRigidBody(1, "rigidbody", testCase.Input.velocity, &store, collider, &collisionSys, false)

		rb.SetParent(&parent)
		collider.SetParent(&parent)
		rb.Pos = parent.Pos

		rb.Dir = testCase.Input.dir

		store.On("GetGameObject", mock.Anything)
		collisionSys.Mock.On("DetectCollisions", mock.Anything)

		rb.OnUpdate(testCase.Input.dt, nil)

		require.Equal(t, testCase.Expected.pos, parent.Pos)
		collisionSys.AssertExpectations(t)
	})
}
