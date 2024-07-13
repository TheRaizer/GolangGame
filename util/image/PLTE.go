package image

// The PLTE chunk contains from 1 to 256 palette entries, each a three-byte series
type PLTE struct {
	// Store each palette entry which holds a three-byte series.
	// The three bytes represent RGB values respectively
	palette [][3]byte
}
