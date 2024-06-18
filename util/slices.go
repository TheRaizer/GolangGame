package util

type Slice[T any] []T

func (slice Slice[T]) RemoveIdx(idx int) Slice[T] {
	slice[idx] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]

	return slice
}
