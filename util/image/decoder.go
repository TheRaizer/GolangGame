package image

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/TheRaizer/GolangGame/util"
)

// TODO: perhaps make this async using goroutines?
// Decodes a PNG file into a slice of RGBA values
// If PNG uses 16 bit depth RGB(A) then it is downscaled to 8 bit depth RGB(A)
func DecodePNG(name string) PNG {
	file, err := os.Open(name)

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	util.CheckErr(err)

	buffer := make([]byte, 50)
	_, err = file.Read(buffer)
	util.CheckErr(err)

	png := PNG{}
	var cmpltIdat []byte // the complete chunk of all the compressed IDAT data concatenated

	checkHeader(buffer[:8])
	// keep track of the current byte we are on
	// skip the first 8 bytes which were the header
	// chunks now start
	// read the first 4 bytes after out current byte which should be chunk length by specification
	// read into the buffer the current chunk by using this chunk length
	// the first 8 bytes contain chunk length and type, do specifics depending on chunk type
	// add the chunk length to the current chunk
	// reiterate through the next chunk

	// skip the header bytes
	var i int64 = 8
	for {
		header := make([]byte, 8)
		_, err := file.ReadAt(header, i)

		if isEOF(err) {
			break
		}

		chunkLength := util.ConvertBytesToUint[uint32](header[:4])
		typeBuf := header[4:8]
		chunkType := string(typeBuf)
		i += 8

		dataBuf := make([]byte, chunkLength)
		_, err = file.ReadAt(dataBuf, i)

		if isEOF(err) {
			break
		}

		if i == 16 && chunkType != "IHDR" {
			panic("first chunk must be IHDR chunk")
		}

		switch strings.ToUpper(chunkType) {
		case "IHDR":
			ihdrChunk, err := decodeIHDR(dataBuf)
			util.CheckErr(err)
			png.IHDR = ihdrChunk

			fmt.Println(png.colorType)
		case "PLTE":
			if png.IHDR.colorType == 0 || png.IHDR.colorType == 4 {
				panic("PLTE chunk must not occur when color type 0 or 4")
			}
			plteChunk, err := parsePLTE(dataBuf)
			util.CheckErr(err)
			png.PLTE = plteChunk
		case "IDAT":
			// ihdr chunk must have been read
			if png.IHDR == nil {
				panic("IHDR chunk should have been encountered before IDAT chunk")
			}
			// PLTE must appear for color type 3
			if png.PLTE == nil && png.IHDR.colorType == 3 {
				panic("PLTE chunk should have been encountered before IDAT chunk")
			}

			cmpltIdat = append(cmpltIdat, dataBuf...)
		case "IEND":
		case "TRNS":
			if png.IHDR.colorType == 4 || png.IHDR.colorType == 6 {
				panic("TRNS chunk cannot exist when color type is 4 or 6")
			}

			if png.PLTE == nil {
				panic("PLTE chunk must follow the TRNS chunk")
			}

			png.TRNS, err = parseTRNS(dataBuf, *png.PLTE, png.colorType)
			util.CheckErr(err)
		default:
			// check if the 5th bit (from LSB to MSB i.e. right to left) of the first byte is 1
			// 0 = critical, 1 = ancillary
			if typeBuf[0]&0b00100000 == 0 {
				panic("for critical chunk, encountered unknown chunk type " + chunkType)
			}
		}

		i += int64(chunkLength)

		crcBuf := make([]byte, 4)
		_, err = file.ReadAt(crcBuf, i)
		checkCRC(typeBuf, dataBuf, crcBuf)

		i += 4
	}

	rawScanlines, err := processIDAT(*png.IHDR, cmpltIdat)
	util.CheckErr(err)

	pixels := getPixels(rawScanlines, png)
	util.CheckErr(err)
	png.Data = &pixels

	return png
}

