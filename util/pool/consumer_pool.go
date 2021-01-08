package pool

import (
	"context"
	"sync"
	"time"

	pctx "github.com/appootb/substratum/plugin/context"
	"github.com/appootb/substratum/util/hash"
)

const (
	DefaultConsumerQueueLength = 1000
	DefaultConsumerMaxMerge    = 10
	DefaultConsumerMaxDuration = time.Second
)

type Consumer interface {
	Handle(context.Context, map[interface{}][]interface{})
}

type ConsumerFunc func(context.Context, map[interface{}][]interface{})

func (fn ConsumerFunc) Handle(ctx context.Context, arg map[interface{}][]interface{}) {
	fn(ctx, arg)
}

type ConsumerOption func(slot *consumerSlot)

func WithConsumerContext(ctx context.Context) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.ctx = ctx
	}
}

func WithConsumerQueueLength(queueLen int) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.length = queueLen
	}
}

func WithConsumerMaxMerge(merge int) ConsumerOption {
	return func(slot *consumerSlot) {
		slot.merge = merge
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

type consumerSlot struct {
	ctx  context.Context
	stop context.CancelFunc

	length    int
	merge     int
	duration  time.Duration
	component string

	ops uint64
	add uint64

	sync.RWMutex
	queue  chan interface{}
	values map[interface{}][]interface{}
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
		ctx:      context.Background(),
		values:   make(map[interface{}][]interface{}),
		length:   DefaultConsumerQueueLength,
		merge:    DefaultConsumerMaxMerge,
		duration: DefaultConsumerMaxDuration,
	}
	for _, o := range opts {
		o(slot)
	}
	slot.ctx, slot.stop = context.WithCancel(slot.ctx)
	slot.queue = make(chan interface{}, slot.length)
	go slot.run(handler)
	return slot
}

func (slot *consumerSlot) run(handler Consumer) {
	for {
		var (
			keys   []interface{}
			sleepy bool
		)

		select {
		case key := <-slot.queue:
			keys = append(keys, key)
		case <-slot.ctx.Done():
			// slot is closing
			select {
			case key := <-slot.queue:
				keys = append(keys, key)
			default:
				return
			}
		}

	MergeLabel:
		for i := 0; i < slot.merge-1; i++ {
			select {
			case key := <-slot.queue:
				keys = append(keys, key)
			default:
				sleepy = true
				break MergeLabel
			}
		}

		values := make(map[interface{}][]interface{}, len(keys))

		slot.Lock()
		for _, key := range keys {
			if v, ok := slot.values[key]; ok {
				if len(v) != 0 {
					values[key] = v
				}
				delete(slot.values, key)
			}
		}
		slot.ops += uint64(len(values))
		slot.Unlock()

		ctx := pctx.WithImplementContext(slot.ctx, slot.component)
		handler.Handle(ctx, values)
		if sleepy {
			time.Sleep(slot.duration)
		}
	}
}

func (pool *ConsumerPool) Add(key interface{}, values ...interface{}) bool {
	idx := hash.Sum(key) % int64(len(pool.slots))
	slot := pool.slots[idx]

	select {
	case <-slot.ctx.Done():
		return false
	default:
	}

	slot.Lock()
	defer slot.Unlock()

	vals, ok := slot.values[key]
	if !ok {
		select {
		case slot.queue <- key:
		default:
			// chan is full
			return false
		}
		slot.values[key] = values
	} else {
		slot.values[key] = append(vals, values...)
	}

	slot.add++
	return true
}
