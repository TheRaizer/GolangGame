package observer

type Subject[T any] struct {
	observers map[string]Observer[T]
}

func (s *Subject[T]) registerObserver(observer Observer[T]) {
	s.observers[observer.GetID()] = observer
}

func (s *Subject[T]) deregisterObserver(observer Observer[T]) {
	delete(s.observers, observer.GetID())
}

func (s *Subject[T]) notifyAll(data T) {
	for _, val := range s.observers {
		val.Emit(data)
	}
}
