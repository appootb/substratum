package cache

import "time"

type Option func(*base)

type LoaderFunc func(interface{}) (interface{}, time.Duration, error)

func WithSize(size int) Option {
	return func(b *base) {
		b.size = size
	}
}
