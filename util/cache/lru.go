package cache

import (
	"container/list"
	"time"
)

// lru is the LRU store: doubly-linked list (front = MRU) + key map.
type lru struct {
	size      int
	items     map[interface{}]*list.Element
	evictList *list.List

	loaderLock syncLocker
}

func newLRU(size int) *lru {
	return &lru{
		size:      size,
		items:     make(map[interface{}]*list.Element, size+1),
		evictList: list.New(),
	}
}

func (c *lru) set(key, value interface{}, expire time.Duration) {
	if expire == 0 {
		c.delete(key)
		return
	}

	// Existing item (including soft-expired still in the map).
	if el, ok := c.items[key]; ok {
		c.evictList.MoveToFront(el)
		item := el.Value.(*entry)
		item.value = value
		item.expire = expireAt(expire)
		return
	}

	item := &entry{
		key:    key,
		value:  value,
		expire: expireAt(expire),
	}
	c.items[key] = c.evictList.PushFront(item)
	c.enforceSize()
}

// enforceSize keeps at most c.size elements. Expired entries are removed first
// (from LRU end toward MRU); then the least-recently-used live entry is evicted.
func (c *lru) enforceSize() {
	if c.size <= 0 {
		c.purge()
		return
	}
	now := time.Now()
	for c.evictList.Len() > c.size {
		if c.removeOneExpired(now) {
			continue
		}
		el := c.evictList.Back()
		if el == nil {
			return
		}
		c.removeElement(el)
	}
}

func (c *lru) removeOneExpired(now time.Time) bool {
	for el := c.evictList.Back(); el != nil; el = el.Prev() {
		if el.Value.(*entry).ExpiredAt(now) {
			c.removeElement(el)
			return true
		}
	}
	return false
}

func (c *lru) get(key interface{}) (interface{}, bool) {
	if el, ok := c.items[key]; ok {
		item := el.Value.(*entry)
		if item.ExpiredAt(time.Now()) {
			c.removeElement(el)
			return nil, false
		}
		c.evictList.MoveToFront(el)
		return item.value, true
	}
	return nil, false
}

func (c *lru) peek(key interface{}, withExpired ...bool) (interface{}, bool) {
	if el, ok := c.items[key]; ok {
		item := el.Value.(*entry)
		if (len(withExpired) > 0 && withExpired[0]) || !item.ExpiredAt(time.Now()) {
			return item.value, true
		}
	}
	return nil, false
}

func (c *lru) delete(key interface{}) bool {
	if el, ok := c.items[key]; ok {
		c.removeElement(el)
		return true
	}
	return false
}

func (c *lru) contain(key interface{}) bool {
	el, ok := c.items[key]
	if !ok {
		return false
	}
	return !el.Value.(*entry).ExpiredAt(time.Now())
}

func (c *lru) length(withExpired ...bool) int {
	if len(withExpired) > 0 && withExpired[0] {
		return c.evictList.Len()
	}
	length := 0
	now := time.Now()
	for _, el := range c.items {
		if !el.Value.(*entry).ExpiredAt(now) {
			length++
		}
	}
	return length
}

func (c *lru) keys(withExpired ...bool) []interface{} {
	now := time.Now()
	keys := make([]interface{}, 0, len(c.items))
	// Oldest first (LRU end → MRU).
	for el := c.evictList.Back(); el != nil; el = el.Prev() {
		item := el.Value.(*entry)
		if (len(withExpired) > 0 && withExpired[0]) || !item.ExpiredAt(now) {
			keys = append(keys, item.key)
		}
	}
	return keys
}

func (c *lru) purge() {
	for k := range c.items {
		delete(c.items, k)
	}
	c.evictList.Init()
}

func (c *lru) removeElement(el *list.Element) {
	c.evictList.Remove(el)
	item := el.Value.(*entry)
	delete(c.items, item.key)
}
