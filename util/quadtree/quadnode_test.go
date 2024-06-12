package quadtree

import (
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/require"
)

func TestIsLeaf(t *testing.T) {
	var cases = []util.TestCase[[4]*QuadNode, bool]{
		{Name: "Should be leaf", Input: [4]*QuadNode{nil, nil, nil, nil}, Expected: true},
		{Name: "Should not be leaf", Input: [4]*QuadNode{{}, {}, {}, {}}, Expected: false},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[4]*QuadNode, bool]) {
		node := QuadNode{children: testCase.Input}
		isLeaf := node.isLeaf()

		require.Equal(t, testCase.Expected, isLeaf)
	})
}

func TestComputeQuadRect(t *testing.T) {
	var parentRects = []Rect{
		{0, 0, 10, 10},
		{1, 30, 88, 13},
		{23, 11, 30, 22},
	}

	var cases = []util.TestCase[int, []*Rect]{
		{
			Name:     "Should return nil if quadrantIdx is not valid",
			Input:    4,
			Expected: nil,
		},
		{
			Name:     "Should return correct 0 index (NW) rect",
			Input:    0,
			Expected: []*Rect{{0, 0, 5, 5}, {1, 30, 44, 6}, {23, 11, 15, 11}},
		},
		{
			Name:     "Should return correct 1 index (NE) rect",
			Input:    1,
			Expected: []*Rect{{5, 0, 5, 5}, {45, 30, 44, 6}, {38, 11, 15, 11}},
		},
		{
			Name:     "Should return correct 2 index (SW) rect",
			Input:    2,
			Expected: []*Rect{{0, 5, 5, 5}, {1, 36, 44, 6}, {23, 22, 15, 11}},
		},
		{
			Name:     "Should return correct 3 index (SE) rect",
			Input:    3,
			Expected: []*Rect{{5, 5, 5, 5}, {45, 36, 44, 6}, {38, 22, 15, 11}},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[int, []*Rect]) {

		for i, parentRect := range parentRects {
			quadRect := ComputeQuadRect(parentRect, testCase.Input)

			if testCase.Input > 3 || testCase.Input < 0 {
				require.Nil(t, quadRect)
				continue
			}

			require.Equal(t, *testCase.Expected[i], *quadRect)

		}

	})
}

func TestQuadrantContaining(t *testing.T) {
	type TestInput struct {
		nodeRect Rect
		el       QuadElement
	}

	var cases = []util.TestCase[TestInput, int32]{
		{
			Name: "Should panic if element is not contained in the nodeRect",
			Input: TestInput{
				nodeRect: Rect{0, 0, 10, 10},
				el:       QuadElement{Rect: Rect{20, 20, 5, 5}, Id: "id"},
			},
			Expected: -1,
		},
		{
			Name: "Should return 0 if the element rect is contained in NW quadrant",
			Input: TestInput{
				nodeRect: Rect{0, 0, 10, 10},
				el:       QuadElement{Rect: Rect{0, 0, 3, 3}, Id: "id"},
			},
			Expected: 0,
		},
		{
			Name: "Should return 1 if the element rect is contained in NE quadrant",
			Input: TestInput{
				nodeRect: Rect{5, 5, 15, 15},
				el:       QuadElement{Rect: Rect{13, 5, 5, 5}, Id: "id"},
			},
			Expected: 1,
		},
		{
			Name: "Should return 2 if the element rect is contained in SW quadrant",
			Input: TestInput{
				nodeRect: Rect{6, 8, 10, 10},
				el:       QuadElement{Rect: Rect{7, 15, 3, 3}, Id: "id"},
			},
			Expected: 2,
		},
		{
			Name: "Should return 3 if the element rect is contained in SE quadrant",
			Input: TestInput{
				nodeRect: Rect{9, 2, 8, 8},
				el:       QuadElement{Rect: Rect{13, 6, 2, 2}, Id: "id"},
			},
			Expected: 3,
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, int32]) {
		if testCase.Expected == -1 {
			defer func() {
				r := recover()
				require.NotNil(t, r)
				require.Equal(t, "element is not contained in the given nodeRect", r)
			}()
		}
		quadrantIdx, _ := QuadrantContaining(testCase.Input.nodeRect, testCase.Input.el)
		require.Equal(t, testCase.Expected, quadrantIdx)
	})
}
