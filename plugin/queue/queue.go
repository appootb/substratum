package queue

import (
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/plugin/context"
	"github.com/appootb/substratum/v2/queue"
)

const (
	DebugLog = "_QUEUE_.debug"
	ErrorLog = "_QUEUE_.error"

	LogTopic  = logger.LogTag + "topic"
	LogGroup  = logger.LogTag + "group"
	LogError  = logger.LogTag + "error"
	LogKey    = logger.LogTag + "key"
	LogStatus = logger.LogTag + "status"
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

// Publish writes a message body to the specified topic.
func (m *Queue) Publish(topic string, content []byte, opts ...queue.PublishOption) error {
	options := queue.EmptyPublishOptions()
	for _, o := range opts {
		o(options)
	}
	if err := queue.BackendImplementor().Ping(); err != nil {
		return err
	}
	return queue.BackendImplementor().Write(topic, content, options)
}

// Subscribe consumes the messages of the specified topic.
func (m *Queue) Subscribe(topic string, handler queue.Consumer, opts ...queue.SubscribeOption) error {
	options := queue.EmptySubscribeOptions()
	for _, o := range opts {
		o(options)
	}
	if err := queue.BackendImplementor().Ping(); err != nil {
		return err
	}
	messageChan := make(chan queue.MessageWrapper)
	if err := queue.BackendImplementor().Read(topic, messageChan, options); err != nil {
		return err
	}
	for i := 0; i < options.Concurrency; i++ {
		go m.process(messageChan, handler, options)
	}
	return nil
}

func (m *Queue) process(ch <-chan queue.MessageWrapper, h queue.Consumer, opts *queue.SubscribeOptions) {
	ctx := context.WithImplementContext(opts.Context, opts.Component)

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

		if err != nil {
			logger.Error(ErrorLog, logger.Content{
				LogError:  err.Error(),
				LogTopic:  msg.Topic(),
				LogGroup:  msg.Group(),
				LogKey:    msg.UniqueID(),
				LogStatus: status,
			})
		} else {
			logger.Debug(DebugLog, logger.Content{
				LogTopic:  msg.Topic(),
				LogGroup:  msg.Group(),
				LogKey:    msg.UniqueID(),
				LogStatus: status,
			})
		}
	}
}
