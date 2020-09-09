package discovery

import (
	"strings"
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

type Node struct {
	UniqueID int64
	Address  string
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
func (m *Debug) RegisteredNode(component string) (int64, string) {
	if addr, ok := m.lc.Load(component); ok {
		node := addr.(*Node)
		return node.UniqueID, node.Address
	}
	return 0, ""
}

// Register rpc address of the component node.
func (m *Debug) RegisterNode(component, rpcAddr string, rpcSvc []string, ttl time.Duration) error {
	uniqueID, err := m.rc.RegisterNode(component, rpcAddr, grc.WithNodeTTL(ttl),
		grc.WithNodeMetadata(map[string]string{"services": strings.Join(rpcSvc, ",")}))
	if err != nil {
		return err
	}
	m.lc.Store(component, &Node{
		UniqueID: uniqueID,
		Address:  rpcAddr,
	})
	return nil
}

// Get rpc service nodes.
func (m *Debug) GetNodes(svc string) map[string]int {
	parts := strings.Split(svc, ".")
	component := parts[0]
	nodes := m.rc.GetNodes(component)
	result := make(map[string]int, len(nodes))

	if len(parts) > 1 {
		for _, node := range nodes {
			if node.Metadata == nil || len(node.Metadata) == 0 {
				continue
			}
			services := strings.Split(node.Metadata["services"], ",")
			for _, name := range services {
				if svc == name {
					result[node.Address] = node.Weight
				}
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	// Get component by default
	for _, node := range nodes {
		result[node.Address] = node.Weight
	}
	return result
}
