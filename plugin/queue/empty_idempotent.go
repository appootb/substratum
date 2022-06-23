package queue

import (
	"github.com/appootb/substratum/v2/queue"
)

type EmptyIdempotent struct{}

// BeforeProcess is invoked before process message.
func (p *EmptyIdempotent) BeforeProcess(_ queue.Message) bool {
	return true
}

// AfterProcess is invoked after processing.
func (p *EmptyIdempotent) AfterProcess(_ queue.Message, _ queue.ProcessStatus) {
}
