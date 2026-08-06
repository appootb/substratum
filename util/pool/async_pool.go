// Package pool provides bounded worker pools for async handling and keyed batching.
package pool

import (
	"context"
	"runtime/debug"

	sctx "github.com/appootb/substratum/v2/context"
	ictx "github.com/appootb/substratum/v2/internal/context"
)

const (
	DefaultAsyncConcurrency = 1
	DefaultAsyncChanLength  = 100
)

// AsyncHandler processes a single async task.
type AsyncHandler interface {
	Handle(context.Context, interface{})
}

// AsyncFunc adapts a function to AsyncHandler.
type AsyncFunc func(context.Context, interface{})

// Handle implements AsyncHandler.
func (fn AsyncFunc) Handle(ctx context.Context, arg interface{}) {
	fn(ctx, arg)
}

// AsyncOption configures AsyncPool construction.
type AsyncOption func(pool *AsyncPool)

// WithAsyncConcurrency sets the number of worker goroutines (minimum 1).
func WithAsyncConcurrency(concurrency int) AsyncOption {
	return func(pool *AsyncPool) {
		pool.concurrency = concurrency
	}
}

// WithAsyncChanLength sets the buffered task queue length (minimum 1).
func WithAsyncChanLength(chanLen int) AsyncOption {
	return func(pool *AsyncPool) {
		pool.chanLength = chanLen
	}
}

// WithAsyncComponent sets the component name used for ServerContext in handlers.
func WithAsyncComponent(component string) AsyncOption {
	return func(pool *AsyncPool) {
		pool.component = component
	}
}

// AsyncPool is a non-blocking task queue with a fixed worker set.
// Add returns false when the queue is full (tasks are not queued).
type AsyncPool struct {
	concurrency int
	chanLength  int
	component   string

	ch chan interface{}
}

// NewAsync starts workers and returns a pool. Invalid options fall back to defaults.
func NewAsync(handler AsyncHandler, opts ...AsyncOption) *AsyncPool {
	pool := &AsyncPool{
		concurrency: DefaultAsyncConcurrency,
		chanLength:  DefaultAsyncChanLength,
	}
	for _, o := range opts {
		o(pool)
	}
	if pool.concurrency <= 0 {
		pool.concurrency = DefaultAsyncConcurrency
	}
	if pool.chanLength <= 0 {
		pool.chanLength = DefaultAsyncChanLength
	}
	pool.ch = make(chan interface{}, pool.chanLength)
	for i := 0; i < pool.concurrency; i++ {
		go pool.run(handler)
	}
	return pool
}

// Add enqueues data without blocking. Returns false if the buffer is full.
func (pool *AsyncPool) Add(data interface{}) bool {
	select {
	case pool.ch <- data:
		return true
	default:
		return false
	}
}

func (pool *AsyncPool) run(h AsyncHandler) {
	for {
		select {
		case d := <-pool.ch:
			pool.safeHandle(h, d)
		case <-ictx.Context.Done():
			return
		}
	}
}

func (pool *AsyncPool) safeHandle(h AsyncHandler, d interface{}) {
	defer func() {
		if recover() != nil {
			debug.PrintStack()
		}
	}()
	h.Handle(sctx.ServerContext(pool.component), d)
}
