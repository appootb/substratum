package cpool

import (
	"context"
)

const (
	DefaultConcurrency = 1
	DefaultChanLength  = 100
)

type ConsumeCallback func(interface{})

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

type ConsumePool struct {
	ctx  context.Context
	stop context.CancelFunc

	concurrency int
	chanLength  int

	fn ConsumeCallback
	ch chan interface{}
}

func New(ctx context.Context, callback ConsumeCallback, opts ...Option) *ConsumePool {
	cp := &ConsumePool{
		fn:          callback,
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

func (cp *ConsumePool) run(callback ConsumeCallback) {
	for {
		select {
		case d := <-cp.ch:
			callback(d)
		case <-cp.ctx.Done():
			return
		}
	}
}
