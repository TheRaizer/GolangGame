package collision

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCollisionSystem struct {
	CollisionSystem
	mock.Mock
}

func (mockSys *MockCollisionSystem) RegisterObject(collider *Collider) {
	// dereference for asssertions
	mockSys.Called(*collider)
}

func (mockSys *MockCollisionSystem) UpdateCollider(id string, oldRect quadtree.Rect, newRect quadtree.Rect) {
	mockSys.Called(id, oldRect, newRect)
}

func TestNewColliderShouldAddObjectIntoCollisionSys(t *testing.T) {
	mockSys := &MockCollisionSystem{}

	mockSys.On("RegisterObject", mock.AnythingOfType("Collider"))

	collider := NewCollider(
		0,
		"name",
		quadtree.Rect{},
		mockSys,
		mockSys,
		make([]func(els []quadtree.QuadElement), 0),
		nil,
	)

	mockSys.AssertExpectations(t)

	registeredCollider := mockSys.Mock.Calls[0].Arguments[0]

	require.Equal(t, *collider, registeredCollider)
}

func TestUpdatePosShouldUpdateColliderInSys(t *testing.T) {
	type TestInput struct {
		rect     quadtree.Rect
		id       string
		distance util.Vec2[float32]
	}

	type Expected struct {
		id      string
		oldRect quadtree.Rect
		newRect quadtree.Rect
	}

	const NAME = "Should update collider with distance: %+v"

	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.distance)
	}

	cases := []util.TestCase[TestInput, Expected]{
		{
			Name: getName,
			Input: TestInput{
				quadtree.Rect{X: 0, Y: 0, W: 45, H: 22},
				"id1",
				util.Vec2[float32]{X: 3, Y: 12},
			},
			Expected: Expected{
				"id1",
				quadtree.Rect{X: 0, Y: 0, W: 45, H: 22},
				quadtree.Rect{X: 3, Y: 12, W: 45, H: 22},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				quadtree.Rect{X: 4, Y: 1, W: 12, H: 10},
				"id2",
				util.Vec2[float32]{X: -3, Y: 7},
			},
			Expected: Expected{
				"id2",
				quadtree.Rect{X: 4, Y: 1, W: 12, H: 10},
				quadtree.Rect{X: 1, Y: 8, W: 12, H: 10},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, Expected]) {
		mockSys := &MockCollisionSystem{}

		mockSys.On("RegisterObject", mock.Anything)
		mockSys.On(
			"UpdateCollider",
			testCase.Expected.id,
			testCase.Expected.oldRect,
			testCase.Expected.newRect,
		)

		collider := NewCollider(
			0,
			testCase.Input.id,
			testCase.Input.rect,
			mockSys,
			mockSys,
			make([]func(els []quadtree.QuadElement), 0),
			nil,
		)

		collider.UpdatePos(testCase.Input.distance.X, testCase.Input.distance.Y)
		mockSys.AssertExpectations(t)
	})

}

func TestUpdatePosShouldUpdateColliderRect(t *testing.T) {
	type TestInput struct {
		rect     quadtree.Rect
		distance util.Vec2[float32]
	}

	type Expected struct {
		newRect quadtree.Rect
	}

	const NAME = "Should update collider with distance: %+v"

	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.distance)
	}

	cases := []util.TestCase[TestInput, Expected]{
		{
			Name: getName,
			Input: TestInput{
				quadtree.Rect{X: 0, Y: 0, W: 45, H: 22},
				util.Vec2[float32]{X: 3, Y: 12},
			},
			Expected: Expected{
				quadtree.Rect{X: 3, Y: 12, W: 45, H: 22},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				quadtree.Rect{X: 4, Y: 1, W: 12, H: 10},
				util.Vec2[float32]{X: -3, Y: 7},
			},
			Expected: Expected{
				quadtree.Rect{X: 1, Y: 8, W: 12, H: 10},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, Expected]) {
		mockSys := &MockCollisionSystem{}

		mockSys.On("RegisterObject", mock.Anything)
		mockSys.On(
			"UpdateCollider",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		)

		collider := NewCollider(
			0,
			"id",
			testCase.Input.rect,
			mockSys,
			mockSys,
			make([]func(els []quadtree.QuadElement), 0),
			nil,
		)

		collider.UpdatePos(testCase.Input.distance.X, testCase.Input.distance.Y)
		require.Equal(t, collider.Rect, testCase.Expected.newRect)
	})

}

func TestAddCollisionEventShouldRegisterEvent(t *testing.T) {
	mockSys := &MockCollisionSystem{}
	mockSys.On("RegisterObject", mock.Anything)

	collider := NewCollider(
		0,
		"id",
		quadtree.Rect{},
		mockSys,
		mockSys,
		make([]func(els []quadtree.QuadElement), 0),
		nil,
	)

	event := func(els []quadtree.QuadElement) {}
	expectedEventPtr := reflect.ValueOf(event).Pointer()

	collider.AddCollisionEvent(event)

	givenEventPtr := reflect.ValueOf(collider.collisionEvents[0]).Pointer()

	require.Equal(t, expectedEventPtr, givenEventPtr)
}

func TestOnCollisionShouldCallCollisionEvents(t *testing.T) {
	mockSys := &MockCollisionSystem{}
	mockSys.On("RegisterObject", mock.Anything)

	collider := NewCollider(
		0,
		"id",
		quadtree.Rect{},
		mockSys,
		mockSys,
		make([]func(els []quadtree.QuadElement), 0),
		nil,
	)

	expectedEls := []quadtree.QuadElement{
		{
			Rect: quadtree.Rect{X: 3, Y: 2, W: 4, H: 4},
			Id:   "id1",
		},
	}
	calls := [3]bool{}

	event1 := func(els []quadtree.QuadElement) {
		calls[0] = true
		require.ElementsMatch(t, expectedEls, els)
	}
	event2 := func(els []quadtree.QuadElement) {
		calls[1] = true
		require.ElementsMatch(t, expectedEls, els)

	}
	event3 := func(els []quadtree.QuadElement) {
		calls[2] = true
		require.ElementsMatch(t, expectedEls, els)

	}

	collider.collisionEvents = []func(els []quadtree.QuadElement){
		event1,
		event2,
		event3,
	}

	collider.OnCollision(expectedEls)

	for _, called := range calls {
		require.True(t, called)
	}
}
