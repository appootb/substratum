package resolver

import (
	"google.golang.org/grpc/resolver"
)

var (
	impl resolver.Builder
)

// Return the service implementor.
func Implementor() resolver.Builder {
	return impl
}

// Register service implementor.
func RegisterImplementor(r resolver.Builder) {
	impl = r
}
