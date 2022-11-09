package configure

import (
	"reflect"
	"sync"

	"github.com/appootb/substratum/configure"
)

type callback struct {
	sync.RWMutex
	events map[configure.DynamicType][]configure.UpdateEvent

	evt chan configure.DynamicType
	fn  chan *configure.CallbackFunc
}

func newCallback() configure.Callback {
	c := &callback{
		evt:    make(chan configure.DynamicType, 50),
		fn:     make(chan *configure.CallbackFunc, 50),
		events: make(map[configure.DynamicType][]configure.UpdateEvent),
	}
	go c.watchEvent()
	return c
}

func (c *callback) RegChan() chan<- *configure.CallbackFunc {
	return c.fn
}

func (c *callback) EvtChan() chan<- configure.DynamicType {
	return c.evt
}

func (c *callback) watchEvent() {
	for {
		select {
		case fn := <-c.fn:
			if reflect.ValueOf(fn.Value).IsNil() {
				panic("substratum: cannot register callback of a pointer type.")
			}
			c.Lock()
			c.events[fn.Value] = append(c.events[fn.Value], fn.Event)
			c.Unlock()

		case val := <-c.evt:
			c.RLock()
			if events, ok := c.events[val]; ok {
				for _, evt := range events {
					evt()
				}
			}
			c.RUnlock()
		}
	}
}
