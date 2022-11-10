package discovery

import "google.golang.org/grpc/resolver"

var (
	impl Discovery
)

// Implementor returns the discovery service implementor.
func Implementor() Discovery {
	return impl
}

// RegisterImplementor registers the discovery service implementor.
func RegisterImplementor(svc Discovery) {
	impl = svc
}

type Discovery interface {
	// PassthroughAddr returns service unique ID and rpc address registered for the component.
	PassthroughAddr(component string) (int64, string)

	// Register rpc address of the component service address.
	Register(component, addr string, opts ...Option) (int64, error)

	// GetAddresses returns rpc service addresses.
	GetAddresses(service string) []resolver.Address
}
