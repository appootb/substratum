package cache

import "time"

// Option configures cache construction.
type Option func(*options)

type options struct {
	size int
}

// LoaderFunc loads a value for key. The returned duration is the cache TTL
// (DurationPersistence for no expiry; DurationMemoryLock / 0 to not store).
type LoaderFunc func(interface{}) (interface{}, time.Duration, error)

// WithSize sets the maximum number of entries. Non-positive values fall back to DefaultSize.
func WithSize(size int) Option {
	return func(o *options) {
		if size <= 0 {
			o.size = DefaultSize
			return
		}
		o.size = size
	}
}

func buildOptions(opts ...Option) options {
	o := options{size: DefaultSize}
	for _, opt := range opts {
		opt(&o)
	}
	if o.size <= 0 {
		o.size = DefaultSize
	}
	return o
}
