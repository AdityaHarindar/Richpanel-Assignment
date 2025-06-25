package store

import (
	"sync"
	"time"
)

// Cache is an interface for a KV cache layer
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

// NewCache accepts the intended cache TTL, returns a pointer to a new CacheStore
func NewCache(ttl time.Duration) *CacheStore {
	return &CacheStore{
		data: make(map[string]cacheObject),
		ttl:  ttl,
	}
}

// Get accepts a cache key and returns the cache response, true/false for cache HIT/MISS
func (c *CacheStore) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(item.expires) {
		return nil, false
	}
	return item.bytes, true
}

// Set accepts a cache key/value and writes it to a cache store
func (c *CacheStore) Set(key string, val []byte) {
	c.mu.Lock()
	c.data[key] = cacheObject{bytes: val, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

// InvalidateAll clears the entire cache
func (c *CacheStore) InvalidateAll() {
	c.mu.Lock()
	c.data = make(map[string]cacheObject)
	c.mu.Unlock()
}
