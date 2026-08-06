package cache

import (
	"testing"
	"time"
)

func TestLRUSetRefreshesExpireOnUpdate(t *testing.T) {
	c := New(LRU, WithSize(10))
	c.Set("k", "v1", 40*time.Millisecond)
	time.Sleep(15 * time.Millisecond)
	c.Set("k", "v2", time.Hour)
	time.Sleep(35 * time.Millisecond)
	v, ok := c.Get("k")
	if !ok || v != "v2" {
		t.Fatalf("want refreshed entry, got ok=%v v=%v", ok, v)
	}
}

func TestLRUSetRevivesSoftExpiredEntry(t *testing.T) {
	c := New(LRU, WithSize(10))
	c.Set("k", "old", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	c.Set("k", "new", time.Hour)
	v, ok := c.Get("k")
	if !ok || v != "new" {
		t.Fatalf("want revived entry, got ok=%v v=%v", ok, v)
	}
}

func TestLRUSetPersistenceClearsExpire(t *testing.T) {
	c := New(LRU, WithSize(10))
	c.Set("k", "v", 20*time.Millisecond)
	c.Set("k", "v2", DurationPersistence)
	time.Sleep(30 * time.Millisecond)
	if _, ok := c.Get("k"); !ok {
		t.Fatal("persistence update should clear TTL")
	}
}

func TestLRUExpiredPreferEvictedOverLive(t *testing.T) {
	c := New(LRU, WithSize(2)).(*LRUCache)
	c.Set("a", 1, 10*time.Millisecond)
	c.Set("b", 2, DurationPersistence)
	time.Sleep(20 * time.Millisecond)
	// a is expired but still mapped; inserting c must drop a, not b
	c.Set("c", 3, DurationPersistence)
	if _, ok := c.Get("b"); !ok {
		t.Fatal("live key b should not be evicted while expired a occupies capacity")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("new key c should be present")
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("expired a should be gone")
	}
}

func TestLRUExpiredDoNotExceedCapacityAfterInsert(t *testing.T) {
	c := New(LRU, WithSize(4)).(*LRUCache)
	for i := 0; i < 4; i++ {
		c.Set(i, i, 5*time.Millisecond)
	}
	time.Sleep(15 * time.Millisecond)
	if c.Len() != 0 {
		t.Fatalf("Len=%d want 0", c.Len())
	}
	// Insert one live key: enforceSize reclaims as needed; internal list stays <= size.
	c.Set("live", 1, DurationPersistence)
	c.mu.Lock()
	n := c.evictList.Len()
	c.mu.Unlock()
	if n > 4 {
		t.Fatalf("internal len=%d exceeds size", n)
	}
	if _, ok := c.Get("live"); !ok {
		t.Fatal("live key missing")
	}
}

// TestLRUExpiredLazyRetention documents that soft-expired entries stay in the
// internal list until Get/enforceSize reclaims them (Len reports 0).
func TestLRUExpiredLazyRetention(t *testing.T) {
	const size = 4
	c := New(LRU, WithSize(size)).(*LRUCache)
	for i := 0; i < size; i++ {
		c.Set(i, i, 5*time.Millisecond)
	}
	time.Sleep(15 * time.Millisecond)

	if c.Len() != 0 {
		t.Fatalf("Len=%d want 0 (logical empty)", c.Len())
	}
	c.mu.Lock()
	internal := c.evictList.Len()
	mapLen := len(c.items)
	c.mu.Unlock()
	if internal != size {
		t.Fatalf("evictList.Len()=%d want %d (lazy retention before reclaim)", internal, size)
	}
	if mapLen != size {
		t.Fatalf("items map len=%d want %d", mapLen, size)
	}
}
