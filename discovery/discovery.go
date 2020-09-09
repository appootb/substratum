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
	// Return service unique ID and rpc address registered for the component.
	RegisteredNode(component string) (int64, string)

	// Register rpc address of the component node.
	RegisterNode(component, rpcAddr string, rpcSvc []string, ttl time.Duration) error

	// Get rpc service nodes.
	GetNodes(svc string) map[string]int
}
