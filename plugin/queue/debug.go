package queue

import (
	"sync"
	"time"

	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/util/timer"
)

type Debug struct {
	queues sync.Map
}

func (m *Debug) Type() string {
	return "debug"
}

func (m *Debug) Ping() error {
	return nil
}

func (m *Debug) MaxDelay() time.Duration {
	return -1
}

func (m *Debug) GetQueues() ([]string, error) {
	var queues []string
	m.queues.Range(func(queue, _ interface{}) bool {
		queues = append(queues, queue.(string))
		return true
	})
	return queues, nil
}

func (m *Debug) GetTopics() (map[string][]string, error) {
	topics := map[string][]string{}
	m.queues.Range(func(queue, chs interface{}) bool {
		ts, _ := chs.(*sync.Map)
		ts.Range(func(topic, _ interface{}) bool {
			topics[queue.(string)] = append(topics[queue.(string)], topic.(string))
			return true
		})
		return true
	})
	return topics, nil
}

func (m *Debug) GetQueueLength(queue string) (map[string]int64, error) {
	length := map[string]int64{}
	topics, ok := m.queues.Load(queue)
	if !ok {
		return length, nil
	}
	ts, _ := topics.(*sync.Map)
	ts.Range(func(key, value interface{}) bool {
		ch := value.(chan *Message)
		length[key.(string)] = int64(len(ch))
		return true
	})
	return length, nil
}

func (m *Debug) GetTopicLength(queue, topic string) (int64, error) {
	topics, ok := m.queues.Load(queue)
	if !ok {
		return 0, nil
	}
	ts, _ := topics.(*sync.Map)
	ch, ok := ts.Load(topic)
	if !ok {
		return 0, nil
	}
	return int64(len(ch.(chan *Message))), nil
}

func (m *Debug) Read(queue, topic string, ch chan<- queue.MessageWrapper) error {
	topics, ok := m.queues.Load(queue)
	if !ok {
		topics = &sync.Map{}
		m.queues.Store(queue, topics)
	}
	ts, _ := topics.(*sync.Map)
	var cache chan *Message
	if c, ok := ts.Load(topic); !ok {
		cache = make(chan *Message, 100)
		ts.Store(topic, cache)
	} else {
		cache = c.(chan *Message)
	}
	go m.dequeue(ch, cache)
	return nil
}

func (m *Debug) Write(queue string, delay time.Duration, content []byte) error {
	topics, ok := m.queues.Load(queue)
	if !ok {
		topics = &sync.Map{}
		m.queues.Store(queue, topics)
	}
	ts, _ := topics.(*sync.Map)
	ts.Range(func(key, value interface{}) bool {
		ch := value.(chan *Message)
		msg := &Message{
			svc:       m,
			queue:     queue,
			topic:     key.(string),
			content:   content,
			timestamp: time.Now(),
			delay:     delay,
		}
		if delay > 0 {
			timer.AfterFunc(delay, func() {
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
	topics, ok := m.queues.Load(msg.queue)
	if !ok {
		return
	}
	ts, _ := topics.(*sync.Map)
	c, ok := ts.Load(msg.topic)
	if !ok {
		return
	}
	ch := c.(chan *Message)
	m.enqueue(ch, msg)
}
