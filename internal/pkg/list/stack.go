package list

type Stack[T any] struct {
	l List[T]
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		l: New[T](),
	}
}

func (s Stack[T]) Len() int {
	return s.l.Len()
}

func (s Stack[T]) Push(v T) {
	s.l.PushBack(v)
}

func (s Stack[T]) Pop() (T, bool) {
	e := s.l.Back()
	if e.e == nil {
		var v T
		return v, false
	}

	s.l.Remove(e)
	return e.Value()
}

func (s Stack[T]) Peek() (T, bool) {
	e := s.l.Back()
	if e.e == nil {
		var v T
		return v, false
	}

	return e.Value()
}

func (s Stack[T]) Clear() {
	s.l.Init()
}
