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

	processIDAT(*png.ihdr, cmpltIdat)
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
// returns error if any error occurs during processing.
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

	scanlines, err := defilterPixelData(buf, ihdr.width, bpp)
	if err != nil {
		return nil, err
	}

	return scanlines, nil
}

// Defilters each scanline according to their specified filter, and returns a 2D slice of the defiltered (raw)
// scanlines with the filter type ommited.
func defilterPixelData(decompressedData []byte, width uint32, bpp float32) ([][]byte, error) {
	// one stride corresponds to the length of one scanline excluding filter byte (one row of the image)
	stride := int(float32(width) * bpp)

	rawScanlines := make([][]byte, int(len(decompressedData)/int(stride+1)))

	offset := 0 // points to the filter byte of the scanline
	i := 0

	for i < len(rawScanlines) {
		filterType := uint8(decompressedData[0])
		filteredScanline := decompressedData[offset+1 : offset+1+stride]
		fmt.Println(filterType)

		switch filterType {
		case 0: // None
			rawScanlines[i] = filteredScanline
		case 1: // Sub
			bppRounded := int(math.Ceil(float64(bpp)))
			rawScanlines[i] = inverseSub(filteredScanline, bppRounded)
		// TODO: implement
		case 2: // Up
		case 3: // Average
		case 4: // Paeth
		}
		offset += 1 + stride
		i += 1
	}

	if offset != len(decompressedData) {
		return nil, fmt.Errorf("Did not iterate correctly through compressed data")
	}

	return rawScanlines, nil
}

// The inverse of the sub filter, which defilters a scanline that was filtered by the sub algorithm.
// bpp is the bytes per pixel rounded up to 1 as per the docs.
// filteredScanline is the sub filtered scanline ommiting the filter type byte.
// http://www.libpng.org/pub/png/spec/1.2/PNG-Filters.html#Filter-type-1-Sub
func inverseSub(filteredScanline []byte, bpp int) []byte {
	rawScanline := make([]byte, len(filteredScanline))
	for i, x := range filteredScanline {
		if i < bpp {
			rawScanline[i] = x
		} else {
			rawScanline[i] = filteredScanline[i] + rawScanline[i-bpp] // NOTE: the docs say this mod 256 to stop overflow?
		}
	}
	return rawScanline
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
