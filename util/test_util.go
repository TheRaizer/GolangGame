package util

import "testing"

type TestCase[T any, S any] struct {
	Name     func(input T) string
	Input    T
	Expected S
}

type TestFunc[T any, S any] func(testCase TestCase[T, S])

// iterates a set of test cases and executes the given test exec on them
func IterateTestCases[T any, S any](cases []TestCase[T, S], t *testing.T, testFunc TestFunc[T, S]) {
	for _, testCase := range cases {
		t.Run(testCase.Name(testCase.Input), func(_ *testing.T) {
			testFunc(testCase)
		})
	}
}
