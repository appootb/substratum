package queue

import (
	"context"
	"time"
)

var (
	backendImpl Backend
)

// BackendImplementor returns the queue backend service implementor.
func BackendImplementor() Backend {
	return backendImpl
}

// RegisterBackendImplementor registers the queue backend service implementor.
func RegisterBackendImplementor(backend Backend) {
	backendImpl = backend
}

// MessageWrapper interface.
type MessageWrapper interface {
	Message
	MessageOperation
}

// Backend interface.
type Backend interface {
	// Type returns backend type.
	Type() string
	// Ping connects the backend server if not connected.
	// Will be called before every Read/Write operation.
	Ping() error
	// MaxDelay returns the max delay duration supported by the backend.
	// A negative value means no limitation.
	// A zero value means delay operation is not supported.
	MaxDelay() time.Duration
	// GetQueues returns all queue names in backend storage.
	GetQueues() ([]string, error)
	// GetTopics returns all queue/topics in backend storage.
	GetTopics() (map[string][]string, error)
	// GetQueueLength returns all topic length of specified queue in backend storage.
	GetQueueLength(queue string) (map[string]int64, error)
	// GetTopicLength returns the specified queue/topic length in backend storage.
	GetTopicLength(queue, topic string) (int64, error)

	// Read subscribes the message of the specified queue and topic.
	Read(ctx context.Context, queue, topic string, ch chan<- MessageWrapper) error
	// Write publishes content data to the specified queue.
	Write(ctx context.Context, queue string, delay time.Duration, content []byte) error
}
