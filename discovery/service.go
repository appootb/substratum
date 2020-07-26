package discovery

import (
	"time"
)

var (
	serviceImpl Service
)

// Return the service implementor.
func Implementor() Service {
	return serviceImpl
}

// Register service implementor.
func RegisterImplementor(svc Service) {
	serviceImpl = svc
}

type Service interface {
	// Register rpc address of the component node.
	RegisterNode(component, rpcAddr string, ttl time.Duration) error

	// Get component nodes.
	GetNodes(component string) map[string]int
}
