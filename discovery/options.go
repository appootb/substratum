package discovery

import (
	"time"
)

type Option func(options *Options)

type Options struct {
	Isolate  bool
	Path     string
	TTL      time.Duration
	Services []string
}

var EmptyOptions = func() *Options {
	return &Options{
		TTL: time.Second * 3,
	}
}

func WithIsolate(isolate bool) Option {
	return func(opts *Options) {
		opts.Isolate = isolate
	}
}

func WithPath(path string) Option {
	return func(opts *Options) {
		opts.Path = path
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(opts *Options) {
		opts.TTL = ttl
	}
}

func WithServices(services []string) Option {
	return func(opts *Options) {
		opts.Services = services
	}
}
