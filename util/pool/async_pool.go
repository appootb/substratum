package pool

import (
	"context"

	sctx "github.com/appootb/substratum/context"
	ictx "github.com/appootb/substratum/internal/context"
)

const (
	DefaultAsyncConcurrency = 1
	DefaultAsyncChanLength  = 100
)

type AsyncHandler interface {
	Handle(context.Context, interface{})
}

type AsyncFunc func(context.Context, interface{})

func (fn AsyncFunc) Handle(ctx context.Context, arg interface{}) {
	fn(ctx, arg)
}

type AsyncOption func(pool *AsyncPool)

func WithAsyncConcurrency(concurrency int) AsyncOption {
	return func(pool *AsyncPool) {
		pool.concurrency = concurrency
	}
}

func WithAsyncChanLength(chanLen int) AsyncOption {
	return func(pool *AsyncPool) {
		pool.chanLength = chanLen
	}
}

func WithAsyncComponent(component string) AsyncOption {
	return func(pool *AsyncPool) {
		pool.component = component
	}
}

func WithAsyncProduct(product string) AsyncOption {
	return func(pool *AsyncPool) {
		pool.product = product
	}
}

type AsyncPool struct {
	concurrency int
	chanLength  int
	component   string
	product     string

	ch chan interface{}
}

func NewAsync(handler AsyncHandler, opts ...AsyncOption) *AsyncPool {
	pool := &AsyncPool{
		concurrency: DefaultAsyncConcurrency,
		chanLength:  DefaultAsyncChanLength,
	}
	for _, o := range opts {
		o(pool)
	}
	pool.ch = make(chan interface{}, pool.chanLength)
	for i := 0; i < pool.concurrency; i++ {
		go pool.run(handler)
	}
	return pool
}

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
			h.Handle(sctx.ServerContext(pool.component, pool.product), d)

		case <-ictx.Context.Done():
			return
		}
	}
}
