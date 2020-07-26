package queue

import (
	"github.com/appootb/substratum/queue"
)

type EmptyIdempotent struct{}

// Invoked before process message.
func (p *EmptyIdempotent) BeforeProcess(_ queue.Message) bool {
	return true
}

// Invoked after processing.
func (p *EmptyIdempotent) AfterProcess(_ queue.Message, _ queue.ProcessStatus) {
}
