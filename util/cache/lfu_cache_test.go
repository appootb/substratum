package cache

import (
	"sync"
	"testing"
	"time"
)

func TestLFUEvictsLeastFrequent(t *testing.T) {
	c := New(LFU, WithSize(2))
	c.Set("a", 1, DurationPersistence)
	c.Set("b", 2, DurationPersistence)
	// bump a twice
	c.Get("a")
	c.Get("a")
	// b stays at freq 1; inserting c should drop b
	c.Set("c", 3, DurationPersistence)

	if _, ok := c.Get("b"); ok {
		t.Fatal("least-frequent b should be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("frequent a should remain")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("new c should remain")
	}
}

func TestLFUPeekDoesNotIncreaseFreq(t *testing.T) {
	c := New(LFU, WithSize(2))
	c.Set("a", 1, DurationPersistence)
	c.Set("b", 2, DurationPersistence)
	// Peek a many times — must not protect a over b
	for i := 0; i < 10; i++ {
		c.Peek("a")
	}
	c.Set("c", 3, DurationPersistence)
	// a and b both freq 1; FIFO among ties: a was inserted first → a evicted
	if _, ok := c.Get("a"); ok {
		t.Fatal("a should be evicted (peek must not raise freq)")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("b should remain")
	}
}

func TestLFUGetOrLoadSingleFlight(t *testing.T) {
	count := 0
	fn := func(key interface{}) (interface{}, time.Duration, error) {
		count++
		return 1, DurationPersistence, nil
	}
	c := New(LFU, WithSize(2))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := c.GetOrLoad(1, fn)
			if err != nil || v.(int) != 1 {
				t.Error("load failed")
			}
		}()
	}
	wg.Wait()
	if count != 1 {
		t.Fatalf("loader called %d times, want 1", count)
	}
}

func TestLFUCapacity(t *testing.T) {
	c := New(LFU, WithSize(128))
	for i := 0; i < 256; i++ {
		c.Set(i, i, DurationPersistence)
	}
	if c.Len() != 128 {
		t.Fatalf("len=%d want 128", c.Len())
	}
	// first 128 inserts should be gone (all freq 1, FIFO within freq)
	for i := 0; i < 128; i++ {
		if _, ok := c.Get(i); ok {
			t.Fatalf("key %d should be evicted", i)
		}
	}
	for i := 128; i < 256; i++ {
		if _, ok := c.Peek(i); !ok {
			t.Fatalf("key %d should remain", i)
		}
	}
}

func TestNewUnsupportedReturnsNil(t *testing.T) {
	if New(ARC) != nil {
		t.Fatal("ARC should be unimplemented")
	}
}
