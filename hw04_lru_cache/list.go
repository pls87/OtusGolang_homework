package hw04lrucache

import (
	"sync"
)

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
	mu    *sync.Mutex
	len   int
	front *ListItem
	back  *ListItem
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
	l.mu.Lock()
	defer l.mu.Unlock()

	newLi := &ListItem{Value: v, Next: l.front, Prev: nil}

	if !l.push2Empty(newLi) {
		if l.front == nil {
			_ = l.front
		}
		l.front.Prev = newLi
		l.front = newLi
		l.len++
	}

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLi := &ListItem{Value: v, Next: nil, Prev: l.back}

	if !l.push2Empty(newLi) {
		if l.back == nil {
			_ = l.back
		}
		l.back.Next = newLi
		l.back = newLi
		l.len++
	}

	return l.back
}

func (l *list) Remove(li *ListItem) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.removeBack(li) {
		return
	}

	if l.removeFront(li) {
		return
	}

	li.Next.Prev, li.Prev.Next = li.Prev, li.Next
	l.len--
}

func (l *list) MoveToFront(li *ListItem) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch {
	case li.IsFirst():
		return
	case l.len == 2:
		l.back.Prev = nil
		l.back.Next = l.front
		l.front.Next = nil
		l.front.Prev = l.back
		l.front, l.back = l.back, l.front
	case li.IsLast():
		prev := l.back.Prev
		l.back.Prev = nil
		l.back.Next = l.front
		l.front.Prev = l.back
		l.front = l.back
		l.back = prev
		l.back.Next = nil
	default:
		prev, next := li.Prev, li.Next
		li.Next.Prev = prev
		li.Prev.Next = next
		l.front.Prev = li
		li.Next = l.front
		li.Prev = nil
		l.front = li
	}
}

func (l *list) removeBack(li *ListItem) bool {
	if li.IsLast() {
		l.back = l.back.Prev
		if l.back != nil {
			l.back.Next = nil
		} else {
			l.front = nil
		}
		l.len--
		return true
	}
	return false
}

func (l *list) removeFront(li *ListItem) bool {
	if li.IsFirst() {
		l.front = l.front.Next
		if l.front != nil {
			l.front.Prev = nil
		} else {
			l.back = nil
		}
		l.len--
		return true
	}
	return false
}

func (l *list) push2Empty(li *ListItem) bool {
	if l.len == 0 {
		l.back, l.front = li, li
		l.len++
		return true
	}
	return false
}

func NewList() List {
	return &list{len: 0, front: nil, back: nil, mu: &sync.Mutex{}}
}
