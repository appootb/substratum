package discovery

import (
	"sync"
	"time"

	"github.com/appootb/grc"
	"github.com/appootb/substratum/discovery"
)

func Init() {
	if discovery.Implementor() == nil {
		discovery.RegisterImplementor(NewDebug())
	}
}

type Debug struct {
	lc sync.Map
	rc *grc.RemoteConfig
}

func NewDebug() *Debug {
	debug := &Debug{}
	debug.rc, _ = grc.New(grc.WithDebugProvider())
	return debug
}

// Return local rpc address registered for the component.
func (m *Debug) RegisteredAddr(component string) string {
	if addr, ok := m.lc.Load(component); ok {
		return addr.(string)
	}
	return ""
}

// Register rpc address of the component node.
func (m *Debug) RegisterNode(component, rpcAddr string, ttl time.Duration) error {
	err := m.rc.RegisterNode(component, grc.WithNodeAddress(rpcAddr), grc.WithNodeTTL(ttl))
	if err != nil {
		return err
	}
	m.lc.Store(component, rpcAddr)
	return nil
}

// Get component nodes.
func (m *Debug) GetNodes(component string) map[string]int {
	nodes := m.rc.GetNodes(component)
	result := make(map[string]int, len(nodes))
	for _, node := range nodes {
		result[node.Address] = node.Weight
	}
	return result
}
