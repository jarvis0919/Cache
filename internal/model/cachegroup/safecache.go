package cachegroup

import (
	"sync"

	"cache/internal/model/lru"
)

// 线程安全的缓存，添加锁
type SafeCache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

//	func NewSafeCache(cacheBytes int64) *SafeCache {
//		return &SafeCache{
//			cacheBytes: cacheBytes,
//		}
//	}

// 添加缓存
func (c *SafeCache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// 获得缓存
func (c *SafeCache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
