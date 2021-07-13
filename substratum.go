package substratum

import (
	"net/http"

	"github.com/appootb/substratum/proto/go/permission"
	"github.com/appootb/substratum/service"
)

type Service interface {
	service.Implementor

	// Handle registers the scoped http handler.
	Handle(scope permission.VisibleScope, pattern string, handler http.Handler)

	// HandleFunc registers the scoped http handle function.
	HandleFunc(scope permission.VisibleScope, pattern string, handler http.HandlerFunc)

	// AddMux adds scoped ServeMux.
	AddMux(permission.VisibleScope, uint16, uint16) error

	// Register component.
	Register(Component, ...string) error

	// Serve start the mux server.
	Serve(isolate ...bool) error
}
