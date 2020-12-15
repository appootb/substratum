package cache

import "time"

type Option func(*base)

type ExpiredLoaderFunc func(interface{}) (interface{}, time.Duration, error)

func WithSize(size int) Option {
	return func(b *base) {
		b.size = size
	}
}

func WithExpiredLoader(fn ExpiredLoaderFunc) Option {
	return func(b *base) {
		b.expiredLoader = fn
	}
}
