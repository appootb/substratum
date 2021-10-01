package queue

import (
	"context"
	"time"
)

// Consumer handler interface.
type Consumer interface {
	Handle(context.Context, Message) error
}

type ConsumerFunc func(context.Context, Message) error

func (fn ConsumerFunc) Handle(ctx context.Context, m Message) error {
	return fn(ctx, m)
}

// Message interface.
type Message interface {
	// Queue name of this message.
	Queue() string
	// Topic name of this message.
	Topic() string

	// UniqueID returns the unique ID of this message.
	UniqueID() string
	// Content returns the message body content.
	Content() []byte
	// Timestamp indicates the creation time of the message.
	Timestamp() time.Time
	// NotBefore indicates the message should not be processed before this timestamp.
	NotBefore() time.Time
	// Retry times.
	Retry() int
	// IsPing returns true for a ping message.
	IsPing() bool
}

// MessageOperation interface.
type MessageOperation interface {
	// Begin to process the message.
	Begin()
	// Cancel indicates the message should be ignored.
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

// UniqueID returns the unique ID of this message.
func (m *PingMessage) UniqueID() string {
	return ""
}

// Content returns the message body content.
func (m *PingMessage) Content() []byte {
	return nil
}

// Timestamp indicates the creation time of the message.
func (m *PingMessage) Timestamp() time.Time {
	return time.Now()
}

// NotBefore indicates the message should not be processed before this timestamp.
func (m *PingMessage) NotBefore() time.Time {
	return time.Now()
}

// Retry times.
func (m *PingMessage) Retry() int {
	return 0
}

// IsPing returns true for a ping message.
func (m *PingMessage) IsPing() bool {
	return true
}

// Begin to process the message.
func (m *PingMessage) Begin() {}

// Cancel indicates the message should be ignored.
func (m *PingMessage) Cancel() {}

// End indicates a successful process.
func (m *PingMessage) End() {}

// Requeue indicates the message should be retried.
func (m *PingMessage) Requeue() {}

// Fail indicates a failed process.
func (m *PingMessage) Fail() {}
