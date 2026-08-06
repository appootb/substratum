package cache

import (
	"sync"
	"time"
)

// LFUCache is a thread-safe LFU cache with optional per-entry TTL.
// Get increments frequency; Set updates value/TTL without changing frequency.
type LFUCache struct {
	*lfu
	mu sync.RWMutex
}

func newLFUCache(size int) *LFUCache {
	core := newLFU(size)
	c := &LFUCache{lfu: core}
	core.loaderLock = syncLocker{
		cache: c,
		m:     make(map[interface{}]*caller),
	}
	return c
}

// Set stores key-value with expiration. expire == 0 deletes the key.
func (c *LFUCache) Set(key, value interface{}, expire time.Duration) {
	c.mu.Lock()
	c.lfu.set(key, value, expire)
	c.mu.Unlock()
}

// Get returns a non-expired value and increments its frequency.
func (c *LFUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lfu.get(key)
}

// GetOrLoad returns a cached value or loads via loader (single-flight).
func (c *LFUCache) GetOrLoad(key interface{}, loader LoaderFunc) (interface{}, error) {
	c.mu.Lock()
	value, ok := c.lfu.get(key)
	c.mu.Unlock()
	if ok {
		return value, nil
	}
	value, err := c.lfu.loaderLock.Invoke(key, loader)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Peek returns a value without changing frequency.
func (c *LFUCache) Peek(key interface{}, withExpired ...bool) (interface{}, bool) {
	c.mu.Lock()
	value, ok := c.lfu.peek(key, withExpired...)
	c.mu.Unlock()
	return value, ok
}

// Del removes the key.
func (c *LFUCache) Del(key interface{}) bool {
	c.mu.Lock()
	exist := c.lfu.delete(key)
	c.mu.Unlock()
	return exist
}

// Contain reports whether a non-expired entry exists.
func (c *LFUCache) Contain(key interface{}) bool {
	c.mu.RLock()
	exist := c.lfu.contain(key)
	c.mu.RUnlock()
	return exist
}

// Len returns the entry count.
func (c *LFUCache) Len(withExpired ...bool) int {
	c.mu.RLock()
	length := c.lfu.length(withExpired...)
	c.mu.RUnlock()
	return length
}

// Keys returns keys ordered by ascending frequency (then FIFO within a frequency).
func (c *LFUCache) Keys(withExpired ...bool) []interface{} {
	c.mu.RLock()
	keys := c.lfu.keys(withExpired...)
	c.mu.RUnlock()
	return keys
}

// Purge clears all entries.
func (c *LFUCache) Purge() {
	c.mu.Lock()
	c.lfu.purge()
	c.mu.Unlock()
}
