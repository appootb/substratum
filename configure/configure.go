package configure

var (
	impl Configure
)

// Implementor returns the configuration service implementor.
func Implementor() Configure {
	return impl
}

// RegisterImplementor registers the configuration service implementor.
func RegisterImplementor(cfg Configure) {
	impl = cfg
}

type Configure interface {
	// Register the configuration pointer.
	Register(component string, v interface{}) error
}
