package util

type Slice[T any] []T

func (slice Slice[T]) RemoveIdx(idx int) Slice[T] {
	slice[idx] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]

	return slice
}

func ConcatSlices[T any](slices ...[]T) []T {
	// calculate length of new concatenated slice
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}

	// create new
	tmp := make([]T, totalLen)

	// copy each slice into their correct position in tmp
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}
