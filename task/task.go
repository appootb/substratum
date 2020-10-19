package task

import "context"

var (
	impl Task
)

// Return the service implementor.
func Implementor() Task {
	return impl
}

// Register service implementor.
func RegisterImplementor(c Task) {
	impl = c
}

type JobFunc func(ctx context.Context, arg interface{}) error

type Task interface {
	// Schedule a task.
	// Supported spec, refer: https://github.com/robfig/cron/tree/v3.0.1
	Schedule(spec string, fn JobFunc, opts ...Option) error
}
