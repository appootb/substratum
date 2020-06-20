package substratum

import (
	"context"

	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/storage"
)

// Service component.
type Component interface {
	// Return the component name.
	Name() string

	// Init component.
	Init(ctx context.Context) error

	// Init storage.
	InitStorage(storage.Storage) error

	// Register service.
	RegisterService(service.Authenticator, service.Implementor) error
}
