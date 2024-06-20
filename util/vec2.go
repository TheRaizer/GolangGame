package util

type Vec2[T int8 | int16 | int32 | int64 | float32 | float64] struct {
	X, Y T
}

func (vec *Vec2[T]) Add(otherVec *Vec2[T]) *Vec2[T] {
	vec.X += otherVec.X
	vec.Y += otherVec.Y

	return vec
}

func (vec *Vec2[T]) Multiply(num T) *Vec2[T] {
	vec.X *= num
	vec.Y *= num

	return vec
}

func (vec *Vec2[T]) Divide(num T) *Vec2[T] {
	vec.X /= num
	vec.Y /= num

	return vec
}
