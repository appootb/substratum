package cache

import (
	"container/list"
	"testing"
	"time"
)

func newBase(size int) *base {
	return &base{
		size:      size,
		items:     make(map[interface{}]*list.Element, size+1),
		evictList: list.New(),
	}
}

func TestBase(t *testing.T) {
	c := newBase(128)

	for i := 0; i < 256; i++ {
		c.set(i, i, DurationPersistence)
	}
	if l := c.length(); l != 128 {
		t.Fatalf("bad length: %v", l)
	}

	for i, k := range c.keys() {
		if v, ok := c.get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := c.get(i)
		if ok {
			t.Fatal("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := c.get(i)
		if !ok {
			t.Fatal("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		ok := c.delete(i)
		if !ok {
			t.Fatal("should be contained")
		}
		ok = c.delete(i)
		if ok {
			t.Fatal("should not be contained")
		}
		_, ok = c.get(i)
		if ok {
			t.Fatal("should be deleted")
		}
	}

	c.get(192) // expect 192 to be last key in l.Keys()

	for i, k := range c.keys() {
		if (i < 63 && k != i+193) || (i == 63 && k != 192) {
			t.Fatalf("out of order key: %v", k)
		}
	}

	c.purge()
	if l := c.length(); l != 0 {
		t.Fatalf("bad len: %v", l)
	}
	if _, ok := c.get(200); ok {
		t.Fatal("should contain nothing")
	}
}

func TestBase_Contains(t *testing.T) {
	c := newBase(2)

	c.set(1, 1, DurationPersistence)
	c.set(2, 2, DurationPersistence)
	if !c.contain(1) {
		t.Fatal("1 should be contained")
	}

	c.set(3, 3, DurationPersistence)
	if c.contain(1) {
		t.Fatal("Contains should not have updated recent-ness of 1")
	}
}

func TestLRU_Peek(t *testing.T) {
	c := newBase(2)

	c.set(1, 1, DurationPersistence)
	c.set(2, 2, DurationPersistence)
	if v, ok := c.peek(1); !ok || v != 1 {
		t.Fatalf("1 should be set to 1: %v, %v", v, ok)
	}

	c.set(3, 3, DurationPersistence)
	if c.contain(1) {
		t.Fatal("should not have updated recent-ness of 1")
	}
}

func TestLRU_PeekExpired(t *testing.T) {
	c := newBase(2)

	c.set(1, 1, time.Nanosecond)
	if v, ok := c.peek(1, true); !ok || v != 1 {
		t.Fatalf("1 should be set to 1: %v, %v", v, ok)
	}

	if _, ok := c.peek(1, false); ok {
		t.Fatalf("1 should be expired")
	}
}
