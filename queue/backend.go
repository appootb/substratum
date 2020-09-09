package queue

import (
	"time"
)

var (
	backendImpl Backend
)

// Return the service implementor.
func BackendImplementor() Backend {
	return backendImpl
}

// Register service implementor.
func RegisterBackendImplementor(backend Backend) {
	backendImpl = backend
}

// Message wrapper.
type MessageWrapper interface {
	Message
	MessageOperation
}

// Queue backend interface.
type Backend interface {
	// Backend type.
	Type() string
	// Ping connect the backend server if not connected.
	// Will be called before every Read/Write operation.
	Ping() error
	// Return the max delay duration supported by the backend.
	// A negative value means no limitation.
	// A zero value means delay operation is not supported.
	MaxDelay() time.Duration
	// Return all queue names in backend storage.
	GetQueues() ([]string, error)
	// Return all queue/topics in backend storage.
	GetTopics() (map[string][]string, error)
	// Return all topic length of specified queue in backend storage.
	GetQueueLength(queue string) (map[string]int64, error)
	// Return the specified queue/topic length in backend storage.
	GetTopicLength(queue, topic string) (int64, error)

	// Read subscribes the message of the specified queue and topic.
	Read(queue, topic string, ch chan<- MessageWrapper) error
	// Write publishes content data to the specified queue.
	Write(queue string, delay time.Duration, content []byte) error
}
