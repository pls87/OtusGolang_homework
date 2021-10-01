package hw04lrucache

type Key string

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

func (lc *lruCache) Set(key Key, value interface{}) bool {
	if lc.items[key] == nil {
		lc.items[key] = lc.queue.PushFront(&cacheItem{key: key, value: value})
		if lc.queue.Len() > lc.capacity {
			back := lc.queue.Back()
			lc.queue.Remove(back)
			delete(lc.items, back.Value.(*cacheItem).key)
		}
		return false
	} else {
		lc.items[key].Value.(*cacheItem).value = value
		lc.queue.MoveToFront(lc.items[key])
		return true
	}
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	if lc.items[key] != nil {
		lc.queue.MoveToFront(lc.items[key])
		return lc.items[key].Value.(*cacheItem).value, true
	}
	return nil, false
}

func (lc *lruCache) Clear() {
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
