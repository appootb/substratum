package queue

import "time"

// Queue message handler.
type MessageHandler func(Message) error

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
