package queue

import (
	"errors"
	"fmt"

	"github.com/appootb/substratum/plugin/context"
	"github.com/appootb/substratum/queue"
)

func Init() {
	if queue.IdempotentImplementor() == nil {
		queue.RegisterIdempotentImplementor(&EmptyIdempotent{})
	}
	if queue.BackendImplementor() == nil {
		queue.RegisterBackendImplementor(&Debug{})
	}
	if queue.Implementor() == nil {
		queue.RegisterImplementor(&Queue{})
	}
}

type Queue struct{}

// Publish writes a message body to the specified queue name.
func (m *Queue) Publish(name string, content []byte, opts ...queue.PublishOption) error {
	options := queue.EmptyPublishOptions()
	for _, o := range opts {
		o(options)
	}
	maxDelay := queue.BackendImplementor().MaxDelay()
	if options.Delay > 0 && maxDelay >= 0 && options.Delay > maxDelay {
		return errors.New(fmt.Sprintf("substratum queue delay: %v, max supported: %v", options.Delay, maxDelay))
	}
	if err := queue.BackendImplementor().Ping(); err != nil {
		return err
	}
	return queue.BackendImplementor().Write(options.Context, name, options.Delay, content)
}

// Subscribe consumes the messages of the specified queue.
func (m *Queue) Subscribe(name string, handler queue.Consumer, opts ...queue.SubscribeOption) error {
	options := queue.EmptySubscribeOptions()
	for _, o := range opts {
		o(options)
	}
	if err := queue.BackendImplementor().Ping(); err != nil {
		return err
	}
	messageChan := make(chan queue.MessageWrapper)
	if err := queue.BackendImplementor().Read(options.Context, name, options.Topic, messageChan); err != nil {
		return err
	}
	for i := 0; i < options.Concurrency; i++ {
		go m.process(messageChan, handler, options)
	}
	return nil
}

func (m *Queue) process(ch <-chan queue.MessageWrapper, h queue.Consumer, opts *queue.SubscribeOptions) {
	ctx := context.WithImplementContext(opts.Context, opts.Component, opts.Product)

	for {
		var (
			err    error
			msg    queue.MessageWrapper
			status queue.ProcessStatus
		)

		select {
		case <-opts.Context.Done():
			return
		case msg = <-ch:
		}

		if msg.IsPing() {
			continue
		}

		msg.Begin()
		status = queue.Processing

		if !opts.Idempotent.BeforeProcess(msg) {
			msg.Cancel()
			status = queue.Canceled
			goto ProcessEnd
		}

		err = h.Handle(ctx, msg)
		if err == nil {
			msg.End()
			status = queue.Succeeded
			goto ProcessEnd
		}

		if msg.Retry() < opts.MaxRetry {
			msg.Requeue()
			status = queue.Requeued
			goto ProcessEnd
		}

		msg.Fail()
		status = queue.Failed

	ProcessEnd:
		opts.Idempotent.AfterProcess(msg, status)
	}
}
