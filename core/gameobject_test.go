package core

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestGameObjectStore struct {
	mock.Mock
}

func (testStore *TestGameObjectStore) AddGameObject(gameObject GameObject) {
	testStore.Called(gameObject)
}
func (testStore *TestGameObjectStore) RemoveGameObject(id string) {
	testStore.Called(id)
}
func (testStore *TestGameObjectStore) GetGameObject(id string) GameObject {
	testStore.Called(id)
	return &BaseGameObject{}
}

type MockGameObject struct {
	BaseGameObject
	mock.Mock
}

// mock the update function
func (obj *MockGameObject) UpdatePos(distX float32, distY float32) {
	// call this so the call is registered on the mock
	obj.Called(distX, distY)
}

func TestUpdatePosShouldUpdateCurrentObjectPos(t *testing.T) {
	type TestInput struct {
		initialPos util.Vec2[float32]
		distance   util.Vec2[float32]
	}

	const NAME string = "should update from position and distance %+v correctly"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input)
	}

	cases := []util.TestCase[TestInput, util.Vec2[float32]]{
		{
			Name: getName,
			Input: TestInput{
				initialPos: util.Vec2[float32]{X: 0, Y: 0},
				distance:   util.Vec2[float32]{X: 5, Y: 11},
			},
			Expected: util.Vec2[float32]{
				X: 5,
				Y: 11,
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, util.Vec2[float32]]) {

		gameObject := NewBaseGameObject(0, "test_name", testCase.Input.initialPos, &TestGameObjectStore{})

		gameObject.UpdatePos(testCase.Input.distance.X, testCase.Input.distance.Y)

		require.Equal(t, testCase.Expected, gameObject.Pos)
	})
}

func TestUpdatePosShouldUpdateChildrenPositions(t *testing.T) {
	type TestInput struct {
		parentObj BaseGameObject
		distance  util.Vec2[float32]
	}

	const NAME string = "should update children from position and distance %+v correctly"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input)
	}

	cases := []util.TestCase[TestInput, map[string]util.Vec2[float32]]{
		{
			Name: getName,
			Input: TestInput{
				parentObj: BaseGameObject{
					Pos:  util.Vec2[float32]{X: 0, Y: 0},
					name: "name",
					children: map[string]GameObject{
						"id1": &BaseGameObject{
							name: "id1",
							Pos:  util.Vec2[float32]{X: 0, Y: 0},
						},
					},
				},
				distance: util.Vec2[float32]{X: 5, Y: 11},
			},
			Expected: map[string]util.Vec2[float32]{
				"id1": {
					X: 5,
					Y: 11,
				},
			},
		},
		{
			Name: getName,
			Input: TestInput{
				parentObj: BaseGameObject{
					Pos:  util.Vec2[float32]{X: 0, Y: 0},
					name: "name",
					children: map[string]GameObject{
						"id1": &BaseGameObject{
							name: "id1",
							Pos:  util.Vec2[float32]{X: 0, Y: 11},
						},

						"id2": &BaseGameObject{
							name: "id2",
							Pos:  util.Vec2[float32]{X: 1, Y: 3},
						},
					},
				},
				distance: util.Vec2[float32]{X: 3, Y: -1},
			},
			Expected: map[string]util.Vec2[float32]{
				"id1": {
					X: 3,
					Y: 10,
				},
				"id2": {
					X: 4,
					Y: 2,
				},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, map[string]util.Vec2[float32]]) {
		testCase.Input.parentObj.UpdatePos(testCase.Input.distance.X, testCase.Input.distance.Y)

		for id, child := range testCase.Input.parentObj.children {
			require.NotNil(t, testCase.Expected[id])

			child, ok := child.(*BaseGameObject)
			require.True(t, ok)
			require.Equal(t, testCase.Expected[id], child.Pos)
		}
	})
}

