package timer

import (
	"container/list"
	"context"
	"time"
)

type eventType int

const (
	add eventType = iota
	remove
	reset
)

type event struct {
	typ eventType
	wt  *wheelTimer
}

type wheel struct {
	ctx      context.Context
	stop     context.CancelFunc
	interval time.Duration
	ticker   *time.Ticker

	pos    int
	slots  []*list.List
	timers map[interface{}]int

	eventCh chan *event
}

var (
	defaultWheel = newWheel(context.Background(), time.Second, 600)
)

func newWheel(ctx context.Context, interval time.Duration, slotNum int) *wheel {
	w := &wheel{
		interval: interval,
		slots:    make([]*list.List, slotNum),
		timers:   make(map[interface{}]int),
		eventCh:  make(chan *event, 100),
	}
	w.ctx, w.stop = context.WithCancel(ctx)
	for i := 0; i < slotNum; i++ {
		w.slots[i] = list.New()
	}
	w.ticker = time.NewTicker(w.interval)
	go w.loop()
	return w
}

func (w *wheel) addTimer(t *wheelTimer) {
	w.eventCh <- &event{
		typ: add,
		wt:  t,
	}
}

func (w *wheel) removeTimer(t *wheelTimer) {
	w.eventCh <- &event{
		typ: remove,
		wt:  t,
	}
}

func (w *wheel) resetTimer(t *wheelTimer) {
	w.eventCh <- &event{
		typ: reset,
		wt:  t,
	}
}

func (w *wheel) loop() {
	for {
		select {
		case <-w.ctx.Done():
			w.ticker.Stop()
			return

		case <-w.ticker.C:
			w.tickerHandler()

		case evt := <-w.eventCh:
			switch evt.typ {
			case add:
				w.add(evt.wt)
			case remove:
				w.remove(evt.wt)
			case reset:
				w.remove(evt.wt)
				w.add(evt.wt)
			}
		}
	}
}

func (w *wheel) tickerHandler() {
	slot := w.slots[w.pos]
	w.pos = (w.pos + 1) % len(w.slots)

	for el := slot.Front(); el != nil; {
		wt := el.Value.(*wheelTimer)

		if wt.circle > 0 {
			wt.circle--
			el = el.Next()
			continue
		}

		// Fire event
		w.fire(wt)

		next := el.Next()
		slot.Remove(el)
		el = next
		delete(w.timers, wt)
	}
}

func (w *wheel) fire(t *wheelTimer) {
	if t.f != nil {
		go t.f()
		return
	}

	select {
	case t.c <- time.Now():
	default:
	}
}

func (w *wheel) add(t *wheelTimer) {
	multi := int(t.dur.Nanoseconds() / w.interval.Nanoseconds())
	pos := (w.pos + multi) % len(w.slots)
	t.circle = multi / len(w.slots)
	w.slots[pos].PushBack(t)
	w.timers[t] = pos
}

func (w *wheel) remove(t *wheelTimer) {
	pos, ok := w.timers[t]
	if !ok {
		select {
		case <-t.c:
		default:
		}
		return
	}
	//
	slot := w.slots[pos]
	for el := slot.Front(); el != nil; {
		wt := el.Value.(*wheelTimer)
		if wt == t {
			slot.Remove(el)
			delete(w.timers, t)
			return
		}
		el = el.Next()
	}
}
