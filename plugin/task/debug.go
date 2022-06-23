package task

import (
	"context"

	"github.com/appootb/substratum/v2/task"
)

func Init() {
	if task.LockerImplementor() == nil {
		task.RegisterLockerImplementor(&Debug{})
	}
	if task.Implementor() == nil {
		task.RegisterImplementor(&Task{})
	}
}

type Debug struct{}

// Lock tries to get the locker of the scheduler,
// should be blocked before acquired the locker.
func (m *Debug) Lock(ctx context.Context, _ string) context.Context {
	return ctx
}

// Unlock gives up the schedule locker.
func (m *Debug) Unlock(_ string) {
}
