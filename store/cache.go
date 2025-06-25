package store

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, val []byte)
	InvalidateAll()
}

type cacheObject struct {
	bytes   []byte
	expires time.Time
}

type CacheStore struct {
	mu   sync.RWMutex
	data map[string]cacheObject
	ttl  time.Duration
}

func NewCache(ttl time.Duration) *CacheStore {
	return &CacheStore{
		data: make(map[string]cacheObject),
		ttl:  ttl,
	}
}

func (c *CacheStore) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(item.expires) {
		return nil, false
	}
	return item.bytes, true
}

func (c *CacheStore) Set(key string, val []byte) {
	c.mu.Lock()
	c.data[key] = cacheObject{bytes: val, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *CacheStore) InvalidateAll() {
	c.mu.Lock()
	c.data = make(map[string]cacheObject)
	c.mu.Unlock()
}
