package cache

import "time"

type Option func(*base)

type LoaderFunc func(interface{}) (interface{}, time.Duration, error)

func WithSize(size int) Option {
	return func(b *base) {
		b.size = size
	}
}

type op struct {
	loader LoaderFunc
}

func (opt *op) apply(opts []OpOption) {
	for _, o := range opts {
		o(opt)
	}
}

type OpOption func(*op)

func WithLoaderFunc(fn LoaderFunc) OpOption {
	return func(op *op) {
		op.loader = fn
	}
}
