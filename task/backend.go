package task

import "context"

var (
	backendImpl Backend
)

// Return the service implementor.
func BackendImplementor() Backend {
	return backendImpl
}

// Register service implementor.
func RegisterBackendImplementor(backend Backend) {
	backendImpl = backend
}

// Backend interface.
type Backend interface {
	// Get the locker of the scheduler,
	// should be blocked before acquired the locker.
	Lock(scheduler string) error
	// Keep alive the schedule locker.
	KeepAlive(scheduler string) context.Context
	// Give up the schedule locker.
	Unlock(scheduler string)
}
