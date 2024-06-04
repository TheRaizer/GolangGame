package observer

type Observer[T any] interface {
	Emit(data T)
	GetID() string
}
