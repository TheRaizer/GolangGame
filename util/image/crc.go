package image

// NOTE: PNG version from http://www.libpng.org/pub/png/spec/1.2/PNG-CRCAppendix.html
// crc initialized to 1's, data from each byte is processed from LSBit to MSBit, then ones complement is taken

/* Table of CRCs of all 8-bit messages. */
var crcTable [256]uint32

/* Flag: has the table been computed? Initially false. */
var crcTableComputed bool = false

/* Make the table for a fast CRC. */
func makeCRCTable() {
	for n := 0; n < 256; n++ {
		c := uint32(n)
		for k := 0; k < 8; k++ {
			if c&1 == 1 {
				c = 0xedb88320 ^ (c >> 1) // 0xedb88320 represents the polynomial
			} else {
				c = c >> 1
			}
		}
		crcTable[n] = c
	}
	crcTableComputed = true
}

/*
Update a running CRC with the bytes buf[0..len-1]--the CRC

	should be initialized to all 1's, and the transmitted value
	is the 1's complement of the final running CRC (see the
	crc() routine below)).
*/
func updateCRC(crc uint32, buf []byte) uint32 {
	c := crc
	if !crcTableComputed {
		makeCRCTable()
	}
	for _, b := range buf {
		c = crcTable[(c^uint32(b))&0xff] ^ (c >> 8)
	}
	return c
}

func Crc32(buf []byte) uint32 {
	return updateCRC(0xffffffff, buf) ^ 0xffffffff
}
