package cache

import (
	"container/list"
	"sort"
	"time"
)

// lfuNode is an LFU entry with frequency and list linkage within its freq bucket.
type lfuNode struct {
	entry
	freq int
	el   *list.Element // element in freqLists[freq]
}

// lfu is an O(1) LFU store: key map + per-frequency doubly-linked lists.
// Within the same frequency, the list front is evicted first (FIFO among ties).
type lfu struct {
	size      int
	items     map[interface{}]*lfuNode
	freqLists map[int]*list.List
	minFreq   int

	loaderLock syncLocker
}

func newLFU(size int) *lfu {
	return &lfu{
		size:      size,
		items:     make(map[interface{}]*lfuNode, size+1),
		freqLists: make(map[int]*list.List),
	}
}

func (c *lfu) freqList(freq int) *list.List {
	lst, ok := c.freqLists[freq]
	if !ok {
		lst = list.New()
		c.freqLists[freq] = lst
	}
	return lst
}

func (c *lfu) set(key, value interface{}, expire time.Duration) {
	if expire == 0 {
		c.delete(key)
		return
	}

	if n, ok := c.items[key]; ok {
		n.value = value
		n.expire = expireAt(expire)
		// Set does not change frequency.
		return
	}

	n := &lfuNode{
		entry: entry{
			key:    key,
			value:  value,
			expire: expireAt(expire),
		},
		freq: 1,
	}
	n.el = c.freqList(1).PushBack(n)
	c.items[key] = n
	c.minFreq = 1
	c.enforceSize()
}

func (c *lfu) get(key interface{}) (interface{}, bool) {
	n, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if n.ExpiredAt(time.Now()) {
		c.removeNode(n)
		return nil, false
	}
	c.increment(n)
	return n.value, true
}

func (c *lfu) peek(key interface{}, withExpired ...bool) (interface{}, bool) {
	n, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if (len(withExpired) > 0 && withExpired[0]) || !n.ExpiredAt(time.Now()) {
		return n.value, true
	}
	return nil, false
}

func (c *lfu) delete(key interface{}) bool {
	n, ok := c.items[key]
	if !ok {
		return false
	}
	c.removeNode(n)
	return true
}

func (c *lfu) contain(key interface{}) bool {
	n, ok := c.items[key]
	if !ok {
		return false
	}
	return !n.ExpiredAt(time.Now())
}

func (c *lfu) length(withExpired ...bool) int {
	if len(withExpired) > 0 && withExpired[0] {
		return len(c.items)
	}
	now := time.Now()
	n := 0
	for _, node := range c.items {
		if !node.ExpiredAt(now) {
			n++
		}
	}
	return n
}

func (c *lfu) keys(withExpired ...bool) []interface{} {
	now := time.Now()
	freqs := make([]int, 0, len(c.freqLists))
	for f := range c.freqLists {
		freqs = append(freqs, f)
	}
	sort.Ints(freqs)

	keys := make([]interface{}, 0, len(c.items))
	// Low frequency first; within a freq, front → back (eviction order).
	for _, freq := range freqs {
		lst := c.freqLists[freq]
		for el := lst.Front(); el != nil; el = el.Next() {
			node := el.Value.(*lfuNode)
			if (len(withExpired) > 0 && withExpired[0]) || !node.ExpiredAt(now) {
				keys = append(keys, node.key)
			}
		}
	}
	return keys
}

func (c *lfu) purge() {
	c.items = make(map[interface{}]*lfuNode, c.size+1)
	c.freqLists = make(map[int]*list.List)
	c.minFreq = 0
}

func (c *lfu) increment(n *lfuNode) {
	oldFreq := n.freq
	oldList := c.freqLists[oldFreq]
	oldList.Remove(n.el)
	if oldList.Len() == 0 {
		delete(c.freqLists, oldFreq)
		if c.minFreq == oldFreq {
			c.minFreq = oldFreq + 1
		}
	}
	n.freq++
	n.el = c.freqList(n.freq).PushBack(n)
}

func (c *lfu) removeNode(n *lfuNode) {
	if lst, ok := c.freqLists[n.freq]; ok {
		lst.Remove(n.el)
		if lst.Len() == 0 {
			delete(c.freqLists, n.freq)
			if c.minFreq == n.freq {
				c.recomputeMinFreq()
			}
		}
	}
	delete(c.items, n.key)
}

func (c *lfu) recomputeMinFreq() {
	c.minFreq = 0
	for f := range c.freqLists {
		if c.minFreq == 0 || f < c.minFreq {
			c.minFreq = f
		}
	}
}

func (c *lfu) enforceSize() {
	if c.size <= 0 {
		c.purge()
		return
	}
	now := time.Now()
	for len(c.items) > c.size {
		if c.removeOneExpired(now) {
			continue
		}
		// Evict least-frequently-used; among ties, the oldest in that bucket (front).
		if c.minFreq == 0 {
			c.recomputeMinFreq()
		}
		lst := c.freqLists[c.minFreq]
		if lst == nil || lst.Len() == 0 {
			c.recomputeMinFreq()
			lst = c.freqLists[c.minFreq]
			if lst == nil || lst.Len() == 0 {
				return
			}
		}
		el := lst.Front()
		c.removeNode(el.Value.(*lfuNode))
	}
}

func (c *lfu) removeOneExpired(now time.Time) bool {
	for _, lst := range c.freqLists {
		for el := lst.Front(); el != nil; el = el.Next() {
			n := el.Value.(*lfuNode)
			if n.ExpiredAt(now) {
				c.removeNode(n)
				return true
			}
		}
	}
	return false
}
