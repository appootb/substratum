package discovery

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/appootb/substratum/v2/discovery"
	ictx "github.com/appootb/substratum/v2/internal/context"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/resolver"
)

func Init() {
	if discovery.BackendImplementor() == nil {
		discovery.RegisterBackendImplementor(newDebug())
	}
	if discovery.Implementor() == nil {
		discovery.RegisterImplementor(&Discovery{})
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

type Discovery struct {
	initialized int32

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
	if err := m.initialize(options.Path); err != nil {
		return 0, err
	}
	//
	idKey := fmt.Sprintf("%s/%s/%s", options.Path, ServiceNodeIDKey, component)
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
	nodeKey := fmt.Sprintf("%s/%s/%s/%s", options.Path, ServicePrefix, component, addr)
	if err = discovery.BackendImplementor().KeepAlive(nodeKey, info.String(), options.TTL); err != nil {
		return 0, err
	}
	//
	m.lc.Store(component, info)
	return info.UniqueID, nil
}

// GetAddresses returns rpc service addresses.
func (m *Discovery) GetAddresses(service string) []string {
	parts := strings.Split(service, ".")
	//
	var nodes map[string]*NodeInfo
	if cache, ok := m.rc.Load(parts[0]); ok {
		nodes = cache.(map[string]*NodeInfo)
	}
	if nodes == nil || len(nodes) == 0 {
		return []string{}
	}
	//
	addrs := make([]string, 0, len(nodes))
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
				addrs = append(addrs, info.Address)
			}
		}
		if len(addrs) > 0 {
			return addrs
		}
	}
	// Get component by default
	for _, info := range nodes {
		addrs = append(addrs, info.Address)
	}
	return addrs
}

func (m *Discovery) initialize(path string) error {
	if atomic.AddInt32(&m.initialized, 1) != 1 {
		return nil
	}
	//
	// Get services.
	path = fmt.Sprintf("%s/%s/", path, ServicePrefix)
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
		m.rc.Store(name, svc)
	}
	return pairs.Version, nil
}

func (m *Discovery) updateService(path, service string) error {
	kvs, err := discovery.BackendImplementor().Get(path+service+"/", true)
	if err != nil {
		return err
	}
	//
	addrs := make([]string, 0, len(kvs.KVs))
	svc := make(map[string]*NodeInfo, len(kvs.KVs))
	for _, kv := range kvs.KVs {
		var n NodeInfo
		if err = json.Unmarshal([]byte(kv.Value), &n); err != nil {
			return err
		}
		svc[n.Address] = &n
		addrs = append(addrs, n.Address)
	}
	//
	m.rc.Store(service, svc)
	//
	return resolver.Implementor().UpdateAddresses(service, addrs)
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
