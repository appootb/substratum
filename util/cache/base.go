package cache

import (
	"container/list"
	"time"
)

// entry is used to hold a value in the evictList.
type entry struct {
	key    interface{}
	value  interface{}
	expire time.Time
}

func (e *entry) ExpiredAt(ts time.Time) bool {
	if e.expire.IsZero() {
		return false
	}
	if ts.IsZero() {
		ts = time.Now()
	}
	return e.expire.Before(ts)
}

type base struct {
	size      int
	items     map[interface{}]*list.Element
	evictList *list.List

	loaderLock syncLocker
}

func (c *base) set(key, value interface{}, expire time.Duration) {
	// Check for existing item
	if el, ok := c.items[key]; ok {
		c.evictList.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}

	// Add new item
	item := &entry{
		key:   key,
		value: value,
	}
	if expire > 0 {
		item.expire = time.Now().Add(expire)
	}
	c.items[key] = c.evictList.PushFront(item)

	// Verify size not exceeded
	if c.evictList.Len() > c.size {
		el := c.evictList.Back()
		if el != nil {
			c.removeElement(el)
		}
	}
}

func (c *base) get(key interface{}) (interface{}, bool) {
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

func (c *base) peek(key interface{}, withExpired ...bool) (interface{}, bool) {
	if el, ok := c.items[key]; ok {
		item := el.Value.(*entry)
		if (len(withExpired) > 0 && withExpired[0]) || !item.ExpiredAt(time.Now()) {
			return item.value, true
		}
	}
	return nil, false
}

func (c *base) delete(key interface{}) bool {
	if el, ok := c.items[key]; ok {
		c.removeElement(el)
		return true
	}
	return false
}

func (c *base) contain(key interface{}) bool {
	el, ok := c.items[key]
	if !ok {
		return false
	}
	return !el.Value.(*entry).ExpiredAt(time.Now())
}

func (c *base) length(withExpired ...bool) int {
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

func (c *base) keys(withExpired ...bool) []interface{} {
	now := time.Now()
	keys := make([]interface{}, 0, len(c.items))
	for el := c.evictList.Back(); el != nil; el = el.Prev() {
		item := el.Value.(*entry)
		if (len(withExpired) > 0 && withExpired[0]) || !item.ExpiredAt(now) {
			keys = append(keys, item.key)
		}
	}
	return keys
}

func (c *base) purge() {
	for k := range c.items {
		delete(c.items, k)
	}
	c.evictList.Init()
}

func (c *base) removeElement(el *list.Element) {
	c.evictList.Remove(el)
	item := el.Value.(*entry)
	delete(c.items, item.key)
}