func packBytesToUint32(bytes [4]byte) uint32 {
	// combine the bytes into a single uint32
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) + (uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func grayscaleToRgba[T uint8 | uint16](pixel T, bitDepth uint8, alpha uint8) uint32 {
	pixel8 := rescaleToByte(bitDepth, pixel)
	rgbValue := packBytesToUint32([4]byte{pixel8, pixel8, pixel8, alpha})
	return rgbValue
}

// compresses
func rescaleToByte[T uint8 | uint16 | uint32 | uint64](bitDepth uint8, pixel T) byte {
	maxValue := math.Pow(2, float64(bitDepth)) - 1
	normalizedPixel := float64(pixel) / maxValue // get pixel between 0 and 1
	pixel8 := byte(normalizedPixel * 255)        // get pixel between 0 and 255

	return pixel8
}

func paletteIndicesToRgba[T uint8 | uint16](idx T, palette [][3]byte, alpha uint8) uint32 {
	rgb := palette[idx]
	rgbValue := packBytesToUint32([4]byte{rgb[0], rgb[1], rgb[2], alpha})
	return rgbValue
}

func checkCRC(typeBuf []byte, dataBuf []byte, crcBuf []byte) {
	var crcInput []byte = append(typeBuf, dataBuf...)
	crc := Crc32(crcInput)

	if crc != util.ConvertBytesToUint[uint32](crcBuf) {
		panic("CRC's did not match in a chunk")
	}

}

func isEOF(err error) bool {
	if err == io.EOF {
		return true
	} else if err != nil {
		panic(err)
	}

	return false
}

// TODO: move to IHDR.go call parseIHDR
// decode the IHDR data into its separate data per
// http://www.libpng.org/pub/png/spec/1.2/PNG-Chunks.html
func decodeIHDR(data []byte) (*IHDR, error) {
	if len(data) != 13 {
		return nil, fmt.Errorf("IHDR data length must be 13")
	}

	width := util.ConvertBytesToUint[uint32](data[0:4])
	height := util.ConvertBytesToUint[uint32](data[4:8])

	ihdr := NewIHDR(width, height, data[8], data[9], data[10], data[11], data[12])

	return &ihdr, nil
}

// TODO: move to PLTE.go
func parsePLTE(data []byte) (*PLTE, error) {
	if len(data)%3 != 0 {
		return nil, fmt.Errorf("PLTE data length must be divisible by 3")
	}

	plte := PLTE{palette: make([][3]byte, len(data)/3)}

	paletteIdx := -1

	for i, b := range data {
		if i%3 == 0 {
			paletteIdx++
		}
		plte.palette[paletteIdx][i%3] = b
	}

	return &plte, nil
}

// Processes the complete IDAT data into a list of scanlines.
// Returns the list of raw scanlines or error if any error occurs during processing.
func processIDAT(ihdr IHDR, data []byte) ([][]byte, error) {
	bReader := bytes.NewReader(data)
	z, err := zlib.NewReader(bReader)
	if err != nil {
		return nil, fmt.Errorf("Error when decompressing IDAT: %w", err)
	}

	defer z.Close()

	buf, err := io.ReadAll(z)
	if err != nil {
		return nil, fmt.Errorf("Error when reading IDAT: %w", err)
	}

	bpp, err := bytesPerPixel(ihdr.bitDepth, ihdr.colorType)

	if err != nil {
		return nil, err
	}

	scanlines, err := defilterPixelData(buf, ihdr.Width, ihdr.Height, bpp)
	if err != nil {
		return nil, err
	}

	return scanlines, nil
}

func getPrevScanline(scanlines [][]byte, i int) []byte {
	var prevScanline []byte = nil
	if i > 0 {
		prevScanline = scanlines[i-1]
	}
	return prevScanline
}

// Returns a matrix of the pixel values parsed from the given raw scanlines.
func getPixels(rawScanlines [][]byte, png PNG) []uint32 {
	var pixels []uint32
	for _, scanline := range rawScanlines {
		var scanlinePixels []uint32 = nil
		if png.bitDepth < 8 {
			scanlinePixels = fetchPixelsFromSubBytes(scanline, png)
		} else {
			scanlinePixels = fetchPixelsFromFullBytes(scanline, png)
		}
		pixels = append(pixels, scanlinePixels...)
	}
	return pixels
}

// fetch pixels from images using < 8 bit depth where pixel data is stored on the bit level
func fetchPixelsFromSubBytes(scanline []byte, png PNG) []uint32 {
	// the number of split bytes per pixel is 8 / the bit depth
	// in this function bit depth is expected to be 1, 2, or 4 so splitBpp is always a factor of 8
	scanlinePixels := make([]uint32, png.Width)

	c := 0
	for i := 0; i < len(scanline); i++ {
		b := scanline[i]

		bytes, err := splitByte(b, int(png.bitDepth))
		util.CheckErr(err)

		j := 0 // tracks the byte we are evaluating
		for j < len(bytes) {
			pixel, err := getPixelData(png, []byte{bytes[j]}) // when we split, each split byte has meaning
			util.CheckErr(err)

			scanlinePixels[c] = pixel
			c++
			j++
		}
	}

	return scanlinePixels
}

// fetch a row of pixels from images using 8 or 16 bit depth where pixel data is stored on the byte level
func fetchPixelsFromFullBytes(scanline []byte, png PNG) []uint32 {
	bpp, err := bytesPerPixel(png.bitDepth, png.colorType)
	util.CheckErr(err)

	scanlinePixels := make([]uint32, png.Width)

	j := 0 // tracks the byte we are evaluating
	c := 0
	for j < len(scanline) {
		bytes := scanline[j : j+int(bpp)]
		var pixel uint32 = 0
		if png.bitDepth == 16 {
			// every 2 bytes has meaninful data so we compress every 2 bytes into a single byte (as to not handle 16 bit depth)
			// but rather scale down to 8 bit depth
			pixel, err = getPixelData(png, compress16BitDepthBytes(bytes, int(bpp), png.bitDepth))
		} else {
			// in this case bitDepth is 8 so each byte has meaningful data
			pixel, err = getPixelData(png, bytes)
		}
		util.CheckErr(err)
		scanlinePixels[c] = pixel
		c++
		j += int(bpp)
	}

	return scanlinePixels
}

// compresses a slice of EVEN bytes with bit depth of 16 (every 2 bytes contains meaningful data)
// into a bit depth of 8 (every byte contains meaninful data) via normalization from a uint16 to uint8
func compress16BitDepthBytes(bytes []byte, bpp int, bitDepth uint8) []byte {
	if len(bytes)%2 != 0 {
		panic("when compressing from 16 bit, the slice of uncompressed bytes must have an even length")
	}

	var compressedBytes []byte = make([]byte, len(bytes)/2)
	var prev byte = 0
	bIdx := 0
	for i, b := range bytes {
		if i%bpp == 0 && i != 0 {
			p := util.ConvertBytesToUint[uint16]([]byte{prev, b})
			compressedBytes[bIdx] = rescaleToByte(bitDepth, p)
			bIdx++
		}
		prev = b
	}

	return compressedBytes
}

// Converts a slice of bytes to RGBA data
// returns a single uint32 representing that pixel's RGBA data
func getPixelData(png PNG, bytes []byte) (uint32, error) {
	switch png.colorType {
	case 0:
		var alpha uint8 = 255
		if png.TRNS != nil {
			if uint16(bytes[0]) == png.gray {
				alpha = 0
			}
		}
		return grayscaleToRgba(bytes[0], png.bitDepth, alpha), nil
	case 2:
		color := [4]byte{
			bytes[0],
			bytes[1],
			bytes[2],
			255,
		}

		if png.TRNS != nil {
			if uint16(bytes[0]) == png.rgb[0] && uint16(bytes[1]) == png.rgb[1] && uint16(bytes[2]) == png.rgb[2] {
				color[3] = 0 // make transparent
			}
		}
		// every 3 pixelData's represents RGB of a single pixel
		return packBytesToUint32(
			color,
		), nil
	case 3:
		if png.PLTE == nil {
			return 0, fmt.Errorf("Should have PLTE chunk with color type 3")
		}

		paletteIdx := uint8(bytes[0])
		var alpha uint8 = 255
		if png.TRNS != nil && int(paletteIdx) < len(png.alphas) {
			alpha = png.alphas[paletteIdx]
		}

		return paletteIndicesToRgba(paletteIdx, png.palette, alpha), nil
	case 4:
		// every 2 pixelData's represents gray scale and alpha of a single pixel
		pixel8 := bytes[0]
		return packBytesToUint32(
			[4]byte{
				pixel8,
				pixel8,
				pixel8,
				rescaleToByte(png.bitDepth, bytes[1]),
			},
		), nil

	case 6:
		// every 4 pixelData's represents 3 RGB channels and an alpha channel
		return packBytesToUint32(
			[4]byte{
				bytes[0],
				bytes[1],
				bytes[2],
				bytes[3],
			},
		), nil
	default:
		return 0, fmt.Errorf("Error color type %d is invalid", png.colorType)
	}
}

// Splits a given byte into n different sub-bytes aligned to the LSB.
// Returns the sub-bytes in order from MSB to LSB of the initial byte.
func splitByte(b byte, n int) ([]byte, error) {
	if n < 1 || n > 8 {
		return nil, fmt.Errorf("Bit depth %d must be from 1-8", n)
	}

	if n != 1 && n%2 != 0 {
		return nil, fmt.Errorf("Bit depth %d must be divisible by 2", n)
	}

	bytes := make([]byte, int(8/n))
	mask := byte((1 << n) - 1)

	for i := len(bytes) - 1; i >= 0; i-- {
		// right shift as to align each split byte to the LSB
		// eg. 00100101 with n=4 -> splits to 00000010 and 00000101
		bytes[i] = (b & mask) >> (n * (len(bytes) - 1 - i))
		mask = mask << n
	}

	return bytes, nil
}

// Defilters each scanline according to their specified filter, and returns a 2D slice of the defiltered (raw)
// scanlines with the filter type ommited.
func defilterPixelData(decompressedData []byte, width uint32, height uint32, bpp float32) ([][]byte, error) {
	// one stride corresponds to the length of one scanline excluding filter byte (one row of the image)
	stride := int(math.Ceil(float64(width) * float64(bpp)))

	bppRounded := int(math.Ceil(float64(bpp)))
	rawScanlines := make([][]byte, height)
	offset := 0 // points to the filter byte of the scanline
	i := 0

	for i < len(rawScanlines) {
		filterType := uint8(decompressedData[offset])
		filteredScanline := decompressedData[offset+1 : offset+1+stride]

		switch filterType {
		case 0: // None
			rawScanlines[i] = filteredScanline
		case 1: // Sub
			rawScanlines[i] = inverseSub(filteredScanline, bppRounded)
		case 2: // Up
			rawScanlines[i] = inverseUp(filteredScanline, getPrevScanline(rawScanlines, i))
		case 3: // Average
			rawScanlines[i] = inverseAverage(filteredScanline, getPrevScanline(rawScanlines, i), bppRounded)
		case 4: // Paeth
			rawScanlines[i] = inversePaeth(filteredScanline, getPrevScanline(rawScanlines, i), bppRounded)
		default:
			return nil, fmt.Errorf("Unexpected filter type %d", filterType)
		}
		offset += 1 + stride
		i++
	}

	if offset != len(decompressedData) {
		return nil, fmt.Errorf("Did not iterate correctly through compressed data")
	}

	return rawScanlines, nil
}

// INVERSE FILTERS
// http://www.libpng.org/pub/png/spec/1.2/PNG-Filters.html

// The inverse of the sub filter, which defilters a scanline that was filtered by the sub algorithm.
// Returns the defiltered (raw) scanline
func inverseSub(filteredScanline []byte, bpp int) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i := 0; i < len(filteredScanline); i++ {
		if i < bpp {
			rawScanline[i] = filteredScanline[i]
		} else {
			rawScanline[i] = filteredScanline[i] + rawScanline[i-bpp]
		}
	}
	return rawScanline
}

