package discovery

import "time"

var (
	backendImpl Backend
)

// BackendImplementor returns the discovery backend service implementor.
func BackendImplementor() Backend {
	return backendImpl
}

// RegisterBackendImplementor registers the discovery backend service implementor.
func RegisterBackendImplementor(backend Backend) {
	backendImpl = backend
}

type EventType int

const (
	Update  EventType = iota + 1 // Key/value updated
	Delete                       // Key/value deleted
	Refresh                      // Refresh pah
)

type KVPair struct {
	Key     string
	Value   string
	Version uint64
}

type KVPairs struct {
	KVs     []*KVPair
	Version uint64
}

type WatchEvent struct {
	EventType
	KVPair
}

type EventChan chan *WatchEvent

// Backend interface.
type Backend interface {
	// Type returns the backend provider type.
	Type() string

	// Set value for the specified key with a specified ttl.
	Set(key, value string, ttl time.Duration) error

	// Get the value of the specified key or directory.
	Get(key string, dir bool) (*KVPairs, error)

	// Incr invokes an atomic value increase for the specified key.
	Incr(key string) (int64, error)

	// Watch for changes of the specified key or directory.
	Watch(key string, version uint64, dir bool) (EventChan, error)

	// KeepAlive sets value and updates the ttl for the specified key.
	KeepAlive(key, value string, ttl time.Duration) error

	// Close the provider connection.
	Close()
}
