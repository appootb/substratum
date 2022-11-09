package task

type Option func(options *Options)

type Options struct {
	Component string
	Name      string
	Singleton bool
	Argument  interface{}
}

var EmptyOptions = func() *Options {
	return &Options{
		Singleton: false,
		Argument:  nil,
	}
}

func WithComponent(component string) Option {
	return func(opts *Options) {
		opts.Component = component
	}
}

func WithName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
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
