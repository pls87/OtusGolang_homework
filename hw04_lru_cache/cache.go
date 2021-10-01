package hw04lrucache

import "sync"

type Key string

var cacheMutex sync.Mutex

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) (status bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if lc.items[key] == nil {
		lc.items[key] = lc.queue.PushFront(&cacheItem{key: key, value: value})
		if lc.queue.Len() > lc.capacity {
			back := lc.queue.Back()
			lc.queue.Remove(back)
			delete(lc.items, back.Value.(*cacheItem).key)
		}
		status = false
	} else {
		lc.items[key].Value.(*cacheItem).value = value
		lc.queue.MoveToFront(lc.items[key])
		status = true
	}

	return status
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if lc.items[key] == nil {
		return nil, false
	}

	lc.queue.MoveToFront(lc.items[key])
	return lc.items[key].Value.(*cacheItem).value, true
}

func (lc *lruCache) Clear() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
