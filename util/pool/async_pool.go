package pool

import (
	"context"

	pctx "github.com/appootb/substratum/v2/plugin/context"
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

func WithAsyncContext(ctx context.Context) AsyncOption {
	return func(pool *AsyncPool) {
		pool.ctx = ctx
	}
}

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

type AsyncPool struct {
	ctx  context.Context
	stop context.CancelFunc

	concurrency int
	chanLength  int
	component   string

	ch chan interface{}
}

func NewAsync(handler AsyncHandler, opts ...AsyncOption) *AsyncPool {
	pool := &AsyncPool{
		ctx:         context.Background(),
		concurrency: DefaultAsyncConcurrency,
		chanLength:  DefaultAsyncChanLength,
	}
	for _, o := range opts {
		o(pool)
	}
	pool.ch = make(chan interface{}, pool.chanLength)
	pool.ctx, pool.stop = context.WithCancel(pool.ctx)
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

func (pool *AsyncPool) Stop() {
	pool.stop()
}

func (pool *AsyncPool) run(h AsyncHandler) {
	for {
		select {
		case d := <-pool.ch:
			ctx := pctx.WithImplementContext(pool.ctx, pool.component)
			h.Handle(ctx, d)
		case <-pool.ctx.Done():
			return
		}
	}
}
