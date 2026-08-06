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
	if impl != nil {
		panic("configure: implementor already registered")
	}
	impl = cfg
}

type Configure interface {
	// Register the configuration pointer.
	Register(component string, v interface{}, opts ...Option) error
}
