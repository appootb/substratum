package task

import (
	"context"
)

type Option func(options *Options)

type Options struct {
	context.Context
	Component string
	Product   string
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

func WithComponent(component string) Option {
	return func(opts *Options) {
		opts.Component = component
	}
}

func WithProduct(product string) Option {
	return func(opts *Options) {
		opts.Product = product
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
