package cron

import (
	"context"
)

type Option func(options *Options)

type Options struct {
	context.Context
	Name      string
	Singleton bool
	Argument  interface{}
}

var EmptyOptions = func() *Options {
	return &Options{
		Context:   context.Background(),
		Singleton: false,
		Argument:  nil,
	}
}

func WithName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func WithContext(ctx context.Context) Option {
	return func(opts *Options) {
		opts.Context = ctx
	}
}

func WithSingleton() Option {
	return func(opts *Options) {
		opts.Singleton = true
	}
}

func WithArgument(arg interface{}) Option {
	return func(opts *Options) {
		opts.Argument = arg
	}
}
