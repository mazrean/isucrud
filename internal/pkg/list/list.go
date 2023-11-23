package list

import "container/list"

type List[T any] struct {
	l *list.List
}

type Element[T any] struct {
	e *list.Element
}

func (e Element[T]) Value() (T, bool) {
	if e.e == nil {
		var v T
		return v, false
	}

	return e.e.Value.(T), true
}

func New[T any]() List[T] {
	return List[T]{
		l: list.New(),
	}
}

func (l List[T]) Back() Element[T] {
	return Element[T]{
		e: l.l.Back(),
	}
}

func (l List[T]) Front() Element[T] {
	return Element[T]{
		e: l.l.Front(),
	}
}

func (l List[T]) Init() List[T] {
	l.l.Init()
	return l
}

func (l List[T]) Len() int {
	return l.l.Len()
}

func (l List[T]) PushBack(v T) Element[T] {
	return Element[T]{
		e: l.l.PushBack(v),
	}
}

func (l List[T]) PushBackList(other List[T]) {
	l.l.PushBackList(other.l)
}
func (l List[T]) PushFront(v T) Element[T] {
	return Element[T]{
		e: l.l.PushFront(v),
	}
}

func (l List[T]) PushFrontList(other List[T]) {
	l.l.PushFrontList(other.l)
}

func (l List[T]) Remove(e Element[T]) T {
	return l.l.Remove(e.e).(T)
}

func (l List[T]) InsertAfter(v T, mark Element[T]) Element[T] {
	return Element[T]{
		e: l.l.InsertAfter(v, mark.e),
	}
}

func (l List[T]) InsertBefore(v T, mark Element[T]) Element[T] {
	return Element[T]{
		e: l.l.InsertBefore(v, mark.e),
	}
}

func (l List[T]) MoveAfter(e, mark Element[T]) {
	l.l.MoveAfter(e.e, mark.e)
}

func (l List[T]) MoveBefore(e Element[T], mark Element[T]) {
	l.l.MoveBefore(e.e, mark.e)
}

func (l List[T]) MoveToBack(e Element[T]) {
	l.l.MoveToBack(e.e)
}

func (l List[T]) MoveToFront(e Element[T]) {
	l.l.MoveToFront(e.e)
}
