package cache

import (
	"sync"
	"time"
)

type caller struct {
	ch  chan bool
	val interface{}
	dur time.Duration
	err error
}

type syncLocker struct {
	cache Cache

	mu sync.Mutex
	m  map[interface{}]*caller
}

func (sg *syncLocker) Invoke(key interface{}, fn LoaderFunc) (interface{}, error) {
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
	c.val, c.dur, c.err = fn(key)
	sg.cache.Set(key, c.val, c.dur)
	close(c.ch)
	// remove from caller map
	sg.mu.Lock()
	delete(sg.m, key)
	sg.mu.Unlock()
	return c.val, c.err
}
