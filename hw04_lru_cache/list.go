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
	front *ListItem
	back  *ListItem
}

func (l *list) insertAfter(item *ListItem, newItem *ListItem) {
	newItem.Prev = item

	if item.Next == nil {
		newItem.Next = nil
		l.back = newItem
	} else {
		newItem.Next = item.Next
		item.Next.Prev = newItem
	}

	item.Next = newItem
}

func (l *list) insertBefore(item *ListItem, newItem *ListItem) {
	newItem.Next = item

	if item.Prev == nil {
		newItem.Prev = nil
		l.front = newItem
	} else {
		newItem.Prev = item.Prev
		item.Prev.Next = newItem
	}

	item.Prev = newItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{
		Value: v,
	}

	if l.front == nil {
		l.front = &newItem
		l.back = &newItem
		newItem.Prev = nil
		newItem.Next = nil
	} else {
		l.insertBefore(l.front, &newItem)
	}

	l.len++

	return &newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{
		Value: v,
	}
	if l.back == nil {
		return l.PushFront(v)
	}
	l.insertAfter(l.back, &newItem)
	l.len++
	return &newItem
}

func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}
	if item.Prev == nil {
		l.front = item.Next
	} else {
		item.Prev.Next = item.Next
	}

	if item.Next == nil {
		l.back = item.Prev
	} else {
		item.Next.Prev = item.Prev
	}

	l.len--
}

func (l *list) MoveToFront(item *ListItem) {
	if item == nil || item.Prev == nil || l.len == 1 {
		return
	}

	if item.Next == nil {
		item.Prev.Next = nil
		l.back = item.Prev
	} else {
		item.Prev.Next = item.Next
		item.Next.Prev = item.Prev
	}

	item.Prev = nil
	item.Next = l.front

	l.front.Prev = item
	l.front = item
}

func NewList() List {
	return new(list)
}
