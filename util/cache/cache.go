package cache

import (
	"container/list"
	"time"
)

type EvictType string

const (
	LRU EvictType = "lru"
	ARC           = "arc"
)

const (
	DefaultSize = 100
)

// Cache is the interface for LRU/ARC cache.
type Cache interface {
	// Set key-value pair with an expiration.
	Set(key, value interface{}, expire time.Duration)

	// Get value from the cache by the key.
	Get(key interface{}) (interface{}, bool)

	// Return value without updating the "recently used"-ness of the key.
	Peek(key interface{}) (interface{}, bool)

	// Delete the specified key from the cache.
	Del(key interface{}) bool

	// Check if a key exists in the cache.
	Contain(key interface{}) bool

	// Return the number of items in the cache.
	Len(withoutExpired ...bool) int

	// Return a slice of the keys in the cache.
	Keys(withoutExpired ...bool) []interface{}

	// Clear the cache entities.
	Purge()
}

func New(evictType EvictType, opts ...Option) (c Cache) {
	//
	b := base{
		size:      DefaultSize,
		evictList: list.New(),
	}
	for _, o := range opts {
		o(&b)
	}
	b.items = make(map[interface{}]*list.Element, b.size+1)
	//
	switch evictType {
	case LRU:
		c = &LRUCache{
			base: &b,
		}
	default:
		return nil
	}
	//
	b.syncLock = syncLocker{
		cache: c,
		m:     make(map[interface{}]*caller),
	}
	return c
}
