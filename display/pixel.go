package display

import (
	"math"
)

type Pixel struct {
	X, Y int
}

func (v1 Pixel) Add(v2 Pixel) Pixel {
	var sum Pixel

	sum.X = v1.X + v2.X
	sum.Y = v1.Y + v2.Y

	return sum
}

func (v Pixel) Negative() Pixel {
	return Pixel{X: -v.X, Y: -v.Y}
}

// scales the vector
func (v *Pixel) Scale(scalar int) {
	v.X *= scalar
	v.Y *= scalar
}

// executes an absolute value on the vectors coordinates
func (v *Pixel) Abs() {
	v.X = int(math.Abs(float64(v.X)))
	v.Y = int(math.Abs(float64(v.Y)))

}
