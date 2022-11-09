package configure

type Option func(options *Options)

type Options struct {
	Path         string
	AutoCreation bool
}

var EmptyOptions = func() *Options {
	return &Options{}
}

func WithPath(path string) Option {
	return func(opts *Options) {
		opts.Path = path
	}
}

func WithAutoCreation(autoCreate bool) Option {
	return func(opts *Options) {
		opts.AutoCreation = autoCreate
	}
}
