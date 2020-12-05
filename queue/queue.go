package queue

import (
	"context"
	"time"
)

const (
	DefaultTopic       = "default"
	DefaultConcurrency = 1
	DefaultMaxRetry    = 3
)

var (
	impl Queue
)

// Return the service implementor.
func Implementor() Queue {
	return impl
}

// Register service implementor.
func RegisterImplementor(s Queue) {
	impl = s
}

type Queue interface {
	// Publish writes a message body to the specified queue.
	Publish(queue string, content []byte, opts ...PublishOption) error
	// Subscribe consumes the messages of the specified queue.
	Subscribe(queue string, handler Consumer, opts ...SubscribeOption) error
}

type PublishOption func(*PublishOptions)

var EmptyPublishOptions = func() *PublishOptions {
	return &PublishOptions{
		Context: context.Background(),
	}
}

type PublishOptions struct {
	context.Context
	Delay time.Duration
}

func WithPublishContext(ctx context.Context) PublishOption {
	return func(opts *PublishOptions) {
		opts.Context = ctx
	}
}

func WithPublishDelay(delay time.Duration) PublishOption {
	return func(opts *PublishOptions) {
		opts.Delay = delay
	}
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	context.Context
	Component   string
	Topic       string
	Concurrency int
	MaxRetry    int
	Idempotent  Idempotent
}

var EmptySubscribeOptions = func() *SubscribeOptions {
	return &SubscribeOptions{
		Context:     context.Background(),
		Topic:       DefaultTopic,
		Concurrency: DefaultConcurrency,
		MaxRetry:    DefaultMaxRetry,
		Idempotent:  IdempotentImplementor(),
	}
}

func WithConsumeContext(ctx context.Context) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Context = ctx
	}
}

func WithConsumeComponent(component string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Component = component
	}
}

func WithConsumeTopic(topic string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Topic = topic
	}
}

func WithConsumeConcurrency(concurrency int) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Concurrency = concurrency
	}
}

func WithConsumeRetry(retry int) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.MaxRetry = retry
	}
}

func WithConsumeIdempotent(impl Idempotent) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Idempotent = impl
	}
}
