package gocache

import (
	"gocache/lru"
	"sync"
)

const (
	defaultReplicas = 3 // 一致性hash 算法默认复制次数
)

type cache struct {
	mu         sync.RWMutex
	lru        *lru.Cache
	cacheBytes int64
}

// add 往缓存中添加数据
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)

}

// get 获取缓存中的数据
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
