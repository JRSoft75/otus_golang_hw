package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; found {
		// Update the value and move to front
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}

	// Add new item
	newItem := &cacheItem{key: key, value: value}
	listItem := c.queue.PushFront(newItem)
	c.items[key] = listItem

	// Check capacity
	if c.queue.Len() > c.capacity {
		// Remove the least recently used item
		backItem := c.queue.Back()
		if backItem != nil {
			c.queue.Remove(backItem)
			delete(c.items, backItem.Value.(*cacheItem).key)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; found {
		// Move to front and return value
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
