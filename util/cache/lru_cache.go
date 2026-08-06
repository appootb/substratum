package cache

import (
	"sync"
	"time"
)

// LRUCache is a thread-safe LRU cache with optional per-entry TTL.
type LRUCache struct {
	*lru
	mu sync.RWMutex
}

func newLRUCache(size int) *LRUCache {
	core := newLRU(size)
	c := &LRUCache{lru: core}
	core.loaderLock = syncLocker{
		cache: c,
		m:     make(map[interface{}]*caller),
	}
	return c
}

// Set key-value pair with an expiration (expire > 0).
// If expire equals 0, the key-value pair will be deleted.
func (c *LRUCache) Set(key, value interface{}, expire time.Duration) {
	c.mu.Lock()
	c.lru.set(key, value, expire)
	c.mu.Unlock()
}

// Get value from the cache by the key.
func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lru.get(key)
}

// GetOrLoad gets value from the cache or load by loader.
func (c *LRUCache) GetOrLoad(key interface{}, loader LoaderFunc) (interface{}, error) {
	c.mu.Lock()
	value, ok := c.lru.get(key)
	c.mu.Unlock()
	if ok {
		return value, nil
	}
	value, err := c.lru.loaderLock.Invoke(key, loader)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Peek returns value without updating the "recently used"-ness of the key.
func (c *LRUCache) Peek(key interface{}, withExpired ...bool) (interface{}, bool) {
	c.mu.Lock()
	value, ok := c.lru.peek(key, withExpired...)
	c.mu.Unlock()
	return value, ok
}

// Del the specified key from the cache.
func (c *LRUCache) Del(key interface{}) bool {
	c.mu.Lock()
	exist := c.lru.delete(key)
	c.mu.Unlock()
	return exist
}

// Contain checks if a key exists in the cache.
func (c *LRUCache) Contain(key interface{}) bool {
	c.mu.RLock()
	exist := c.lru.contain(key)
	c.mu.RUnlock()
	return exist
}

// Len returns the number of items in the cache.
func (c *LRUCache) Len(withExpired ...bool) int {
	c.mu.RLock()
	length := c.lru.length(withExpired...)
	c.mu.RUnlock()
	return length
}

// Keys returns a slice of the keys in the cache (oldest first).
func (c *LRUCache) Keys(withExpired ...bool) []interface{} {
	c.mu.RLock()
	keys := c.lru.keys(withExpired...)
	c.mu.RUnlock()
	return keys
}

// Purge clears the cache entities.
func (c *LRUCache) Purge() {
	c.mu.Lock()
	c.lru.purge()
	c.mu.Unlock()
}
