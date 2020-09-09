package substratum

import (
	"net/http"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/service"
)

type Service interface {
	service.Implementor

	// Scoped http handler.
	Handle(scope permission.VisibleScope, pattern string, handler http.Handler)

	// Scoped http handle function.
	HandleFunc(scope permission.VisibleScope, pattern string, handler http.HandlerFunc)

	// Add scoped ServeMux.
	AddMux(permission.VisibleScope, uint16, uint16) error

	// Register component.
	Register(Component, ...string) error

	// Serve start the mux server.
	Serve() error
}
