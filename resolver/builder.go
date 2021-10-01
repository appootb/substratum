package resolver

import (
	"google.golang.org/grpc/resolver"
)

var (
	impl resolver.Builder
)

// Implementor return the gRPC resolver service implementor.
func Implementor() resolver.Builder {
	return impl
}

// RegisterImplementor registers the gRPC resolver service implementor.
func RegisterImplementor(r resolver.Builder) {
	impl = r
}
