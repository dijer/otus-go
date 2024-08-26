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
	head *ListItem
	tail *ListItem
	len  int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{
		Value: v,
		Next:  l.head,
	}

	if l.len == 0 {
		l.head = newNode
		l.tail = newNode
	} else {
		l.head.Prev = newNode
		l.head = newNode
	}

	l.len++

	return newNode
}

func (l *list) PushBack(v interface{}) *ListItem {
	newNode := &ListItem{
		Value: v,
		Prev:  l.tail,
	}

	if l.len == 0 {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.Next = newNode
		l.tail = newNode
	}

	l.len++

	return newNode
}

func (l *list) Remove(i *ListItem) {
	prevNode := i.Prev
	nextNode := i.Next

	if prevNode != nil {
		prevNode.Next = nextNode
	} else {
		l.head = nextNode
	}

	if nextNode != nil {
		nextNode.Prev = prevNode
	} else {
		l.tail = prevNode
	}

	i.Prev = nil
	i.Next = nil

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	headNode := l.head
	prevNode := i.Prev
	nextNode := i.Next

	if prevNode != nil {
		prevNode.Next = nextNode
	}

	if nextNode != nil {
		nextNode.Prev = prevNode
	}

	i.Next = headNode

	l.head = i
}

func NewList() List {
	return new(list)
}
