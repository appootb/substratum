package queue

import (
	"context"
	"strconv"
	"time"

	"github.com/appootb/substratum/v2/util/snowflake"
)

type ConsumeOffset int

const (
	ConsumeFromLatest ConsumeOffset = iota
	ConsumeFromEarliest
)

const (
	DefaultGroup       = "default"
	DefaultConcurrency = 1
	DefaultMaxRetry    = 3
)

var (
	impl Queue
)

// Implementor returns the queue service implementor.
func Implementor() Queue {
	return impl
}

// RegisterImplementor registers the queue service implementor.
func RegisterImplementor(s Queue) {
	impl = s
}

// Queue interface.
type Queue interface {
	// Publish writes a message body to the specified topic.
	Publish(topic string, content []byte, opts ...PublishOption) error
	// Subscribe consumes the messages of the specified topic.
	Subscribe(topic string, handler Consumer, opts ...SubscribeOption) error
}

type PublishOption func(*PublishOptions)

var EmptyPublishOptions = func() *PublishOptions {
	key, _ := snowflake.NextID()
	return &PublishOptions{
		Context:    context.Background(),
		Sequence:   key,
		Key:        strconv.FormatUint(key, 10),
		Properties: map[string]string{},
	}
}

type PublishOptions struct {
	context.Context
	Sequence   uint64
	Key        string
	Delay      time.Duration
	Properties map[string]string
}

func WithPublishContext(ctx context.Context) PublishOption {
	return func(opts *PublishOptions) {
		opts.Context = ctx
	}
}

func WithPublishSequence(seq uint64) PublishOption {
	return func(opts *PublishOptions) {
		opts.Sequence = seq
	}
}

func WithPublishUniqueKey(key string) PublishOption {
	return func(opts *PublishOptions) {
		opts.Key = key
	}
}

func WithPublishDelay(delay time.Duration) PublishOption {
	return func(opts *PublishOptions) {
		opts.Delay = delay
	}
}

func WithProperty(key, value string) PublishOption {
	return func(opts *PublishOptions) {
		if opts.Properties != nil {
			opts.Properties[key] = value
			return
		}
		opts.Properties = map[string]string{
			key: value,
		}
	}
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	context.Context
	Component   string
	Group       string
	Concurrency int
	MaxRetry    int
	InitOffset  ConsumeOffset
	Idempotent  Idempotent
}

var EmptySubscribeOptions = func() *SubscribeOptions {
	return &SubscribeOptions{
		Context:     context.Background(),
		Group:       DefaultGroup,
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

func WithConsumeGroup(name string) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Group = name
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

func WithInitOffset(offset ConsumeOffset) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.InitOffset = offset
	}
}

func WithConsumeIdempotent(impl Idempotent) SubscribeOption {
	return func(opts *SubscribeOptions) {
		opts.Idempotent = impl
	}
}
