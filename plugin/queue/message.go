package queue

import (
	"fmt"
	"time"

	"github.com/appootb/substratum/logger"
)

type Message struct {
	svc       *Debug
	queue     string
	topic     string
	content   []byte
	retry     int
	timestamp time.Time
	delay     time.Duration
}

// Queue name of this message.
func (m *Message) Queue() string {
	return m.queue
}

// Topic name of this message.
func (m *Message) Topic() string {
	return m.topic
}

// Unique ID of this message.
func (m *Message) UniqueID() string {
	return fmt.Sprintf("%s/%s-%d", m.queue, m.topic, m.timestamp.UnixNano())
}

// Message body content.
func (m *Message) Content() []byte {
	return m.content
}

// The creation time of the message.
func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

// The message should not be processed before this timestamp.
func (m *Message) NotBefore() time.Time {
	return m.timestamp.Add(m.delay)
}

// Message retry times.
func (m *Message) Retry() int {
	return m.retry
}

// Return true for a ping message.
func (m *Message) IsPing() bool {
	return false
}

// Begin to process the message.
func (m *Message) Begin() {
}

// Indicate the message should be ignored.
func (m *Message) Cancel() {
}

// End indicates a successful process.
func (m *Message) End() {
}

// Requeue indicates the message should be retried.
func (m *Message) Requeue() {
	m.retry++
	m.svc.requeue(m)
}

// Fail indicates a failed process.
func (m *Message) Fail() {
	logger.Error("queue_failed", logger.Content{
		"queue":     m.queue,
		"topic":     m.topic,
		"timestamp": m.timestamp,
		"retry":     m.retry,
		"delay":     m.delay,
	})
}
