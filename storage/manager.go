package storage

var (
	impl Manager
)

// Implementor returns the storage manage service implementor.
func Implementor() Manager {
	return impl
}

// RegisterImplementor registers the storage manage service implementor.
func RegisterImplementor(mgr Manager) {
	if impl != nil {
		panic("storage: implementor already registered")
	}
	impl = mgr
}

type Manager interface {
	New(component string)
	Get(component string) Storage
}
