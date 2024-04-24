package hw04lrucache

import (
	"fmt"
	"strings"
)

type List interface {
	fmt.Stringer
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	fmt.Stringer
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
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
	newNode := new(ListItem)
	newNode.Value = v
	newNode.Next = l.front
	if l.front != nil {
		l.front.Prev = newNode
	}
	if l.back == nil {
		l.back = newNode
	}
	l.front = newNode
	l.len++
	return newNode
}

func (l *list) PushBack(v interface{}) *ListItem {
	newNode := new(ListItem)
	newNode.Value = v
	newNode.Prev = l.back
	if l.back != nil {
		l.back.Next = newNode
	}
	if l.front == nil {
		l.front = newNode
	}
	l.back = newNode
	l.len++
	return newNode
}

func (l *list) Remove(i *ListItem) {
	l.len--
	if l.len < 1 {
		l.front = nil
		l.back = nil
		return
	}
	if l.front == i {
		l.front = i.Next
		l.front.Prev = nil
		return
	}
	if l.back == i {
		l.back = i.Prev
		l.back.Next = nil
		return
	}
	prev := i.Prev
	next := i.Next
	prev.Next = next
	next.Prev = prev
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	if l.back == i {
		l.back = l.back.Prev
	}
	prev := i.Prev
	next := i.Next
	prev.Next = next
	if next != nil {
		next.Prev = prev
	}
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

func (l *list) String() string {
	arr := make([]string, l.len)
	node := l.front
	for i := 0; node != nil; i, node = i+1, node.Next {
		arr[i] = node.String()
	}
	return "[" + strings.Join(arr, ", ") + "]"
}

func (li *ListItem) String() string {
	return fmt.Sprintf("%v", li.Value)
}
