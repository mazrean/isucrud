package list

type Queue[T any] struct {
	l List[T]
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{
		l: New[T](),
	}
}

func (q Queue[T]) Len() int {
	return q.l.Len()
}

func (q Queue[T]) Push(v T) {
	q.l.PushBack(v)
}

func (q Queue[T]) Pop() (T, bool) {
	e := q.l.Front()
	if e.e == nil {
		var v T
		return v, false
	}

	q.l.Remove(e)
	return e.Value()
}

func (q Queue[T]) Peek() (T, bool) {
	e := q.l.Front()
	if e.e == nil {
		var v T
		return v, false
	}

	return e.Value()
}

func (q Queue[T]) Clear() {
	q.l.Init()
}
