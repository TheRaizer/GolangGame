package display

import (
	"image"
	"image/color"
	"math"
)

// use Bresenhams line algorithm to generate the list of vectors that make an approximate line between two points
// this assumes we are using bitmaps and so vector coordinates must be integers as is the Pixel type
func bresenhamsLine(p1 Pixel, p2 Pixel) []Pixel {
	var start, end *Pixel
	switch {
	case p1.X < p2.X:
		start = &p1
		end = &p2
	case p1.X > p2.X:
		start = &p2
		end = &p1
	default:
		length := int(math.Abs(float64(p1.Y-p2.Y))) + 1
		pixels := make([]Pixel, length)

		yBase := int(math.Min(float64(p1.Y), float64(p2.Y)))

		for i := 0; i < length; i++ {
			pixels[i] = Pixel{X: p1.X, Y: yBase + i}
		}
		return pixels
	}

	numOfPixels := math.Abs(float64(end.X - start.X))
	pixels := make([]Pixel, int(numOfPixels))

	currPixel := 0

	var slope float64 = float64(end.Y-start.Y) / float64(end.X-start.X)

	for x := start.X; x < end.X; x++ {
		var y float64 = (slope * float64(x-start.X)) + float64(start.Y)
		roundedY := int(math.Round(y))

		pixels[currPixel] = Pixel{X: x, Y: roundedY}
		currPixel++
	}

	return pixels
}

// generates a black path from certain locations onto an image, this could be done without pointers
// and probably better without it but this is an example of passing values by reference (by pointer)
// mutates the given image
func GenPaths(img *image.Gray, locations []Pixel) {
	// calc dir vector as |a - b| from each point in order.
	// starting from a, scale the dir vector by a natural num i s.t. i stops when i * a is no longer on the img
	// each new vector a * i should be painted as a black path
	for i := 0; i < len(locations)-1; i++ {
		startPixel := locations[i]
		endPixel := locations[i+1]

		approxLine := bresenhamsLine(startPixel, endPixel)

		for _, pixel := range approxLine {
			// set certain pixels to be black
			(*img).SetGray(pixel.X, pixel.Y, color.Gray{Y: 255}) // could also do img.SetGray as a shortcut of dereferencing
		}
	}
}
