package hw04lrucache

import (
	"sync"
)

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
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

// Set adds or updates cached value by string key and move the cashed item to the top if the cached item exists.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		cacheItem := item.Value.(*cacheItem)
		cacheItem.value = value

		c.queue.MoveToFront(item)
		return true
	}

	newCacheItem := cacheItem{key: key, value: value}
	newItem := c.queue.PushFront(&newCacheItem)
	c.items[key] = newItem
	if c.queue.Len() > c.capacity {
		backItem := c.queue.Back()
		delete(c.items, backItem.Value.(*cacheItem).key)
		c.queue.Remove(c.queue.Back())
	}
	return false
}

// Get gets cached value by string key and move this item to the top if cache item exists.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

// Clear removes all cached values and clear queue.
func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
