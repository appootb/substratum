package configure

var (
	impl Configure
)

// Return the service implementor.
func Implementor() Configure {
	return impl
}

// Register service implementor.
func RegisterImplementor(cfg Configure) {
	impl = cfg
}

type Configure interface {
	// Register the configuration pointer.
	Register(component string, v interface{}) error
}
