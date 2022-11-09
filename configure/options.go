package configure

type Option func(options *Options)

type Options struct {
	AutoCreation bool
}

var EmptyOptions = func() *Options {
	return &Options{}
}

func WithAutoCreation(autoCreate bool) Option {
	return func(opts *Options) {
		opts.AutoCreation = autoCreate
	}
}
