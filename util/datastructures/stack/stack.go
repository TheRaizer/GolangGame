package datastructures

type Stack[T comparable] struct {
	items []T
}

func NewStack[T comparable]() Stack[T] {
	return Stack[T]{}
}

func (stack *Stack[T]) Peek() T {
	if stack.IsEmpty() {
		panic("cannot peek an empty stack")
	}
	return stack.items[len(stack.items)-1]
}

func (stack *Stack[T]) Push(item T) {
	stack.items = append(stack.items, item)
}

func (stack *Stack[T]) Pop() T {
	if stack.IsEmpty() {
		panic("cannot pop off an empty stack")
	}
	item := stack.Peek()
	stack.items = stack.items[0 : len(stack.items)-1]
	return item
}

func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.items) == 0
}
