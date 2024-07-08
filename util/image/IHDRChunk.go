package image

type IHDRChunk struct {
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

func NewIHDRChunk(
	width uint32,
	height uint32,
	bitDepth uint8,
	colorType uint8,
	compressionMethod uint8,
	filterMethod uint8,
	interlaceMethod uint8,
) IHDRChunk {
	if compressionMethod != 0 {
		panic("unsupported compression method " + string(compressionMethod) + " was found!")
	}
	if filterMethod != 0 {
		panic("unsupported filter method " + string(filterMethod) + " was found!")
	}

	return IHDRChunk{
		width,
		height,
		bitDepth,
		colorType,
		compressionMethod,
		filterMethod,
		interlaceMethod,
	}

}
