package datastructures

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/require"
)

func TestCenter(t *testing.T) {
	const NAME string = "should return correct center for Rect%+v"
	getName := func(input Rect) string {
		return fmt.Sprintf(NAME, input)
	}

	var cases = []util.TestCase[Rect, util.Vec2[int32]]{
		{
			Name:  getName,
			Input: Rect{0, 0, 10, 10}, Expected: util.Vec2[int32]{X: 5, Y: 5},
		},
		{
			Name:  getName,
			Input: Rect{5, 3, 17, 14}, Expected: util.Vec2[int32]{X: 13, Y: 10},
		},
	}

	util.IterateTestCases(cases, t,
		func(testCase util.TestCase[Rect, util.Vec2[int32]]) {
			center := testCase.Input.Center()
			require.Equal(t, testCase.Expected, center)
		})
}

func TestContains(t *testing.T) {
	type TestInput struct {
		rect      Rect
		otherRect Rect
	}

	const NAME string = "should return whether Rect%+v contains Rect%+v"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.rect, input.otherRect)
	}
	var cases = []util.TestCase[TestInput, bool]{
		{
			Name: getName,
			Input: TestInput{
				Rect{0, 0, 10, 10},
				Rect{1, 1, 4, 4},
			},
			Expected: true,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{6, 5, 10, 7},
				Rect{2, 1, 4, 4},
			},
			Expected: false,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{0, 0, 5, 5},
				Rect{2, 2, 10, 10},
			},
			Expected: false,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{3, 2, 9, 7},
				Rect{5, 3, 4, 3},
			},
			Expected: true,
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, bool]) {
		contains := testCase.Input.rect.Contains(testCase.Input.otherRect)
		require.Equal(t, testCase.Expected, contains)
	})
}

func TestIntersects(t *testing.T) {
	type TestInput struct {
		rect      Rect
		otherRect Rect
	}

	const NAME string = "should return whether Rect%+v intersects Rect%+v"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.rect, input.otherRect)
	}
	var cases = []util.TestCase[TestInput, bool]{
		{
			Name: getName,
			Input: TestInput{
				Rect{0, 0, 10, 10},
				Rect{10, 10, 4, 4},
			},
			Expected: false,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{6, 5, 10, 7},
				Rect{2, 3, 5, 4},
			},
			Expected: true,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{0, 0, 5, 5},
				Rect{2, 2, 10, 10},
			},
			Expected: true,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{3, 2, 9, 7},
				Rect{5, 3, 4, 3},
			},
			Expected: true,
		},
		{
			Name: getName,
			Input: TestInput{
				Rect{5, 4, 15, 20},
				Rect{75, 40, 4, 3},
			},
			Expected: false,
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[TestInput, bool]) {
		intersects := testCase.Input.rect.Intersects(testCase.Input.otherRect)
		require.Equal(t, testCase.Expected, intersects)
	})
}

func TestRight(t *testing.T) {
	const NAME string = "should return correct right for Rect%+v"
	getName := func(input Rect) string {
		return fmt.Sprintf(NAME, input)
	}

	var cases = []util.TestCase[Rect, int32]{
		{
			Name:  getName,
			Input: Rect{0, 0, 10, 10}, Expected: 10,
		},
		{
			Name:  getName,
			Input: Rect{5, 3, 17, 14}, Expected: 22,
		},
	}

	util.IterateTestCases(cases, t,
		func(testCase util.TestCase[Rect, int32]) {
			right := testCase.Input.Right()
			require.Equal(t, testCase.Expected, right)
		})

}

func TestBottom(t *testing.T) {
	const NAME string = "should return correct bottom for Rect%+v"
	getName := func(input Rect) string {
		return fmt.Sprintf(NAME, input)
	}

	var cases = []util.TestCase[Rect, int32]{
		{
			Name:  getName,
			Input: Rect{0, 0, 10, 10}, Expected: 10,
		},
		{
			Name:  getName,
			Input: Rect{5, 3, 17, 14}, Expected: 17,
		},
	}

	util.IterateTestCases(cases, t,
		func(testCase util.TestCase[Rect, int32]) {
			bottom := testCase.Input.Bottom()
			require.Equal(t, testCase.Expected, bottom)
		})
}
