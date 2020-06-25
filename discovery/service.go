package discovery

import (
	"time"
)

type Service interface {
	// Register rpc address of the component node.
	RegisterNode(component, rpcAddr string, ttl time.Duration) error

	// Get component nodes.
	GetNodes(component string) map[string]int
}
