package configure

import (
	"strings"
	"sync"

	"github.com/appootb/substratum/configure"
	ictx "github.com/appootb/substratum/internal/context"
)

type node struct {
	k string
	v string
}

type watch struct {
	ch     configure.EventChan
	key    string
	prefix bool
}

type Debug struct {
	kvs map[string]*node
	ws  []*watch

	event configure.EventChan
	sync.RWMutex
}

func newDebug() configure.Backend {
	provider := &Debug{
		kvs:   make(map[string]*node),
		event: make(configure.EventChan, 10),
	}
	go provider.checkWatch()
	return provider
}

// Type returns the backend provider type.
func (m *Debug) Type() string {
	return "debug"
}

// Set value for the specified key.
func (m *Debug) Set(key, value string) error {
	m.Lock()
	m.kvs[key] = &node{
		k: key,
		v: value,
	}
	m.Unlock()
	m.event <- &configure.WatchEvent{
		EventType: configure.Update,
		KVPair: configure.KVPair{
			Key:   key,
			Value: value,
		},
	}
	return nil
}

// Get the value of the specified key or directory.
func (m *Debug) Get(key string, dir bool) (*configure.KVPairs, error) {
	m.RLock()
	defer m.RUnlock()
	if !dir {
		if n, ok := m.kvs[key]; !ok {
			return &configure.KVPairs{}, nil
		} else {
			return &configure.KVPairs{
				KVs: []*configure.KVPair{
					{
						Key:   key,
						Value: n.v,
					},
				},
			}, nil
		}
	}
	//
	var kvs []*configure.KVPair
	for k, v := range m.kvs {
		if strings.HasPrefix(k, key) {
			kvs = append(kvs, &configure.KVPair{
				Key:   k,
				Value: v.v,
			})
		}
	}
	return &configure.KVPairs{
		KVs: kvs,
	}, nil
}

// Watch for changes of the specified key or directory.
func (m *Debug) Watch(key string, _ uint64, dir bool) (configure.EventChan, error) {
	m.Lock()
	defer m.Unlock()
	ch := make(configure.EventChan, 10)
	m.ws = append(m.ws, &watch{
		ch:     ch,
		key:    key,
		prefix: dir,
	})
	return ch, nil
}

// Close the provider connection.
func (m *Debug) Close() {}

func (m *Debug) checkWatch() {
	for {
		select {
		case <-ictx.Context.Done():
			return

		case evt := <-m.event:
			m.RLock()
			for _, w := range m.ws {
				if evt.Key == w.key ||
					w.prefix && strings.HasPrefix(evt.Key, w.key) {
					w.ch <- evt
				}
			}
			m.RUnlock()
		}
	}
}
