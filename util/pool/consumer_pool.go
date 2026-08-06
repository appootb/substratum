package pool

import (
	"context"
	"runtime/debug"
	"time"

	sctx "github.com/appootb/substratum/v2/context"
	ictx "github.com/appootb/substratum/v2/internal/context"
	"github.com/appootb/substratum/v2/util/hash"
)

const (
	DefaultConsumerQueueLength = 1000
	DefaultConsumerMaxMerge    = 10
	DefaultConsumerMaxDuration = time.Second
	// DefaultConsumerSlots is used when NewConsumer is given a non-positive slot size.
	DefaultConsumerSlots = 1
)

// Consumer handles a batch of values for one key.
type Consumer interface {
	Handle(context.Context, interface{}, []interface{})
}

// ConsumerFunc adapts a function to Consumer.
type ConsumerFunc func(context.Context, interface{}, []interface{})

// Handle implements Consumer.
func (fn ConsumerFunc) Handle(ctx context.Context, key interface{}, values []interface{}) {
	fn(ctx, key, values)
}

// ConsumerOption configures each consumer slot.
type ConsumerOption func(slot *consumerSlot)

// WithConsumerChanLength sets per-slot queue length (minimum 1).
func WithConsumerChanLength(chanLen int) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.chanLength = chanLen
	}
}

// WithConsumerMaxMerge sets how many values per key trigger an immediate flush (minimum 1).
func WithConsumerMaxMerge(length int) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.length = length
	}
}

// WithConsumerMaxDuration sets the periodic flush interval (must be > 0).
func WithConsumerMaxDuration(dur time.Duration) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.duration = dur
	}
}

// WithConsumerComponent sets the component name for ServerContext in handlers.
func WithConsumerComponent(component string) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.component = component
	}
}

type consumerData struct {
	k, v interface{}
}

type consumerSlot struct {
	component string

	chanLength int
	length     int
	duration   time.Duration

	ch    chan interface{}
	cache map[interface{}][]interface{}
}

// ConsumerPool shards keys across slots and merges values by key before Handle.
type ConsumerPool struct {
	slots []*consumerSlot
}

// NewConsumer creates slotSize independent merge workers.
// Non-positive slotSize falls back to DefaultConsumerSlots.
func NewConsumer(slotSize int, handler Consumer, opts ...ConsumerOption) *ConsumerPool {
	if slotSize <= 0 {
		slotSize = DefaultConsumerSlots
	}
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
	if slot.chanLength <= 0 {
		slot.chanLength = DefaultConsumerQueueLength
	}
	if slot.length <= 0 {
		slot.length = DefaultConsumerMaxMerge
	}
	if slot.duration <= 0 {
		slot.duration = DefaultConsumerMaxDuration
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
				batch := map[interface{}][]interface{}{
					data.k: slot.cache[data.k],
				}
				delete(slot.cache, data.k)
				go slot.callback(handler, batch)
			}
		case <-ticker.C:
			if len(slot.cache) > 0 {
				batch := slot.cache
				slot.cache = map[interface{}][]interface{}{}
				go slot.callback(handler, batch)
			}
		case <-ictx.Context.Done():
			// Flush remaining merged data before exit.
			if len(slot.cache) > 0 {
				slot.callback(handler, slot.cache)
				slot.cache = map[interface{}][]interface{}{}
			}
			return
		}
	}
}

func (slot *consumerSlot) callback(handler Consumer, cache map[interface{}][]interface{}) {
	defer func() {
		if recover() != nil {
			debug.PrintStack()
		}
	}()
	for k, v := range cache {
		handler.Handle(sctx.ServerContext(slot.component), k, v)
	}
}

// slotIndex maps key onto [0, n) using a non-negative remainder.
func slotIndex(key interface{}, n int) int {
	if n <= 0 {
		return 0
	}
	idx := hash.Sum(key) % int64(n)
	if idx < 0 {
		idx += int64(n)
	}
	return int(idx)
}

// Add enqueues a value under key without blocking. Returns false if the slot queue is full.
// Keys are sharded with a non-negative modulo so negative hash sums are safe.
func (pool *ConsumerPool) Add(key interface{}, value interface{}) bool {
	if len(pool.slots) == 0 {
		return false
	}
	slot := pool.slots[slotIndex(key, len(pool.slots))]

	select {
	case slot.ch <- &consumerData{key, value}:
		return true
	default:
		return false
	}
}
