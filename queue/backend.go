package queue

import "github.com/appootb/substratum/v2/configure"

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

// InitBackend initializes the queue backend instance.
func InitBackend(addr configure.Address) error {
	return backendImpl.Init(addr)
}

// MessageWrapper interface.
type MessageWrapper interface {
	Message
	MessageOperation
}

// Backend interface.
type Backend interface {
	// Init queue backend instance.
	Init(cfg configure.Address) error

	// Type returns backend type.
	Type() string
	// Ping connects the backend server if not connected.
	// Will be called before every Read/Write operation.
	Ping() error

	// Read subscribes the message of the specified topic.
	Read(topic string, ch chan<- MessageWrapper, opts *SubscribeOptions) error
	// Write publishes content data to the specified queue.
	Write(topic string, content []byte, opts *PublishOptions) error
}
