package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes  int64
	nbytes    int64
	cache     map[string]*list.Element
	ll        *list.List
	OnEvicted func(key string, value Value)
}

// New is Constructor of Cache
func New(maxBytes int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// Get
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest  缓存淘汰
func (c *Cache) RemoveOldest() {

	ele := c.ll.Back()

	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)

		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}

	}
}

// Add 新增或修改
func (c *Cache) Add(key string, val Value) {
	// 如果缓存已经存在
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)

		kv := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(kv.value.Len())
		kv.value = val
	} else {
		ele := c.ll.PushFront(&entry{value: val, key: key})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(val.Len())
	}

	// 内存淘汰
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 获取缓存条数
func (c *Cache) Len() int {
	return c.ll.Len()
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}
