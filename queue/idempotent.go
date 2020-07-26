package queue

var (
	idempotentImpl Idempotent
)

// Return the service implementor.
func IdempotentImplementor() Idempotent {
	return idempotentImpl
}

// Register service implementor.
func RegisterIdempotentImplementor(idempotent Idempotent) {
	idempotentImpl = idempotent
}

// Status of the message.
type ProcessStatus int

const (
	Created ProcessStatus = iota
	Processing
	Canceled
	Succeeded
	Failed
	Requeued
)

// Idempotent interface.
type Idempotent interface {
	// Invoked before process message.
	// Returns true to continue the message processing.
	// Returns false to invoke Cancel for the message.
	BeforeProcess(Message) bool
	// Invoked after processing.
	AfterProcess(Message, ProcessStatus)
}
