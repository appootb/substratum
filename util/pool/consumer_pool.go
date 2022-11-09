package pool

import (
	"context"
	"time"

	sctx "github.com/appootb/substratum/context"
	ictx "github.com/appootb/substratum/internal/context"
	"github.com/appootb/substratum/util/hash"
)

const (
	DefaultConsumerQueueLength = 1000
	DefaultConsumerMaxMerge    = 10
	DefaultConsumerMaxDuration = time.Second
)

type Consumer interface {
	Handle(context.Context, interface{}, []interface{})
}

type ConsumerFunc func(context.Context, interface{}, []interface{})

func (fn ConsumerFunc) Handle(ctx context.Context, key interface{}, values []interface{}) {
	fn(ctx, key, values)
}

type ConsumerOption func(slot *consumerSlot)

func WithConsumerChanLength(chanLen int) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.chanLength = chanLen
	}
}

func WithConsumerMaxMerge(length int) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.length = length
	}
}

func WithConsumerMaxDuration(dur time.Duration) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.duration = dur
	}
}

func WithConsumerComponent(component string) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.component = component
	}
}

func WithConsumerProduct(product string) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.product = product
	}
}

type consumerData struct {
	k, v interface{}
}

type consumerSlot struct {
	component string
	product   string

	chanLength int
	length     int
	duration   time.Duration

	ch    chan interface{}
	cache map[interface{}][]interface{}
}

type ConsumerPool struct {
	slots []*consumerSlot
}

func NewConsumer(slotSize int, handler Consumer, opts ...ConsumerOption) *ConsumerPool {
	pool := &ConsumerPool{
		slots: make([]*consumerSlot, 0, slotSize),
	}
	for i := 0; i < slotSize; i++ {
		pool.slots = append(pool.slots, newConsumerSlot(handler, opts))
	}
	return pool
}

func newConsumerSlot(handler Consumer, opts []ConsumerOption) *consumerSlot {
	slot := &consumerSlot{
		cache:      map[interface{}][]interface{}{},
		chanLength: DefaultConsumerQueueLength,
		length:     DefaultConsumerMaxMerge,
		duration:   DefaultConsumerMaxDuration,
	}
	for _, o := range opts {
		o(slot)
	}
	slot.ch = make(chan interface{}, slot.chanLength)
	go slot.run(handler)
	return slot
}

func (slot *consumerSlot) run(handler Consumer) {
	ticker := time.NewTicker(slot.duration)
	defer ticker.Stop()

	for {
		select {
		case v := <-slot.ch:
			data := v.(*consumerData)
			slot.cache[data.k] = append(slot.cache[data.k], data.v)
			if len(slot.cache[data.k]) >= slot.length {
				go slot.callback(handler, map[interface{}][]interface{}{
					data.k: slot.cache[data.k],
				})
				delete(slot.cache, data.k)
			}
		case <-ticker.C:
			if len(slot.cache) > 0 {
				go slot.callback(handler, slot.cache)
				slot.cache = map[interface{}][]interface{}{}
			}
		case <-ictx.Context.Done():
			return
		}
	}
}

func (slot *consumerSlot) callback(handler Consumer, cache map[interface{}][]interface{}) {
	for k, v := range cache {
		handler.Handle(sctx.ServerContext(slot.component), k, v)
	}
}

func (pool *ConsumerPool) Add(key interface{}, value interface{}) bool {
	idx := hash.Sum(key) % int64(len(pool.slots))
	slot := pool.slots[idx]

	select {
	case slot.ch <- &consumerData{key, value}:
		return true
	default:
		return false
	}
}
