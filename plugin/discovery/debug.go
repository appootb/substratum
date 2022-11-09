package discovery

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/appootb/substratum/v2/discovery"
	ictx "github.com/appootb/substratum/v2/internal/context"
)

var (
	zeroTime = time.Unix(0, 0)
)

type node struct {
	k      string
	v      string
	expire time.Time
}

type watch struct {
	ch     discovery.EventChan
	key    string
	prefix bool
}

type Debug struct {
	kvs map[string]*node
	ws  []*watch

	event discovery.EventChan
	sync.RWMutex
}

func newDebug() discovery.Backend {
	provider := &Debug{
		kvs:   make(map[string]*node),
		event: make(discovery.EventChan, 10),
	}
	go provider.checkTTL()
	go provider.checkWatch()
	return provider
}

// Type returns the backend provider type.
func (m *Debug) Type() string {
	return "debug"
}

// Set value for the specified key with a specified ttl.
func (m *Debug) Set(key, value string, ttl time.Duration) error {
	expire := time.Now().Add(ttl)
	if ttl == 0 {
		expire = zeroTime
	}
	m.Lock()
	m.kvs[key] = &node{
		k:      key,
		v:      value,
		expire: expire,
	}
	m.Unlock()
	m.event <- &discovery.WatchEvent{
		EventType: discovery.Update,
		KVPair: discovery.KVPair{
			Key:   key,
			Value: value,
		},
	}
	return nil
}

// Get the value of the specified key or directory.
func (m *Debug) Get(key string, dir bool) (*discovery.KVPairs, error) {
	m.RLock()
	defer m.RUnlock()
	if !dir {
		if n, ok := m.kvs[key]; !ok {
			return &discovery.KVPairs{}, nil
		} else {
			return &discovery.KVPairs{
				KVs: []*discovery.KVPair{
					{
						Key:   key,
						Value: n.v,
					},
				},
			}, nil
		}
	}
	//
	var kvs []*discovery.KVPair
	for k, v := range m.kvs {
		if strings.HasPrefix(k, key) {
			kvs = append(kvs, &discovery.KVPair{
				Key:   k,
				Value: v.v,
			})
		}
	}
	return &discovery.KVPairs{
		KVs: kvs,
	}, nil
}

// Incr invokes an atomic value increase for the specified key.
func (m *Debug) Incr(key string) (int64, error) {
	m.Lock()
	defer m.Unlock()
	n, ok := m.kvs[key]
	if !ok {
		n = &node{
			k:      key,
			v:      "0",
			expire: zeroTime,
		}
	}
	v, _ := strconv.ParseInt(n.v, 10, 64)
	v++
	n.v = strconv.FormatInt(v, 10)
	m.kvs[key] = n
	return v, nil
}

// Watch for changes of the specified key or directory.
func (m *Debug) Watch(key string, _ uint64, dir bool) (discovery.EventChan, error) {
	m.Lock()
	defer m.Unlock()
	ch := make(discovery.EventChan, 10)
	m.ws = append(m.ws, &watch{
		ch:     ch,
		key:    key,
		prefix: dir,
	})
	return ch, nil
}

// KeepAlive sets value and updates the ttl for the specified key.
func (m *Debug) KeepAlive(key, value string, ttl time.Duration) error {
	return m.Set(key, value, 0)
}

// Close the provider connection.
func (m *Debug) Close() {}

func (m *Debug) checkTTL() {
	ticker := time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case <-ictx.Context.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			m.Lock()
			for k, v := range m.kvs {
				if v.expire.Sub(zeroTime) > 0 && time.Now().Sub(v.expire) > 0 {
					delete(m.kvs, k)
					m.event <- &discovery.WatchEvent{
						EventType: discovery.Delete,
						KVPair: discovery.KVPair{
							Key:   v.k,
							Value: v.v,
						},
					}
				}
			}
			m.Unlock()
		}
	}
}

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
