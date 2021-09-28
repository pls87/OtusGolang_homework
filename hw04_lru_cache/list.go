package hw04lrucache

import (
	"sync"
)

var mutex sync.Mutex

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

func (li ListItem) IsLast() bool {
	return li.Next == nil
}

func (li ListItem) IsFirst() bool {
	return li.Prev == nil
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	mutex.Lock()
	defer mutex.Unlock()

	newLi := &ListItem{Value: v, Next: l.front, Prev: nil}

	if !l.push2Empty(newLi) {
		l.front.Prev = newLi
		l.front = newLi
		l.len++
	}

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	mutex.Lock()
	defer mutex.Unlock()

	newLi := &ListItem{Value: v, Next: nil, Prev: l.back}

	if !l.push2Empty(newLi) {
		l.back.Next = newLi
		l.back = newLi
		l.len++
	}

	return l.back
}

func (l *list) Remove(li *ListItem) {
	mutex.Lock()
	defer mutex.Unlock()

	if l.removeBack(li) || l.removeFront(li) {
		return
	}

	li.Next.Prev, li.Prev.Next = li.Prev, li.Next
	l.len--
}

func (l *list) MoveToFront(li *ListItem) {
	mutex.Lock()
	defer mutex.Unlock()

	switch {
	case li.IsFirst():
		return
	case li.IsLast():
		l.back = l.back.Prev
		l.back.Next = nil
	default:
		li.Next.Prev, li.Prev.Next = li.Prev, li.Next
	}

	li.Next, li.Prev, l.front = l.front, nil, li
}

func (l *list) removeBack(li *ListItem) bool {
	if li.IsLast() {
		l.back = l.back.Prev
		l.back.Next = nil
		l.len--
		return true
	}
	return false
}

func (l *list) removeFront(li *ListItem) bool {
	if li.IsFirst() {
		l.front = l.front.Next
		l.front.Prev = nil
		l.len--
		return true
	}
	return false
}

func (l *list) push2Empty(li *ListItem) bool {
	if l.Len() == 0 {
		l.back, l.front = li, li
		l.len++
		return true
	}
	return false
}

func NewList() List {
	return &list{len: 0, front: nil, back: nil}
}
