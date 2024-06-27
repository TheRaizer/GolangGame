package collision

import (
	"fmt"
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
		"name",
		quadtree.Rect{},
		mockSys,
		mockSys,
		make([]func(els []quadtree.QuadElement), 0),
		nil,
		false,
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
			testCase.Input.id,
			testCase.Input.rect,
			mockSys,
			mockSys,
			make([]func(els []quadtree.QuadElement), 0),
			nil,
			false,
		)

		collider.UpdatePos(testCase.Input.distance.X, testCase.Input.distance.Y)
		mockSys.AssertExpectations(t)
	})

}
