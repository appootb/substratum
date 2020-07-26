package discovery

var (
	configImpl Config
)

// Return the service implementor.
func ConfigImplementor() Config {
	return configImpl
}

// Register service implementor.
func RegisterConfigImplementor(cfg Config) {
	configImpl = cfg
}

type Config interface {
	// Register the configuration pointer.
	RegisterConfig(component string, v interface{}) error
}
