package quadtree

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/require"
)

func TestInsertion(t *testing.T) {
	type TestInput struct {
		tree QuadTree
		els  []QuadElement
	}

	const NAME string = "should insert elements %+v into a quad tree successfully"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.els)
	}
	var cases = []util.TestCase[TestInput, QuadNode]{
		{
			Name: getName,
			Input: TestInput{
				QuadTree{
					threshold:  2,
					maxDepth:   4,
					globalRect: Rect{0, 0, 100, 100},
					root:       &QuadNode{},
				},
				[]QuadElement{
					{Rect{0, 0, 5, 5}, "id1"},
					{Rect{0, 0, 7, 7}, "id2"},
				},
			},
			Expected: QuadNode{els: []QuadElement{{Rect{0, 0, 5, 5}, "id1"}, {Rect{0, 0, 7, 7}, "id2"}}},
		},
		{
			Name: getName,
			Input: TestInput{
				QuadTree{
					threshold:  2,
					maxDepth:   4,
					globalRect: Rect{0, 0, 100, 100},
					root:       &QuadNode{},
				},
				[]QuadElement{
					{Rect{0, 0, 5, 5}, "id1"},
					{Rect{0, 0, 7, 7}, "id2"},
					{Rect{40, 40, 8, 8}, "id3"},
				}},
			Expected: QuadNode{children: [4]*QuadNode{
				{
					children: [4]*QuadNode{{
						els: []QuadElement{
							{Rect{0, 0, 5, 5}, "id1"},
							{Rect{0, 0, 7, 7}, "id2"},
						}},
						{},
						{},
						{els: []QuadElement{
							{Rect{40, 40, 8, 8}, "id3"},
						}},
					},
				}, {}, {}, {},
			}},
		},
		{
			Name: getName,
			Input: TestInput{
				QuadTree{
					threshold:  2,
					maxDepth:   2,
					globalRect: Rect{0, 0, 100, 100},
					root:       &QuadNode{},
				},
				[]QuadElement{
					{Rect{0, 0, 5, 5}, "id1"},
					{Rect{0, 0, 7, 7}, "id2"},
					{Rect{40, 40, 30, 30}, "id3"},
					{Rect{70, 70, 5, 5}, "id4"},
					{Rect{80, 5, 12, 12}, "id5"},
					{Rect{10, 51, 12, 12}, "id6"},
					{Rect{70, 65, 3, 2}, "id7"},
					{Rect{90, 93, 4, 4}, "id8"},
					{Rect{55, 56, 30, 30}, "id9"},
					{Rect{55, 56, 5, 4}, "id10"},
				}},
			Expected: QuadNode{
				children: [4]*QuadNode{
					{ // NW
						els: []QuadElement{
							{Rect{0, 0, 5, 5}, "id1"},
							{Rect{0, 0, 7, 7}, "id2"},
						},
					},
					{ // NE
						els: []QuadElement{
							{Rect{80, 5, 12, 12}, "id5"},
						},
					},
					{ // SW
						els: []QuadElement{
							{Rect{10, 51, 12, 12}, "id6"},
						},
					},
					{ // SE
						children: [4]*QuadNode{
							{ // when at max depth allow past threshold
								els: []QuadElement{
									{Rect{70, 70, 5, 5}, "id4"},
									{Rect{70, 65, 3, 2}, "id7"},
									{Rect{55, 56, 5, 4}, "id10"},
								},
							},
							{},
							{},
							{
								els: []QuadElement{
									{Rect{90, 93, 4, 4}, "id8"},
								},
							},
						},
						els: []QuadElement{
							{Rect{55, 56, 30, 30}, "id9"},
						},
					},
				},
				els: []QuadElement{{Rect{40, 40, 30, 30}, "id3"}},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, QuadNode]) {
		for _, el := range testCase.Input.els {
			testCase.Input.tree.Insert(el)
		}

		require.Equal(t, testCase.Expected, *testCase.Input.tree.root)
	})
}