func TestAddChild(t *testing.T) {
	const NAME string = "should add gameobjects %+v to the parent correctly"
	getName := func(input []*MockGameObject) string {
		return fmt.Sprintf(NAME, input)
	}

	cases := []util.TestCase[[]*MockGameObject, map[string]*MockGameObject]{
		{
			Name: getName,
			Input: []*MockGameObject{
				{
					BaseGameObject: BaseGameObject{
						name: "id1",
						Pos:  util.Vec2[float32]{X: 3, Y: 2},
					},
				},
				{
					BaseGameObject: BaseGameObject{
						name: "id2",
						Pos:  util.Vec2[float32]{X: 1, Y: 0},
					},
				},
			},
			Expected: map[string]*MockGameObject{
				"id1": {
					BaseGameObject: BaseGameObject{
						name: "id1",
						Pos:  util.Vec2[float32]{X: 3, Y: 2},
					},
				},
				"id2": {
					BaseGameObject: BaseGameObject{
						name: "id2",
						Pos:  util.Vec2[float32]{X: 1, Y: 0},
					},
				},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[]*MockGameObject, map[string]*MockGameObject]) {
		mockStore := TestGameObjectStore{}
		parentObj := NewBaseGameObject(0, "test_name", util.Vec2[float32]{X: 0, Y: 0}, &mockStore)

		for _, child := range testCase.Input {
			// setup call expectations
			mockStore.On("AddGameObject", child)
			child.On("UpdatePos", parentObj.Pos.X, parentObj.Pos.Y)

			parentObj.AddChild(child)

			// evaluate call expectations
			mockStore.AssertExpectations(t)
			child.AssertExpectations(t)
		}

		for id, expectedObj := range testCase.Expected {
			require.NotNil(t, parentObj.children[id])

			// cast so we can compare the BaseGameObject
			obj, ok := parentObj.children[id].(*MockGameObject)
			expectedObj.parent = &parentObj // set parent

			require.Equal(t, ok, true)
			require.Equal(t, expectedObj.BaseGameObject, obj.BaseGameObject)
		}
	})

}

func TestRemoveChild(t *testing.T) {
	type TestInput struct {
		parentObj  *BaseGameObject
		idToRemove string
	}

	const NAME string = "should remove child with id '%+v' from the parent correctly"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.idToRemove)
	}

	cases := []util.TestCase[TestInput, struct{}]{
		{
			Name: getName,
			Input: TestInput{
				parentObj: &BaseGameObject{
					name: "id",
					children: map[string]GameObject{
						"id1": &BaseGameObject{name: "id1"},
					},
				},
				idToRemove: "id1",
			},
		},
		{
			Name: getName,
			Input: TestInput{
				parentObj: &BaseGameObject{
					name: "id",
					children: map[string]GameObject{
						"id1": &BaseGameObject{name: "id1"},
						"id2": &BaseGameObject{name: "id2"},
					},
				},
				idToRemove: "id2",
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, struct{}]) {
		// initialize the mock store
		mockStore := &TestGameObjectStore{}
		mockStore.On("RemoveGameObject", testCase.Input.idToRemove)
		testCase.Input.parentObj.GameObjectStore = mockStore

		require.NotNil(t, testCase.Input.parentObj.children[testCase.Input.idToRemove])

		testCase.Input.parentObj.RemoveChild(testCase.Input.idToRemove)

		require.Nil(t, testCase.Input.parentObj.children[testCase.Input.idToRemove])
		mockStore.AssertExpectations(t)
	})
}

func TestIDRetrieval(t *testing.T) {
	const NAME string = "should return name '%+v' from the object ID correctly"
	getName := func(input BaseGameObject) string {
		return fmt.Sprintf(NAME, input.name)
	}

	cases := []util.TestCase[BaseGameObject, string]{
		{
			Name:     getName,
			Input:    NewBaseGameObject(0, "name1", util.Vec2[float32]{X: 0, Y: 0}, nil),
			Expected: "name1",
		},
		{
			Name:     getName,
			Input:    NewBaseGameObject(0, "another_name", util.Vec2[float32]{X: 0, Y: 0}, nil),
			Expected: "another_name",
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[BaseGameObject, string]) {
		id := testCase.Input.ID()

		require.Equal(t, testCase.Expected, id)
	})

}

func TestGetPos(t *testing.T) {
	const NAME string = "should return pos %+v from the object correctly"
	getName := func(input BaseGameObject) string {
		return fmt.Sprintf(NAME, input.Pos)
	}

	cases := []util.TestCase[BaseGameObject, util.Vec2[float32]]{
		{
			Name:     getName,
			Input:    NewBaseGameObject(0, "name1", util.Vec2[float32]{X: 3, Y: 31}, nil),
			Expected: util.Vec2[float32]{X: 3, Y: 31},
		},
		{
			Name:     getName,
			Input:    NewBaseGameObject(0, "another_name", util.Vec2[float32]{X: 32, Y: 23}, nil),
			Expected: util.Vec2[float32]{X: 32, Y: 23},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[BaseGameObject, util.Vec2[float32]]) {
		require.Equal(t, testCase.Expected, testCase.Input.Pos)
	})

}
