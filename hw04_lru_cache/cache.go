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

type cacheKeyValue struct {
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

func (c *lruCache) Set(key Key, value interface{}) bool {
	if cachedNode, ok := c.items[key]; ok {
		cachedNode.Value = cacheKeyValue{
			key:   key,
			value: value,
		}
		c.queue.MoveToFront(cachedNode)
		return true
	}

	if c.capacity == c.queue.Len() {
		lastNode := c.queue.Back()
		delete(c.items, lastNode.Value.(cacheKeyValue).key)
		c.queue.Remove(lastNode)
	}

	newNode := c.queue.PushFront(cacheKeyValue{
		key:   key,
		value: value,
	})
	c.items[key] = newNode

	return false
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	if val, ok := c.items[key]; ok {
		c.queue.MoveToFront(c.items[key])
		if val.Value != nil {
			return val.Value.(cacheKeyValue).value, true
		}

		return nil, false
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
