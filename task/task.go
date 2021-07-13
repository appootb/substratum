package task

import "context"

var (
	impl Task
)

// Implementor returns the task service implementor.
func Implementor() Task {
	return impl
}

// RegisterImplementor registers the task service implementor.
func RegisterImplementor(c Task) {
	impl = c
}

type Executor interface {
	Execute(context.Context, interface{}) error
}

type ExecutorFunc func(ctx context.Context, arg interface{}) error

func (fn ExecutorFunc) Execute(ctx context.Context, arg interface{}) error {
	return fn(ctx, arg)
}

type Task interface {
	// Schedule a task.
	// Supported spec, refer: https://github.com/robfig/cron/tree/v3.0.1
	Schedule(spec string, exec Executor, opts ...Option) error
}
