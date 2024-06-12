package quadtree

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/require"
)

func TestIsLeaf(t *testing.T) {
	var cases = []util.TestCase[[4]*QuadNode, bool]{
		{
			Name: func(input [4]*QuadNode) string {
				return "Should be leaf"
			},
			Input: [4]*QuadNode{nil, nil, nil, nil}, Expected: true,
		},
		{
			Name: func(input [4]*QuadNode) string {
				return "Should not be leaf"
			},
			Input: [4]*QuadNode{{}, {}, {}, {}}, Expected: false,
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[4]*QuadNode, bool]) {
		node := QuadNode{children: testCase.Input}
		isLeaf := node.isLeaf()

		require.Equal(t, testCase.Expected, isLeaf)
	})
}

func TestSplit(t *testing.T) {
	type ExpectedOutput struct {
		expectedChildren [4]QuadNode
		expectedNodeEls  []QuadElement
	}

	type TestInput struct {
		node     QuadNode
		nodeRect Rect
	}

	testName := func(input TestInput) string {
		return fmt.Sprintf("Should allocate elements: %v correctly across the node rect: %+v", input.node.els, input.nodeRect)
	}

	// ensure that the nodeRect contains all the QuadNode quad elements
	var cases = []util.TestCase[TestInput, *ExpectedOutput]{
		{
			Name: func(input TestInput) string {
				return "Should panic if the given node is not a leaf node"
			},
			Input: TestInput{
				QuadNode{
					// add children so not a leaf
					[4]*QuadNode{{}, {}, {}, {}},
					[]QuadElement{},
				},
				Rect{},
			},
			Expected: nil,
		},
		{
			Name: testName,
			Input: TestInput{
				QuadNode{
					els: []QuadElement{
						{Rect{0, 0, 8, 8}, "id"},
					},
				},
				Rect{0, 0, 10, 10},
			},
			Expected: &ExpectedOutput{
				[4]QuadNode{{}, {}, {}, {}},
				[]QuadElement{{Rect{0, 0, 8, 8}, "id"}},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, *ExpectedOutput]) {
		if testCase.Expected == nil {
			defer func() {
				r := recover()
				require.NotNil(t, r)
				require.Equal(t, "only a leaf node can be split", r)
			}()
		}
		testCase.Input.node.split(testCase.Input.nodeRect)
		derefdChildren := [4]QuadNode{}

		for i := range derefdChildren {
			derefdChildren[i] = *testCase.Input.node.children[i]
		}

		require.ElementsMatch(t, testCase.Expected.expectedChildren, derefdChildren)
		require.ElementsMatch(t, testCase.Expected.expectedNodeEls, testCase.Input.node.els)
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
			Name: func(input int) string {
				return "Should return nil if quadrantIdx is not valid"
			},
			Input:    4,
			Expected: nil,
		},
		{
			Name: func(input int) string {
				return fmt.Sprintf("Should return correct %d index (NW) rect", input)
			},
			Input:    0,
			Expected: []*Rect{{0, 0, 5, 5}, {1, 30, 44, 6}, {23, 11, 15, 11}},
		},
		{
			Name: func(input int) string {
				return fmt.Sprintf("Should return correct %d index (NE) rect", input)
			},
			Input:    1,
			Expected: []*Rect{{5, 0, 5, 5}, {45, 30, 44, 6}, {38, 11, 15, 11}},
		},
		{
			Name: func(input int) string {
				return fmt.Sprintf("Should return correct %d index (SW) rect", input)
			},
			Input:    2,
			Expected: []*Rect{{0, 5, 5, 5}, {1, 36, 44, 6}, {23, 22, 15, 11}},
		},
		{
			Name: func(input int) string {
				return fmt.Sprintf("Should return correct %d index (SE) rect", input)
			},
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
			Name: func(input TestInput) string {
				return "Should panic if element is not contained in the nodeRect"
			},
			Input: TestInput{
				nodeRect: Rect{0, 0, 10, 10},
				el:       QuadElement{Rect: Rect{20, 20, 5, 5}, Id: "id"},
			},
			Expected: -1,
		},
		{
			Name: func(input TestInput) string {
				return "Should return 0 if the element rect is contained in NW quadrant"
			},
			Input: TestInput{
				nodeRect: Rect{0, 0, 10, 10},
				el:       QuadElement{Rect: Rect{0, 0, 3, 3}, Id: "id"},
			},
			Expected: 0,
		},
		{
			Name: func(input TestInput) string {
				return "Should return 1 if the element rect is contained in NE quadrant"
			},
			Input: TestInput{
				nodeRect: Rect{5, 5, 15, 15},
				el:       QuadElement{Rect: Rect{13, 5, 5, 5}, Id: "id"},
			},
			Expected: 1,
		},
		{
			Name: func(input TestInput) string {
				return "Should return 2 if the element rect is contained in SW quadrant"
			},
			Input: TestInput{
				nodeRect: Rect{6, 8, 10, 10},
				el:       QuadElement{Rect: Rect{7, 15, 3, 3}, Id: "id"},
			},
			Expected: 2,
		},
		{
			Name: func(input TestInput) string {
				return "Should return 3 if the element rect is contained in SE quadrant"
			},
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
