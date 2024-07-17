package image

import (
	"fmt"
	"testing"

	"github.com/TheRaizer/GolangGame/util"
	"github.com/stretchr/testify/require"
)

func TestPaethPred(t *testing.T) {
	const NAME string = "should return the correct byte when given a, b, c = %v"
	getName := func(input [3]byte) string {
		return fmt.Sprintf(NAME, input)
	}

	var cases = []util.TestCase[[3]byte, byte]{
		{
			Name:     getName,
			Input:    [3]byte{1, 2, 32},
			Expected: 1,
		},
		{
			Name:     getName,
			Input:    [3]byte{0, 0, 0},
			Expected: 0,
		},
		{
			Name:     getName,
			Input:    [3]byte{255, 150, 9},
			Expected: 255,
		},
		{
			Name:     getName,
			Input:    [3]byte{0, 1, 22},
			Expected: 0,
		},
	}

	util.IterateTestCases(cases, t,
		func(testCase util.TestCase[[3]byte, byte]) {
			require.Equal(
				t,
				testCase.Expected,
				paethPred(testCase.Input[0], testCase.Input[1], testCase.Input[2]),
			)
		})
}

func TestInversePaeth(t *testing.T) {
	type TestInput struct {
		rawPrevScanline  []byte
		filteredScanline []byte
		bpp              int
	}

	const NAME string = "should return the correct paeth defiltered scanline when given filtered scan line: %v, rawPrevScanline: %v, and bytes per pixel: %d"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.filteredScanline, input.rawPrevScanline, input.bpp)
	}

	var cases = []util.TestCase[TestInput, []byte]{
		{
			Name: getName,
			Input: TestInput{
				rawPrevScanline:  nil,
				filteredScanline: []byte{255, 0, 203, 255, 108, 10, 12, 255},
				bpp:              4,
			},
			Expected: []byte{255, 0, 203, 255, 107, 10, 215, 254},
		},
		{
			Name: getName,
			Input: TestInput{
				rawPrevScanline:  []byte{200, 101, 22, 1, 5, 4, 33, 209},
				filteredScanline: []byte{22, 10, 22, 84, 95, 0, 22, 255},
				bpp:              4,
			},
			Expected: []byte{222, 111, 44, 85, 100, 4, 66, 208},
		},
		{
			Name: getName,
			Input: TestInput{
				rawPrevScanline:  []byte{1, 2, 23, 4, 5, 1},
				filteredScanline: []byte{1, 100, 22, 0, 7, 2},
				bpp:              3,
			},
			Expected: []byte{2, 102, 45, 4, 109, 25},
		},
	}

	util.IterateTestCases(cases, t,
		func(testCase util.TestCase[TestInput, []byte]) {
			require.Equal(
				t,
				testCase.Expected,
				inversePaeth(testCase.Input.filteredScanline, testCase.Input.rawPrevScanline, testCase.Input.bpp),
			)
		})
}

func TestInverseSub(t *testing.T) {
	type TestInput struct {
		filteredScanline []byte
		bpp              int
	}

	const NAME string = "should return the correct sub defiltered scanline when given filtered scan line: %v, and bytes per pixel: %d"
	getName := func(input TestInput) string {
		return fmt.Sprintf(NAME, input.filteredScanline, input.bpp)
	}

	var cases = []util.TestCase[TestInput, []byte]{
		{
			Name: getName,
			Input: TestInput{
				filteredScanline: []byte{255, 0, 203, 255, 108, 10, 12, 255},
				bpp:              4,
			},
			Expected: []byte{255, 0, 203, 255, 107, 10, 215, 254},
		},
		{
			Name: getName,
			Input: TestInput{
				filteredScanline: []byte{255, 0, 203, 255, 108, 10, 12, 255},
				bpp:              2,
			},
			Expected: []byte{255, 0, 202, 255, 54, 9, 66, 8},
		},
	}

	util.IterateTestCases(cases, t,
		func(testCase util.TestCase[TestInput, []byte]) {
			require.Equal(
				t,
				testCase.Expected,
				inverseSub(testCase.Input.filteredScanline, testCase.Input.bpp),
			)
		})
}
