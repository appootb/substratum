package discovery

import (
	"time"
)

var (
	impl Discovery
)

// Return the service implementor.
func Implementor() Discovery {
	return impl
}

// Register service implementor.
func RegisterImplementor(svc Discovery) {
	impl = svc
}

type Discovery interface {
	// Return local rpc address registered for the component.
	RegisteredAddr(component string) string

	// Register rpc address of the component node.
	RegisterNode(component, rpcAddr string, ttl time.Duration) error

	// Get component nodes.
	GetNodes(component string) map[string]int
}
