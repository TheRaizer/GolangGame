package image

import (
	"fmt"

	"github.com/TheRaizer/GolangGame/util"
)

// The PLTE chunk contains from 1 to 256 palette entries, each a three-byte series
type TRNS struct {
	// color type 3
	// each index in the alphas corresponds to the respective index in the PLTE palette.
	// value is the alpha value for that palette index.
	// length must be less than or equal to length of the palette.
	// if less then palette length, any other palette indices will have alpha 255.
	alphas []uint8

	// color type 0
	// this gray scale value will be treated as transparent
	gray uint16

	// color type 2
	// this rgb color will be treated as transparent
	rgb [3]uint16
}

func parseTRNS(data []byte, plte PLTE, colorType uint8) (*TRNS, error) {
	if len(data) > len(plte.palette) {
		return nil, fmt.Errorf("alpha channels cannot exceed palette length")
	}

	trns := TRNS{}

	switch colorType {
	case 3:
		trns.alphas = data
	case 0:
		trns.gray = util.ConvertBytesToUint[uint16](data)
	case 2:
		red := util.ConvertBytesToUint[uint16](data[:2])
		green := util.ConvertBytesToUint[uint16](data[2:4])
		blue := util.ConvertBytesToUint[uint16](data[4:6])

		trns.rgb = [3]uint16{red, green, blue}
	}

	return &trns, nil
}
