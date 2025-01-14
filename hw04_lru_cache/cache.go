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

type subValue struct {
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
	var exist bool
	v := subValue{key: key, value: value}
	if item, ok := c.items[key]; ok {
		item.Value = v
		c.queue.MoveToFront(item)
		exist = true
	} else {
		newItem := c.queue.PushFront(v)
		c.items[key] = newItem
		if c.queue.Len() > c.capacity {
			backItem := c.queue.Back()
			c.queue.Remove(backItem)
			delete(c.items, backItem.Value.(subValue).key)
		}
	}

	return exist
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(subValue).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
