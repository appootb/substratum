package cache

import (
	"testing"
	"time"
)

func TestLFUSetRefreshesExpireOnUpdate(t *testing.T) {
	c := New(LFU, WithSize(10))
	c.Set("k", "v1", 40*time.Millisecond)
	time.Sleep(15 * time.Millisecond)
	c.Set("k", "v2", time.Hour)
	time.Sleep(35 * time.Millisecond)
	v, ok := c.Get("k")
	if !ok || v != "v2" {
		t.Fatalf("want refreshed entry, got ok=%v v=%v", ok, v)
	}
}

func TestLFUSetRevivesSoftExpiredEntry(t *testing.T) {
	c := New(LFU, WithSize(10))
	c.Set("k", "old", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	c.Set("k", "new", time.Hour)
	v, ok := c.Get("k")
	if !ok || v != "new" {
		t.Fatalf("want revived entry, got ok=%v v=%v", ok, v)
	}
}

func TestLFUSetPersistenceClearsExpire(t *testing.T) {
	c := New(LFU, WithSize(10))
	c.Set("k", "v", 20*time.Millisecond)
	c.Set("k", "v2", DurationPersistence)
	time.Sleep(30 * time.Millisecond)
	if _, ok := c.Get("k"); !ok {
		t.Fatal("persistence update should clear TTL")
	}
}

func TestLFUTimeout(t *testing.T) {
	c := New(LFU, WithSize(2))
	c.Set(1, 1, DurationPersistence)
	c.Set(2, 2, 100*time.Millisecond)
	if !c.Contain(1) || !c.Contain(2) {
		t.Fatal("1 and 2 should be contained")
	}
	time.Sleep(200 * time.Millisecond)
	if _, ok := c.Get(2); ok {
		t.Fatal("2 should be expired")
	}
	if c.Len() != 1 {
		t.Fatalf("Len=%d want 1", c.Len())
	}
}

func TestLFUPeekExpired(t *testing.T) {
	c := New(LFU, WithSize(2))
	c.Set(1, 1, time.Nanosecond)
	time.Sleep(time.Millisecond)
	if v, ok := c.Peek(1, true); !ok || v != 1 {
		t.Fatalf("withExpired: want 1, got %v ok=%v", v, ok)
	}
	if _, ok := c.Peek(1, false); ok {
		t.Fatal("without withExpired: should not return expired entry")
	}
}

func TestLFUExpiredPreferEvictedOverLive(t *testing.T) {
	c := New(LFU, WithSize(2)).(*LFUCache)
	c.Set("a", 1, 10*time.Millisecond)
	c.Set("b", 2, DurationPersistence)
	c.Get("b") // raise b frequency
	time.Sleep(20 * time.Millisecond)
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

func TestLFUExpiredDoNotExceedCapacityAfterInsert(t *testing.T) {
	c := New(LFU, WithSize(4)).(*LFUCache)
	for i := 0; i < 4; i++ {
		c.Set(i, i, 5*time.Millisecond)
	}
	time.Sleep(15 * time.Millisecond)
	if c.Len() != 0 {
		t.Fatalf("Len=%d want 0", c.Len())
	}
	c.Set("live", 1, DurationPersistence)
	c.mu.Lock()
	n := len(c.items)
	c.mu.Unlock()
	if n > 4 {
		t.Fatalf("internal items=%d exceeds size", n)
	}
	if _, ok := c.Get("live"); !ok {
		t.Fatal("live key missing")
	}
}

// TestLFUExpiredLazyRetention documents that soft-expired entries stay until
// Get/enforceSize reclaims them (Len reports 0).
func TestLFUExpiredLazyRetention(t *testing.T) {
	const size = 4
	c := New(LFU, WithSize(size)).(*LFUCache)
	for i := 0; i < size; i++ {
		c.Set(i, i, 5*time.Millisecond)
	}
	time.Sleep(15 * time.Millisecond)

	if c.Len() != 0 {
		t.Fatalf("Len=%d want 0 (logical empty)", c.Len())
	}
	c.mu.Lock()
	mapLen := len(c.items)
	c.mu.Unlock()
	if mapLen != size {
		t.Fatalf("items map len=%d want %d (lazy retention before reclaim)", mapLen, size)
	}
}