// The inverse of the up filter, which defilters a scanline that was filtered by the up algorithm.
// Returns the defiltered (raw) scanline
func inverseUp(filteredScanline []byte, rawPrevScanline []byte) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i, x := range filteredScanline {
		if rawPrevScanline == nil {
			rawScanline[i] = x
		} else {
			rawScanline[i] = filteredScanline[i] + rawPrevScanline[i]
		}
	}
	return rawScanline
}

// The inverse of the average filter, which defilters a scanline that was filtered by the average algorithm.
// Returns the defiltered (raw) scanline
func inverseAverage(filteredScanline []byte, rawPrevScanline []byte, bpp int) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i := 0; i < len(filteredScanline); i++ {
		var left, up byte
		if i >= bpp {
			left = rawScanline[i-bpp]
		}
		if rawPrevScanline != nil {
			up = rawPrevScanline[i]
		}

		// convert to float to avoid overflow and satisfy type constraint
		floored := int(math.Floor((float64(left) + float64(up)) / 2))

		// mod 256 to fit into a single byte
		rawScanline[i] = byte((int(filteredScanline[i]) + floored) % 256)
	}
	return rawScanline
}

// The inverse of the Paeth filter, which defilters a scanline that was filtered by the Paeth algorithm.
// Returns the defiltered (raw) scanline
func inversePaeth(filteredScanline []byte, rawPrevScanline []byte, bpp int) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i := 0; i < len(filteredScanline); i++ {
		var a, b, c byte
		if rawPrevScanline != nil {
			b = rawPrevScanline[i]
			if i >= bpp {
				c = rawPrevScanline[i-bpp]
			}
		}
		if i >= bpp {
			a = rawScanline[i-bpp]
		}
		rawScanline[i] = filteredScanline[i] + paethPred(a, b, c)
	}
	return rawScanline
}

