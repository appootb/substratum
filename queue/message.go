package queue

import (
	"context"
	"time"
)

// Queue message handler.
type MessageHandler func(context.Context, Message) error

// Queue message struct.
type Message interface {
	// Queue name of this message.
	Queue() string
	// Topic name of this message.
	Topic() string

	// Unique ID of this message.
	UniqueID() string
	// Message body content.
	Content() []byte
	// The creation time of the message.
	Timestamp() time.Time
	// The message should not be processed before this timestamp.
	NotBefore() time.Time
	// Message retry times.
	Retry() int
	// Return true for a ping message.
	IsPing() bool
}

// Queue message operation interface.
type MessageOperation interface {
	// Begin to process the message.
	Begin()
	// Indicate the message should be ignored.
	Cancel()
	// End indicates a successful process.
	End()
	// Requeue indicates the message should be retried.
	Requeue()
	// Fail indicates a failed process.
	Fail()
}

type PingMessage struct{}

// Queue name of this message.
func (m *PingMessage) Queue() string {
	return ""
}

// Topic name of this message.
func (m *PingMessage) Topic() string {
	return ""
}

// Unique ID of this message.
func (m *PingMessage) UniqueID() string {
	return ""
}

// Message body content.
func (m *PingMessage) Content() []byte {
	return nil
}

// The creation time of the message.
func (m *PingMessage) Timestamp() time.Time {
	return time.Now()
}

// The message should not be processed before this timestamp.
func (m *PingMessage) NotBefore() time.Time {
	return time.Now()
}

// Message retry times.
func (m *PingMessage) Retry() int {
	return 0
}

// Return true for a ping message.
func (m *PingMessage) IsPing() bool {
	return true
}

// Begin to process the message.
func (m *PingMessage) Begin() {}

// Indicate the message should be ignored.
func (m *PingMessage) Cancel() {}

// End indicates a successful process.
func (m *PingMessage) End() {}

// Requeue indicates the message should be retried.
func (m *PingMessage) Requeue() {}

// Fail indicates a failed process.
func (m *PingMessage) Fail() {}
