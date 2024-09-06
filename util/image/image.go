package image

type PNG struct {
	*IHDR
	*PLTE
	*TRNS
	Data *[]uint32
}
