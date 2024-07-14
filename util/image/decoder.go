package image

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/TheRaizer/GolangGame/util"
)

func DecodePNG(name string) {
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

		chunkLength := convertBytesToUint[uint32](header[:4])
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
			png.ihdr = ihdrChunk
			fmt.Println(ihdrChunk)
		case "PLTE":
			if png.ihdr.colorType == 0 || png.ihdr.colorType == 4 {
				panic("PLTE chunk must not occur when color type 0 or 4")
			}
			plteChunk, err := parsePLTE(dataBuf)
			util.CheckErr(err)
			png.plte = plteChunk
			fmt.Println(plteChunk)
		case "IDAT":
			// ihdr chunk must have been read
			if png.ihdr == nil {
				panic("IHDR chunk should have been encountered before IDAT chunk")
			}
			// PLTE must appear for color type 3
			if png.plte == nil && png.ihdr.colorType == 3 {
				panic("PLTE chunk should have been encountered before IDAT chunk")
			}
			cmpltIdat = slices.Concat(cmpltIdat, dataBuf) // TODO: is there a faster method then concatenating
		case "IEND":
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

	scanlines, err := processIDAT(*png.ihdr, cmpltIdat)
	util.CheckErr(err)

	switch png.ihdr.colorType {
	case 3:
		colorMatrix := getColorIndices(scanlines, png.ihdr.bitDepth)
		fmt.Println(colorMatrix)
	}
}

func checkCRC(typeBuf []byte, dataBuf []byte, crcBuf []byte) {
	var crcInput []byte = slices.Concat(typeBuf, dataBuf)
	crc := Crc32(crcInput)

	if crc != convertBytesToUint[uint32](crcBuf) {
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

// decode the IHDR data into its separate data per
// http://www.libpng.org/pub/png/spec/1.2/PNG-Chunks.html
func decodeIHDR(data []byte) (*IHDR, error) {
	if len(data) != 13 {
		return nil, fmt.Errorf("IHDR data length must be 13")
	}

	width := convertBytesToUint[uint32](data[0:4])
	height := convertBytesToUint[uint32](data[4:8])

	ihdr := NewIHDR(width, height, data[8], data[9], data[10], data[11], data[12])

	return &ihdr, nil
}

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
// Returns the list of scanlines or error if any error occurs during processing.
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

	scanlines, err := defilterPixelData(buf, ihdr.width, ihdr.height, bpp)
	if err != nil {
		return nil, err
	}

	return scanlines, nil
}

func getPrevScanline(scanlines [][]byte, i int) []byte {
	var prevScanline []byte = nil
	if i > 0 {
		prevScanline = scanlines[i-1]
		return prevScanline
	}
	return nil
}

// When color type is 3 this will get the value of each index.
// Returns a 2D array of palette indices
func getColorIndices(rawScanlines [][]byte, bitDepth uint8) [][]uint8 {
	var paletteIndices [][]uint8 = make([][]uint8, len(rawScanlines))
	for r := 0; r < len(rawScanlines); r++ {
		scanline := rawScanlines[r]
		paletteIndices[r] = make([]uint8, len(scanline))
		j := 0 // tracks the byte we are evaluating
		c := 0 // the column of the image we are on
		for j < len(scanline) {
			if bitDepth < 8 {
				b := scanline[j]
				bytes, err := splitByte(b, int(bitDepth))
				util.CheckErr(err)

				for _, splitByte := range bytes {
					paletteIndices[r][c] = uint8(splitByte)
				}
				j++
			} else {
				bytes := scanline[j : j+int(bitDepth/8)]
				paletteIndices[r][c] = uint8(convertBytesToUint[uint32](bytes))

				j += int(bitDepth)
			}
			c++
		}

	}

	return paletteIndices
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
	stride := int(float32(width) * bpp)

	rawScanlines := make([][]byte, height)

	offset := 0 // points to the filter byte of the scanline
	i := 0
	bppRounded := int(math.Ceil(float64(bpp)))

	for i < len(rawScanlines) {
		filterType := uint8(decompressedData[0])
		filteredScanline := decompressedData[offset+1 : offset+1+stride]

		switch filterType {
		case 0: // None
			rawScanlines[i] = filteredScanline
		case 1: // Sub
			rawScanlines[i] = inverseSub(filteredScanline, bppRounded)
		// TODO: implement
		case 2: // Up
			rawScanlines[i] = inverseUp(filteredScanline, getPrevScanline(rawScanlines, i))
		case 3: // Average
			rawScanlines[i] = inverseAverage(filteredScanline, getPrevScanline(rawScanlines, i), bppRounded)
		case 4: // Paeth
			rawScanlines[i] = inversePaeth(filteredScanline, getPrevScanline(rawScanlines, i), bppRounded)
		}
		offset += 1 + stride
		i += 1
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
			rawScanline[i] = (filteredScanline[i] + rawScanline[i-bpp]) & 0xFF // & 0xFF == % 256 for unsigned integers
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
			rawScanline[i] = (filteredScanline[i] + rawPrevScanline[i]) & 0xFF
		}
	}
	return rawScanline
}

// The inverse of the average filter, which defilters a scanline that was filtered by the average algorithm.
// Returns the defiltered (raw) scanline
func inverseAverage(filteredScanline []byte, rawPrevScanline []byte, bpp int) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i := 0; i < len(filteredScanline); i++ {
		if i < bpp {
			// will run at least once when i = 0
			rawScanline[i] = filteredScanline[i]
		} else {
			// other case will have run when i = 0, we can be sure rawPrevScanline != nil by this point
			floored := byte(math.Floor(float64(rawScanline[i-bpp] + rawPrevScanline[i])))
			rawScanline[i] = (filteredScanline[i] + floored) & 0xFF
		}
	}
	return rawScanline
}

// The inverse of the Paeth filter, which defilters a scanline that was filtered by the Paeth algorithm.
// Returns the defiltered (raw) scanline
func inversePaeth(filteredScanline []byte, rawPrevScanline []byte, bpp int) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i := 0; i < len(filteredScanline); i++ {
		if i < bpp {
			rawScanline[i] = filteredScanline[i]
		} else {
			paethPrediction := paethPredictor(rawScanline[i-bpp], rawPrevScanline[i], rawPrevScanline[i-1])
			rawScanline[i] = (filteredScanline[i] + paethPrediction) & 0xFF
		}
	}
	return rawScanline
}

// Predict the value of a pixel based on the values of neighbouring pixels.
// left is the pixel to the left of the current
// up is the pixel above the current
// upLeft is the pixel to the upper left of the current
func paethPredictor(left, up, upLeft byte) byte {
	p := left + up - upLeft
	pLeftDist := math.Abs(float64(p - left))
	pUpDist := math.Abs(float64(p - up))
	pUpLeftDist := math.Abs(float64(p - upLeft))

	if pLeftDist <= pUpDist && pLeftDist <= pUpLeftDist {
		return left
	} else if pLeftDist <= pUpLeftDist {
		return up
	} else {
		return upLeft
	}
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

// this assumes data is formated in big endian
// length of val * 8 (bits) determines the type to use
func convertBytesToUint[T uint8 | uint16 | uint32 | uint64](val []byte) T {
	num := 0.0

	for i, byte := range val {
		num += float64(byte) * math.Pow(256, float64(len(val)-1-i))
	}

	return T(num)
}
