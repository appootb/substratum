package task

import "context"

var (
	lockerImpl Locker
)

// Return the service implementor.
func LockerImplementor() Locker {
	return lockerImpl
}

// Register service implementor.
func RegisterLockerImplementor(locker Locker) {
	lockerImpl = locker
}

// Locker interface.
type Locker interface {
	// Get the locker of the scheduler,
	// should be blocked before acquired the locker.
	Lock(ctx context.Context, scheduler string) context.Context
	// Give up the schedule locker.
	Unlock(scheduler string)
}
