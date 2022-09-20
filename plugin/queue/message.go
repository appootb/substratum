package queue

import (
	"fmt"
	"time"

	"github.com/appootb/substratum/v2/logger"
)

type Message struct {
	svc       *Debug
	topic     string
	group     string
	content   []byte
	retry     int
	timestamp time.Time
	delay     time.Duration
}

// Topic name of this message.
func (m *Message) Topic() string {
	return m.topic
}

// Group name of this message.
func (m *Message) Group() string {
	return m.group
}

// Key returns the unique key ID of this message.
func (m *Message) Key() string {
	return fmt.Sprintf("%s/%s-%d", m.topic, m.group, m.timestamp.UnixNano())
}

// Content returns the message body content.
func (m *Message) Content() []byte {
	return m.content
}

// Properties returns the properties of this message.
func (m *Message) Properties() map[string]string {
	return map[string]string{}
}

// Timestamp indicates the creation time of the message.
func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

// NotBefore indicates the message should not be processed before this timestamp.
func (m *Message) NotBefore() time.Time {
	return m.timestamp.Add(m.delay)
}

// Retry times.
func (m *Message) Retry() int {
	return m.retry
}

// IsPing returns true for a ping message.
func (m *Message) IsPing() bool {
	return false
}

// Begin to process the message.
func (m *Message) Begin() {
}

// Cancel indicates the message should be ignored.
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
		"topic":     m.topic,
		"group":     m.group,
		"timestamp": m.timestamp,
		"retry":     m.retry,
		"delay":     m.delay,
	})
}
