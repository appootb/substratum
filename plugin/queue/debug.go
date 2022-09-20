package queue

import (
	"sync"
	"time"

	"github.com/appootb/substratum/v2/configure"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/util/timer"
)

type Debug struct {
	queues sync.Map
}

// Init queue backend instance.
func (m *Debug) Init(configure.Address) error {
	return nil
}

func (m *Debug) Type() string {
	return "debug"
}

func (m *Debug) Ping() error {
	return nil
}

func (m *Debug) Read(topic string, ch chan<- queue.MessageWrapper, opts *queue.SubscribeOptions) error {
	groups, ok := m.queues.Load(topic)
	if !ok {
		groups = &sync.Map{}
		m.queues.Store(topic, groups)
	}
	gs, _ := groups.(*sync.Map)
	var cache chan *Message
	if c, ok := gs.Load(opts.Group); !ok {
		cache = make(chan *Message, 100)
		gs.Store(opts.Group, cache)
	} else {
		cache = c.(chan *Message)
	}
	go m.dequeue(ch, cache)
	return nil
}

func (m *Debug) Write(topic string, content []byte, opts *queue.PublishOptions) error {
	groups, ok := m.queues.Load(topic)
	if !ok {
		groups = &sync.Map{}
		m.queues.Store(topic, groups)
	}
	gs, _ := groups.(*sync.Map)
	gs.Range(func(key, value interface{}) bool {
		ch := value.(chan *Message)
		msg := &Message{
			svc:       m,
			topic:     topic,
			group:     key.(string),
			content:   content,
			timestamp: time.Now(),
			delay:     opts.Delay,
		}
		if opts.Delay > 0 {
			timer.AfterFunc(opts.Delay, func() {
				m.enqueue(ch, msg)
			})
		} else {
			m.enqueue(ch, msg)
		}
		return true
	})
	return nil
}

func (m *Debug) dequeue(in chan<- queue.MessageWrapper, out <-chan *Message) {
	ping := &queue.PingMessage{}

	for {
		in <- ping
		in <- <-out
	}
}

func (m *Debug) enqueue(ch chan *Message, msg *Message) {
	for {
		select {
		case ch <- msg:
			return
		default:
			<-ch
		}
	}
}

func (m *Debug) requeue(msg *Message) {
	groups, ok := m.queues.Load(msg.topic)
	if !ok {
		return
	}
	ts, _ := groups.(*sync.Map)
	c, ok := ts.Load(msg.group)
	if !ok {
		return
	}
	ch := c.(chan *Message)
	m.enqueue(ch, msg)
}
