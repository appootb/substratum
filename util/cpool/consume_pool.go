package cpool

import (
	"context"

	appootb "github.com/appootb/substratum/plugin/context"
)

const (
	DefaultConcurrency = 1
	DefaultChanLength  = 100
)

type Callback interface {
	Handle(context.Context, interface{})
}

type CallbackFunc func(context.Context, interface{})

func (fn CallbackFunc) Handle(ctx context.Context, arg interface{}) {
	fn(ctx, arg)
}

type Option func(pool *ConsumePool)

func WithConcurrency(concurrency int) Option {
	return func(cp *ConsumePool) {
		cp.concurrency = concurrency
	}
}

func WithChanLength(chanLen int) Option {
	return func(cp *ConsumePool) {
		cp.chanLength = chanLen
	}
}

func WithComponent(component string) Option {
	return func(cp *ConsumePool) {
		cp.component = component
	}
}

type ConsumePool struct {
	ctx  context.Context
	stop context.CancelFunc

	concurrency int
	chanLength  int
	component   string

	ch chan interface{}
}

func New(ctx context.Context, callback Callback, opts ...Option) *ConsumePool {
	cp := &ConsumePool{
		concurrency: DefaultConcurrency,
		chanLength:  DefaultChanLength,
	}
	for _, o := range opts {
		o(cp)
	}
	cp.ch = make(chan interface{}, cp.chanLength)
	cp.ctx, cp.stop = context.WithCancel(ctx)
	for i := 0; i < cp.concurrency; i++ {
		go cp.run(callback)
	}
	return cp
}

func (cp *ConsumePool) Add(data interface{}) bool {
	select {
	case cp.ch <- data:
		return true
	default:
		return false
	}
}

func (cp *ConsumePool) Stop() {
	cp.stop()
}

func (cp *ConsumePool) run(h Callback) {
	ctx := appootb.WithImplementContext(cp.ctx, cp.component)

	for {
		select {
		case d := <-cp.ch:
			h.Handle(ctx, d)
		case <-cp.ctx.Done():
			return
		}
	}
}
