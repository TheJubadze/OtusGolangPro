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

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, ok := cache.items[key]
	if ok {
		ci := item.Value.(cacheItem)
		ci.value = value
		item.Value = ci
		cache.queue.MoveToFront(item)
		return true
	}
	if cache.queue.Len()+1 > cache.capacity {
		oldest := cache.queue.Back()
		delete(cache.items, oldest.Value.(cacheItem).key)
		cache.queue.Remove(oldest)
	}
	item = cache.queue.PushFront(cacheItem{key: key, value: value})
	cache.items[key] = item
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, ok := cache.items[key]
	if ok {
		cache.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func (cache *lruCache) String() string {
	return cache.queue.String()
}