func paethPred(a, b, c byte) byte {
	// cast to int to avoid overflow
	p := int(a) + int(b) - int(c)
	pa := abs(p - int(a))
	pb := abs(p - int(b))
	pc := abs(p - int(c))

	if pa <= pb && pa <= pc {
		return a
	} else if pb <= pc {
		return b
	} else {
		return c
	}
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Returns the number of bytes per pixel of a PNG (NOT including the filter byte).
// The multiplicative constant represents the number of color channels or additional information (eg. an alpha channel)
func bytesPerPixel(bitDepth, colorType uint8) (float32, error) {
	bd := float32(bitDepth)
	switch colorType {
	case 0: // Grayscale
		return bd / 8, nil
	case 2: // Truecolor (3 color channels)
		return 3 * (bd / 8), nil
	case 3: // Indexed-color
		return bd / 8, nil
	case 4: // Grayscale with alpha
		return 2 * (bd / 8), nil
	case 6: // Truecolor with alpha
		return 4 * (bd / 8), nil
	default:
		return 0, fmt.Errorf("unsupported color type: %d", colorType)
	}
}

// checks the 8 byte header and ensures that they match the PNG specification id
// per http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html
func checkHeader(header []byte) error {
	expectedHeader := []byte{137, 80, 78, 71, 13, 10, 26, 10}

	if !bytes.Equal(header, expectedHeader) {
		return fmt.Errorf("Not a PNG file")
	}
	return nil
}
