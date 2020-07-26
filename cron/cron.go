package cron

var (
	impl Cron
)

// Return the service implementor.
func Implementor() Cron {
	return impl
}

// Register service implementor.
func RegisterImplementor(c Cron) {
	impl = c
}

type JobFunc func(arg interface{})

type Cron interface {
	// Schedule a task.
	// Supported spec, refer: https://github.com/robfig/cron/tree/v3.0.1
	Schedule(spec string, fn JobFunc, opts ...Option) error
}
