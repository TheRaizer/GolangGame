package image

type PNG struct {
	*IHDR
	*PLTE
	Data *[]uint32
}
