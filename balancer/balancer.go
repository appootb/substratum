package balancer

import (
	"google.golang.org/grpc/balancer"
)

var (
	impl balancer.Builder
)

// Implementor returns the balancer builder service implementor.
func Implementor() balancer.Builder {
	return impl
}

// RegisterImplementor registers the balancer builder service implementor.
func RegisterImplementor(svc balancer.Builder) {
	impl = svc
}
