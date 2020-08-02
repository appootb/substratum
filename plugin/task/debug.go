package task

import (
	"context"

	"github.com/appootb/substratum/task"
)

func Init() {
	if task.BackendImplementor() == nil {
		task.RegisterBackendImplementor(&Debug{})
	}
	if task.Implementor() == nil {
		task.RegisterImplementor(&Task{})
	}
}

type Debug struct{}

// Get the locker of the scheduler,
// should be blocked before acquired the locker.
func (m *Debug) Lock(_ string) error {
	return nil
}

// Keep alive the schedule locker.
func (m *Debug) KeepAlive(_ string) context.Context {
	return context.TODO()
}

// Give up the schedule locker.
func (m *Debug) Unlock(_ string) {
}
