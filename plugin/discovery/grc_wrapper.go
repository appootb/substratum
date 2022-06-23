package discovery

import (
	"strings"
	"sync"
	"time"

	"github.com/appootb/grc"
	"github.com/appootb/substratum/v2/discovery"
)

func Init() {
	if discovery.Implementor() == nil {
		debug, _ := grc.New(grc.WithDebugProvider())
		Register(debug)
	}
}

func Register(rc *grc.RemoteConfig) {
	discovery.RegisterImplementor(&GRCWrapper{
		rc: rc,
	})
}

type Node struct {
	UniqueID int64
	Address  string
}

type GRCWrapper struct {
	lc sync.Map
	rc *grc.RemoteConfig
}

// RegisteredNode returns service unique ID and rpc address registered for the component.
func (m *GRCWrapper) RegisteredNode(component string) (int64, string) {
	if addr, ok := m.lc.Load(component); ok {
		node := addr.(*Node)
		return node.UniqueID, node.Address
	}
	return 0, ""
}

// RegisterNode registers rpc address of the component node.
func (m *GRCWrapper) RegisterNode(component, rpcAddr string, rpcSvc []string, ttl time.Duration) (int64, error) {
	uniqueID, err := m.rc.RegisterNode(component, rpcAddr, grc.WithNodeTTL(ttl),
		grc.WithNodeMetadata(map[string]string{"services": strings.Join(rpcSvc, ",")}))
	if err != nil {
		return 0, err
	}
	m.lc.Store(component, &Node{
		UniqueID: uniqueID,
		Address:  rpcAddr,
	})
	return uniqueID, nil
}

// GetNodes returns rpc service nodes.
func (m *GRCWrapper) GetNodes(svc string) map[string]int {
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
