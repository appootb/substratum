package cache

import (
	"sync"
)

type caller struct {
	ch  chan bool
	val interface{}
	err error
}

type syncLocker struct {
	cache Cache

	mu sync.Mutex
	m  map[interface{}]*caller
}

func (sg *syncLocker) Invoke(key interface{}, fn func() (interface{}, error)) (interface{}, error) {
	sg.mu.Lock()
	if v, ok := sg.cache.Peek(key); ok {
		sg.mu.Unlock()
		return v, nil
	}
	if c, ok := sg.m[key]; ok {
		sg.mu.Unlock()
		<-c.ch
		return c.val, c.err
	}

	// add new caller
	c := &caller{
		ch: make(chan bool),
	}
	sg.m[key] = c
	sg.mu.Unlock()

	// do invoke, and close caller channel
	c.val, c.err = fn()
	close(c.ch)
	// remove from caller map
	sg.mu.Lock()
	delete(sg.m, key)
	sg.mu.Unlock()
	return c.val, c.err
}
