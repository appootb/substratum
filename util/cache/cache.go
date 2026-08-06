// Package cache provides process-local caches with TTL and pluggable eviction
// (LRU, LFU). GetOrLoad single-flights concurrent loads for the same key.
package cache

import "time"

// EvictType selects the eviction policy for New.
type EvictType string

const (
	// LRU is least-recently-used eviction.
	LRU EvictType = "lru"
	// LFU is least-frequently-used eviction.
	LFU EvictType = "lfu"
	// ARC is reserved (not implemented).
	ARC EvictType = "arc"
)

const (
	// DefaultSize is used when WithSize is omitted or non-positive.
	DefaultSize = 100
)

const (
	// DurationPersistence means the entry does not expire by time.
	DurationPersistence time.Duration = -1
	// DurationMemoryLock means do not store (GetOrLoad single-flight only / Set deletes).
	DurationMemoryLock time.Duration = 0
)

// Cache is the interface for process-local caches.
type Cache interface {
	// Set stores key-value with an expiration.
	// expire > 0 is a TTL; DurationPersistence keeps the entry until evicted;
	// expire == 0 deletes the key.
	Set(key, value interface{}, expire time.Duration)

	// Get returns the value and true if present and not expired.
	// Successful Get updates eviction metadata (recency or frequency).
	Get(key interface{}) (interface{}, bool)

	// GetOrLoad returns a cached value or loads it via loader (single-flight per key).
	GetOrLoad(key interface{}, loader LoaderFunc) (interface{}, error)

	// Peek returns the value without updating eviction metadata.
	// withExpired=true allows returning soft-expired entries still in memory.
	Peek(key interface{}, withExpired ...bool) (interface{}, bool)

	// Del removes the key. Reports whether it was present (including soft-expired).
	Del(key interface{}) bool

	// Contain reports whether a non-expired entry exists (does not update metadata).
	Contain(key interface{}) bool

	// Len returns the number of non-expired entries, or all internal entries when withExpired.
	Len(withExpired ...bool) int

	// Keys returns keys (order is policy-specific). withExpired includes soft-expired.
	Keys(withExpired ...bool) []interface{}

	// Purge removes all entries.
	Purge()
}

// New creates a cache of the given eviction type.
// Unsupported types return nil. Size defaults to DefaultSize.
func New(evictType EvictType, opts ...Option) Cache {
	o := buildOptions(opts...)
	switch evictType {
	case LRU:
		return newLRUCache(o.size)
	case LFU:
		return newLFUCache(o.size)
	default:
		return nil
	}
}
