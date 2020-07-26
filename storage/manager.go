package storage

var (
	impl Manager
)

// Return the service implementor.
func Implementor() Manager {
	return impl
}

// Register service implementor.
func RegisterImplementor(mgr Manager) {
	impl = mgr
}

type Manager interface {
	New(component string)
	Get(component string) Storage
}
