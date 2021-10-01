package discovery

import (
	"time"
)

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
	// RegisteredNode returns service unique ID and rpc address registered for the component.
	RegisteredNode(component string) (int64, string)

	// RegisterNode registers rpc address of the component node.
	RegisterNode(component, rpcAddr string, rpcSvc []string, ttl time.Duration) (int64, error)

	// GetNodes returns rpc service nodes.
	GetNodes(svc string) map[string]int
}
