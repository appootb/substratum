package task

import "context"

var (
	lockerImpl Locker
)

// LockerImplementor return the task locker service implementor.
func LockerImplementor() Locker {
	return lockerImpl
}

// RegisterLockerImplementor registers the task locker service implementor.
func RegisterLockerImplementor(locker Locker) {
	lockerImpl = locker
}

// Locker interface.
type Locker interface {
	// Lock tries to get the locker of the scheduler,
	// should be blocked before acquired the locker.
	Lock(ctx context.Context, scheduler string) context.Context
	// Unlock gives up the schedule locker.
	Unlock(scheduler string)
}
