package cache

import (
	"sync"
	"time"
)

type LRUCache struct {
	*base
	mu sync.RWMutex
}

// Set key-value pair with an expiration.
func (c *LRUCache) Set(key, value interface{}, expire time.Duration) {
	c.mu.Lock()
	c.base.set(key, value, expire)
	c.mu.Unlock()
}

// Get value from the cache by the key.
func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.base.get(key)
}

func (c *LRUCache) GetOrLoad(key interface{}, loader LoaderFunc) (interface{}, error) {
	c.mu.Lock()
	value, ok := c.base.get(key)
	c.mu.Unlock()
	if ok {
		return value, nil
	}
	// try load
	value, err := c.base.loaderLock.Invoke(key, func() (interface{}, error) {
		v, dur, e := loader(key)
		if e != nil {
			return nil, e
		}
		c.mu.Lock()
		c.base.set(key, v, dur)
		c.mu.Unlock()
		return v, nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Return value without updating the "recently used"-ness of the key.
func (c *LRUCache) Peek(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	value, ok := c.base.peek(key)
	c.mu.Unlock()
	return value, ok
}

// Delete the specified key from the cache.
func (c *LRUCache) Del(key interface{}) bool {
	c.mu.Lock()
	exist := c.base.delete(key)
	c.mu.Unlock()
	return exist
}

// Check if a key exists in the cache.
func (c *LRUCache) Contain(key interface{}) bool {
	c.mu.RLock()
	exist := c.base.contain(key)
	c.mu.RUnlock()
	return exist
}

// Return the number of items in the cache.
func (c *LRUCache) Len(withoutExpired ...bool) int {
	c.mu.RLock()
	length := c.base.length(withoutExpired...)
	c.mu.RUnlock()
	return length
}

// Return a slice of the keys in the cache.
func (c *LRUCache) Keys(withoutExpired ...bool) []interface{} {
	c.mu.RLock()
	keys := c.base.keys(withoutExpired...)
	c.mu.RUnlock()
	return keys
}

// Clear the cache entities.
func (c *LRUCache) Purge() {
	c.mu.Lock()
	c.base.purge()
	c.mu.Unlock()
}
