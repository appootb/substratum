package discovery

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/appootb/substratum/v2/discovery"
	ictx "github.com/appootb/substratum/v2/internal/context"
	"github.com/appootb/substratum/v2/logger"
	builder "github.com/appootb/substratum/v2/resolver"
	"google.golang.org/grpc/resolver"
)

func Init() {
	if discovery.BackendImplementor() == nil {
		discovery.RegisterBackendImplementor(newDebug())
	}
	if discovery.Implementor() == nil {
		discovery.RegisterImplementor(newDiscovery())
	}
}

const (
	ServicePrefix    = "service"
	ServiceNodeIDKey = "node_id"
)

type NodeInfo struct {
	TTL      time.Duration     `json:"ttl,omitempty"`
	UniqueID int64             `json:"unique_id,omitempty"`
	Service  string            `json:"service,omitempty"`
	Address  string            `json:"address,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (info *NodeInfo) String() string {
	v, _ := json.Marshal(info)
	return string(v)
}

func newDiscovery() discovery.Discovery {
	impl := &Discovery{}
	if err := impl.refresh(); err != nil {
		logger.Fatal("discovery initialize failed", logger.Content{
			"error": err.Error(),
		})
	}
	return impl
}

type Discovery struct {
	lc sync.Map // Local registered service
	rc sync.Map // Remote registered service
}

// PassthroughAddr returns service unique ID and rpc address registered for the component.
func (m *Discovery) PassthroughAddr(component string) (int64, string) {
	if addr, ok := m.lc.Load(component); ok {
		info := addr.(*NodeInfo)
		return info.UniqueID, info.Address
	}
	return 0, ""
}

// Register rpc address of the component node.
func (m *Discovery) Register(component, addr string, opts ...discovery.Option) (int64, error) {
	options := discovery.EmptyOptions()
	for _, o := range opts {
		o(options)
	}
	//
	idKey := fmt.Sprintf("%s/%s", ServiceNodeIDKey, component)
	uniqueID, err := discovery.BackendImplementor().Incr(idKey)
	if err != nil {
		return 0, err
	}
	if options.Isolate {
		return uniqueID, nil
	}
	//
	info := &NodeInfo{
		TTL:      options.TTL,
		UniqueID: uniqueID,
		Service:  component,
		Address:  addr,
		Metadata: map[string]string{
			ServicePrefix: strings.Join(options.Services, ","),
		},
	}
	nodeKey := fmt.Sprintf("%s/%s/%s", ServicePrefix, component, addr)
	if err = discovery.BackendImplementor().KeepAlive(nodeKey, info.String(), options.TTL); err != nil {
		return 0, err
	}
	//
	m.lc.Store(component, info)
	return info.UniqueID, nil
}

// GetAddresses returns rpc service addresses.
func (m *Discovery) GetAddresses(service string) []resolver.Address {
	parts := strings.Split(service, ".")
	//
	var nodes map[string]*NodeInfo
	if cache, ok := m.rc.Load(parts[0]); ok {
		nodes = cache.(map[string]*NodeInfo)
	}
	if nodes == nil || len(nodes) == 0 {
		return []resolver.Address{}
	}
	//
	addresses := make([]resolver.Address, 0, len(nodes))
	if len(parts) > 1 {
		for _, info := range nodes {
			if info.Metadata == nil || len(info.Metadata) == 0 {
				continue
			}
			services := strings.Split(info.Metadata[ServicePrefix], ",")
			for _, name := range services {
				if service != name {
					continue
				}
				addresses = append(addresses, resolver.Address{
					Addr: info.Address,
				})
			}
		}
		if len(addresses) > 0 {
			return addresses
		}
	}
	// Get component by default
	for _, info := range nodes {
		addresses = append(addresses, resolver.Address{
			Addr: info.Address,
		})
	}
	return addresses
}

func (m *Discovery) refresh() error {
	// Get services.
	path := ServicePrefix + "/"
	version, err := m.getServices(path)
	if err != nil {
		return err
	}
	//
	evtChan, err := discovery.BackendImplementor().Watch(path, version, true)
	if err != nil {
		return err
	}
	//
	go m.watchEvent(path, evtChan)
	return nil
}

func (m *Discovery) getServices(path string) (uint64, error) {
	services := make(map[string]map[string]*NodeInfo)
	pairs, err := discovery.BackendImplementor().Get(path, true)
	if err != nil {
		return 0, err
	}
	for _, kv := range pairs.KVs {
		var n NodeInfo
		if err = json.Unmarshal([]byte(kv.Value), &n); err != nil {
			return 0, err
		}
		svc, ok := services[n.Service]
		if !ok {
			svc = make(map[string]*NodeInfo)
			services[n.Service] = svc
		}
		svc[n.Address] = &n
	}
	for name, svc := range services {
		if err = m.updateResolver(name, svc); err != nil {
			return 0, err
		}
	}
	return pairs.Version, nil
}

func (m *Discovery) updateService(path, service string) error {
	kvs, err := discovery.BackendImplementor().Get(path+service+"/", true)
	if err != nil {
		return err
	}
	//
	svc := make(map[string]*NodeInfo, len(kvs.KVs))
	for _, kv := range kvs.KVs {
		var n NodeInfo
		if err = json.Unmarshal([]byte(kv.Value), &n); err != nil {
			return err
		}
		svc[n.Address] = &n
	}
	//
	return m.updateResolver(service, svc)
}

func (m *Discovery) updateResolver(service string, nodes map[string]*NodeInfo) error {
	m.rc.Store(service, nodes)
	//
	addresses := make([]resolver.Address, 0, len(nodes))
	for _, info := range nodes {
		addresses = append(addresses, resolver.Address{
			Addr: info.Address,
		})
	}
	if len(addresses) > 0 && builder.Implementor() != nil {
		return builder.Implementor().UpdateAddresses(service, addresses)
	}
	return nil
}

func (m *Discovery) watchEvent(path string, ch discovery.EventChan) {
	var (
		err error
	)

	for {
		select {
		case <-ictx.Context.Done():
			discovery.BackendImplementor().Close()
			return

		case evt := <-ch:
			if evt.EventType == discovery.Refresh {
				_, err = m.getServices(path)
			} else {
				paths := strings.Split(strings.TrimPrefix(evt.Key, path), "/")
				err = m.updateService(path, paths[0])
			}
			if err != nil {
				logger.Error("substratum discovery event failed", logger.Content{
					"error": err.Error(),
					"event": evt.EventType,
					"key":   evt.Key,
				})
			}
		}
	}
}
