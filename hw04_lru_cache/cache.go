package hw04lrucache

import (
	"fmt"
	"sync"
)

type Key string

type Cache interface {
	fmt.Stringer
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
	keys     map[*ListItem]Key
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, ok := cache.items[key]
	if ok {
		item.Value = value
		cache.queue.MoveToFront(item)
		return true
	}
	item = cache.queue.PushFront(value)
	cache.items[key] = item
	cache.keys[item] = key
	if cache.queue.Len() > cache.capacity {
		back := cache.queue.Back()
		delete(cache.items, cache.keys[back])
		delete(cache.keys, back)
		cache.queue.Remove(back)
	}
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.keys = make(map[*ListItem]Key, cache.capacity)
}

func (cache *lruCache) String() string {
	return cache.queue.String()
}
