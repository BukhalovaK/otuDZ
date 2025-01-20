package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	first *ListItem
	last  *ListItem
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.first
}

func (l list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{Value: v}
	prevFirst := l.first
	item.Next = prevFirst
	l.first = &item
	if l.len == 0 {
		l.last = l.first
	} else {
		prevFirst.Prev = &item
	}
	l.len++

	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v}
	prevLast := l.last
	if l.len == 0 {
		l.first = &item
	} else {
		prevLast.Next = &item
		item.Prev = prevLast
	}
	l.last = &item

	l.len++

	return &item
}

func (l *list) Remove(i *ListItem) {
	iNext := i.Next
	iPrev := i.Prev

	if iPrev != nil {
		iPrev.Next = iNext
	} else {
		l.first = iNext
	}

	if iNext != nil {
		iNext.Prev = iPrev
	} else {
		l.last = iPrev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
