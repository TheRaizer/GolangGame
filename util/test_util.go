package util

import "testing"

type TestCase[Input any, Expected any] struct {
	Name     func(input Input) string
	Input    Input
	Expected Expected
}

type TestFunc[Input any, Expected any] func(testCase TestCase[Input, Expected])

// iterates a set of test cases and executes the given test exec on them
func IterateTestCases[Input any, Expected any](
	cases []TestCase[Input, Expected],
	t *testing.T,
	testFunc TestFunc[Input, Expected],
) {
	for _, testCase := range cases {
		t.Run(testCase.Name(testCase.Input), func(_ *testing.T) {
			testFunc(testCase)
		})
	}
}
