package discovery

type Config interface {
	// Register the configuration pointer.
	RegisterConfig(component string, v interface{}) error
}
