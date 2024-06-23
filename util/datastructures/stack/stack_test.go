package stack

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/require"
)

func TestPeek(t *testing.T) {
	const NAME string = "should return correct value for stack %+v"
	getName := func(input []int) string {
		return fmt.Sprintf(NAME, input)
	}
	cases := []util.TestCase[[]int, int]{
		{
			Name:     getName,
			Input:    []int{1, 3, 4, 6},
			Expected: 6,
		},
		{
			Name:     getName,
			Input:    []int{33, 2, 1},
			Expected: 1,
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[]int, int]) {
		stack := Stack[int]{items: make([]int, 0)}
		stack.items = testCase.Input

		require.Equal(t, testCase.Expected, stack.Peek())
	})
}

func TestPeekPanicsOnEmptyStack(t *testing.T) {
	// empty stack
	stack := Stack[int]{items: make([]int, 0)}

	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.Equal(t, "cannot peek an empty stack", r)
	}()

	stack.Peek()
}

func TestPush(t *testing.T) {
	const NAME string = "should push values %+v into stack correctly"
	getName := func(input []int) string {
		return fmt.Sprintf(NAME, input)
	}
	cases := []util.TestCase[[]int, Stack[int]]{
		{
			Name:  getName,
			Input: []int{1, 3, 4, 6},
			Expected: Stack[int]{
				items: []int{1, 3, 4, 6},
			},
		},
		{
			Name:  getName,
			Input: []int{1},
			Expected: Stack[int]{
				items: []int{1},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[]int, Stack[int]]) {
		stack := Stack[int]{items: make([]int, 0)}
		for _, val := range testCase.Input {
			stack.Push(val)
		}

		require.Equal(t, testCase.Expected, stack)
	})
}

func TestPop(t *testing.T) {
	type TestExpected struct {
		val   int
		items []int
	}

	const NAME string = "should return correct popped value for stack %+v"
	getName := func(input []int) string {
		return fmt.Sprintf(NAME, input)
	}
	cases := []util.TestCase[[]int, TestExpected]{
		{
			Name:  getName,
			Input: []int{1, 3, 4, 6},
			Expected: TestExpected{
				val:   6,
				items: []int{1, 3, 4},
			},
		},
		{
			Name:  getName,
			Input: []int{33, 2, 1},
			Expected: TestExpected{
				val:   1,
				items: []int{33, 2},
			},
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[]int, TestExpected]) {
		stack := Stack[int]{items: make([]int, 0)}
		stack.items = testCase.Input

		require.Equal(t, testCase.Expected.val, stack.Pop())
		require.Equal(t, testCase.Expected.items, stack.items)
	})
}

func TestPopPanicsOnEmptyStack(t *testing.T) {
	// empty stack
	stack := Stack[int]{items: make([]int, 0)}

	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.Equal(t, "cannot pop off an empty stack", r)
	}()

	stack.Pop()
}

func TestIsEmpty(t *testing.T) {
	const NAME string = "should return whether stack %+v is empty"
	getName := func(input []int) string {
		return fmt.Sprintf(NAME, input)
	}
	cases := []util.TestCase[[]int, bool]{
		{
			Name:     getName,
			Input:    []int{1, 3, 4, 6},
			Expected: false,
		},
		{
			Name:     getName,
			Input:    []int{},
			Expected: true,
		},
	}

	util.IterateTestCases(cases, t, func(testCase util.TestCase[[]int, bool]) {
		stack := Stack[int]{items: make([]int, 0)}
		stack.items = testCase.Input

		require.Equal(t, testCase.Expected, stack.IsEmpty())
	})
}
