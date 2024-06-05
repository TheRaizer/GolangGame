package quadtree

import (
	"github.com/TheRaizer/GolangGame/util"
	"testing"
)

func TestIsLeaf(t *testing.T) {
	var cases = []util.TestCase[[4]*QuadNode, bool]{
		{Name: "Should be leaf", Input: [4]*QuadNode{nil, nil, nil, nil}, Expected: true},
		{Name: "Should not be leaf", Input: [4]*QuadNode{{}, {}, {}, {}}, Expected: false},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[4]*QuadNode, bool]) {
		node := QuadNode{children: testCase.Input}
		isLeaf := node.isLeaf()

		if isLeaf != testCase.Expected {
			t.Errorf("got %t, want %t", isLeaf, testCase.Expected)
		}

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
				if quadRect != nil {
					t.Errorf("got %v when it should have been nil", quadRect)
					return
				}
				// is nil so do not raise errors and move onto next parentRect
				continue
			}

			if *quadRect != *testCase.Expected[i] {
				t.Errorf("got %v, want %v", *quadRect, *testCase.Expected[i])
			}
		}

	})
}
