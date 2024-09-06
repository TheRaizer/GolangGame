package util

import "math"

// this assumes data is formated in big endian
// length of val * 8 (bits) determines the type to use
func ConvertBytesToUint[T uint8 | uint16 | uint32 | uint64](val []byte) T {
	num := 0.0

	for i, byte := range val {
		num += float64(byte) * math.Pow(256, float64(len(val)-1-i))
	}

	return T(num)
}
