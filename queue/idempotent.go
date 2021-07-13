package queue

var (
	idempotentImpl Idempotent
)

// IdempotentImplementor returns the idempotent service implementor.
func IdempotentImplementor() Idempotent {
	return idempotentImpl
}

// RegisterIdempotentImplementor registers the idempotent service implementor.
func RegisterIdempotentImplementor(idempotent Idempotent) {
	idempotentImpl = idempotent
}

// ProcessStatus type.
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
	// BeforeProcess should be invoked before process message.
	// Returns true to continue the message processing.
	// Returns false to invoke Cancel for the message.
	BeforeProcess(Message) bool
	// AfterProcess should be invoked after processing.
	AfterProcess(Message, ProcessStatus)
}
