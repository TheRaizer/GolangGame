package image

import "fmt"

type IHDR struct {
	width  uint32 // width of PNG
	height uint32 // height of PNG

	// Bit depth is a single-byte integer giving the number of bits per sample or per palette index (not per pixel)
	// Valid values are 1, 2, 4, 8, and 16, although not all values are allowed for all color types.
	bitDepth uint8

	// Color type codes represent sums of the following values: 1 (palette used), 2 (color used), and 4 (alpha channel used).
	// Valid values are 0, 2, 3, 4, and 6.
	colorType uint8

	// Indicates the method used to compress the image data.
	// At present, only compression method 0 (deflate/inflate compression with a sliding window of at most 32768 bytes) is defined.
	compressionMethod uint8

	// Indicates the preprocessing method applied to the image data before compression.
	// At present, only filter method 0 (adaptive filtering with five basic filter types) is defined
	filterMethod uint8

	// Indicates the transmission order of the image data
	// Two values are currently defined: 0 (no interlace) or 1 (Adam7 interlace).
	interlaceMethod uint8
}

func NewIHDR(
	width uint32,
	height uint32,
	bitDepth uint8,
	colorType uint8,
	compressionMethod uint8,
	filterMethod uint8,
	interlaceMethod uint8,
) IHDR {
	if compressionMethod != 0 {
		panic("unsupported compression method " + string(compressionMethod) + " was found!")
	}
	if filterMethod != 0 {
		panic("unsupported filter method " + string(filterMethod) + " was found!")
	}

	err := checkColorType(colorType)
	err = checkBitDepth(bitDepth, colorType)

	if err != nil {
		panic(err)
	}

	return IHDR{
		width,
		height,
		bitDepth,
		colorType,
		compressionMethod,
		filterMethod,
		interlaceMethod,
	}
}

func checkColorType(colorType uint8) error {
	if colorType != 0 && colorType != 2 && colorType != 3 && colorType != 4 && colorType != 6 {
		return fmt.Errorf(
			"Color type %d is an invalid integer. Must be: 0, 2, 3, 4, or 6",
			colorType,
		)
	}
	return nil
}

func checkBitDepth(bitDepth, colorType uint8) error {
	if bitDepth != 1 && bitDepth%2 != 0 {
		return fmt.Errorf(
			"Bit depth %d is an invalid integer. Must be: 1, 2, 4, 8, or 16",
			bitDepth,
		)
	}

	if (colorType == 4 || colorType == 6 || colorType == 2) && (bitDepth != 8 && bitDepth != 16) {
		return fmt.Errorf(
			"Color type 2, with invalid bit depth: %d. Must be: 8 or 16",
			bitDepth,
		)
	}

	if colorType == 3 && (bitDepth != 1 && bitDepth != 2 && bitDepth != 4 && bitDepth != 8) {
		return fmt.Errorf(
			"Color type 3, with invalid bit depth: %d. Must be: 1, 2, 4 or 8",
			bitDepth,
		)
	}

	return nil
}
