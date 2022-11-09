package resolver

import "google.golang.org/grpc/resolver"

var (
	impl Resolver
)

// Implementor return the gRPC resolver service implementor.
func Implementor() Resolver {
	return impl
}

// RegisterImplementor registers the gRPC resolver service implementor.
func RegisterImplementor(r Resolver) {
	impl = r
}

type Resolver interface {
	resolver.Builder

	UpdateAddresses(service string, addresses []string) error
	ReportAddressError(service string, err error)
}
