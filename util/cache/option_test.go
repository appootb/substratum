package cache

import "testing"

func TestWithSizeNonPositiveFallsBack(t *testing.T) {
	for _, typ := range []EvictType{LRU, LFU} {
		for _, sz := range []int{0, -1} {
			c := New(typ, WithSize(sz))
			if c == nil {
				t.Fatalf("%s nil cache", typ)
			}
			var size int
			switch x := c.(type) {
			case *LRUCache:
				size = x.size
			case *LFUCache:
				size = x.size
			}
			if size != DefaultSize {
				t.Fatalf("%s WithSize(%d): size=%d want DefaultSize=%d", typ, sz, size, DefaultSize)
			}
			key := string(typ)
			c.Set(key, 1, DurationPersistence)
			if _, ok := c.Get(key); !ok {
				t.Fatalf("%s WithSize(%d): key should be retained", typ, sz)
			}
		}
	}
}
