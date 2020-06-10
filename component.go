package substratum

import "github.com/appootb/protobuf/go/service"

// Service component.
type Component interface {
	// Return the component name.
	Name() string

	// Init storage.
	InitStorage() error

	// Init service.
	InitService(service.Authenticator, service.Implementor) error
}
