package image

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/TheRaizer/GolangGame/util"
)

// NOTE: first 8 bytes are an identifier for png
// now chunks start, first chunk is IHDR chunk
// next 4 bytes represents the chunk length
// next 4 bytes represents the chunk type
// chunk data
// 4 bytes of CRC

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

		dataLength := convertBytesToUint[uint32](header[:4])
		typeBuf := header[4:8]
		chunkType := string(typeBuf)
		i += 8

		dataBuf := make([]byte, dataLength)
		_, err = file.ReadAt(dataBuf, i)

		if isEOF(err) {
			break
		}

		switch strings.ToUpper(chunkType) {
		case "IHDR":
			ihdrChunk := decodeIHDRChunk(dataBuf)
			fmt.Println(ihdrChunk)
		case "IDAT":

		case "IEND":
		default:
			// check if the 5th bit of the first byte is 1
			v := typeBuf[0] & 16 // use 16's bit representation as a mask
			if v == 0 {
				panic("for critical chunk, encountered unknown chunk type " + chunkType)
			}
		}

		i += int64(dataLength)

		crcBuf := make([]byte, 4)
		_, err = file.ReadAt(crcBuf, i)
		checkCRC(typeBuf, dataBuf, crcBuf)

		i += 4
	}
}

func checkCRC(typeBuf []byte, dataBuf []byte, crcBuf []byte) {
	var crcInput []byte = util.ConcatSlices(typeBuf, dataBuf)
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

// decode the IHDR chunk into its separate data per
// http://www.libpng.org/pub/png/spec/1.2/PNG-Chunks.html
func decodeIHDRChunk(chunk []byte) IHDRChunk {
	if len(chunk) != 13 {
		panic("IHDR chunk length must be 13")
	}

	width := convertBytesToUint[uint32](chunk[0:4])
	height := convertBytesToUint[uint32](chunk[4:8])

	return NewIHDRChunk(width, height, chunk[8], chunk[9], chunk[10], chunk[11], chunk[12])
}

// checks the 8 byte header and ensures that they match the PNG specification id
// per http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html
func checkHeader(header []byte) {
	expectedHeader := []byte{137, 80, 78, 71, 13, 10, 26, 10}

	if !bytes.Equal(header, expectedHeader) {
		panic("not a png file!")
	}
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
