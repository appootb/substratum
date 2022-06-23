package substratum

import (
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/service"
)

type Service interface {
	service.Implementor

	// AddServeMux adds scoped ServeMux.
	AddServeMux(permission.VisibleScope, uint16, uint16) error

	// Register component.
	Register(Component, ...string) error

	// Serve start the mux server.
	Serve(isolate ...bool) error
}
