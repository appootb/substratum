package cache

import (
	"container/list"
	"time"
)

type EvictType string

const (
	LRU EvictType = "lru"
	LFU           = "lfu"
	ARC           = "arc"
)

const (
	DefaultSize = 100
)

const (
	DurationPersistence = -1
	DurationMemoryLock  = 0
)

// Cache is the interface for LRU/ARC cache.
type Cache interface {
	// Set key-value pair with an expiration.
	Set(key, value interface{}, expire time.Duration)

	// Get value from the cache by the key.
	Get(key interface{}) (interface{}, bool)

	// GetOrLoad gets value from the cache or load by loader.
	GetOrLoad(key interface{}, loader LoaderFunc) (interface{}, error)

	// Peek returns value without updating the "recently used"-ness of the key.
	Peek(key interface{}, withExpired ...bool) (interface{}, bool)

	// Del the specified key from the cache.
	Del(key interface{}) bool

	// Contain checks if a key exists in the cache.
	Contain(key interface{}) bool

	// Len returns the number of items in the cache.
	Len(withExpired ...bool) int

	// Keys returns a slice of the keys in the cache.
	Keys(withExpired ...bool) []interface{}

	// Purge clears the cache entities.
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
	b.loaderLock = syncLocker{
		cache: c,
		m:     make(map[interface{}]*caller),
	}
	return c
}
