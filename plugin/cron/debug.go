package cron

import (
	"context"

	"github.com/appootb/substratum/cron"
)

func Init() {
	if cron.BackendImplementor() == nil {
		cron.RegisterBackendImplementor(&Debug{})
	}
	if cron.Implementor() == nil {
		cron.RegisterImplementor(&Cron{})
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
