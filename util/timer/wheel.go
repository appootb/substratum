package timer

import (
	"container/list"
	"time"

	ictx "github.com/appootb/substratum/v2/internal/context"
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
	dur time.Duration // used by add/reset; duration is applied on the wheel goroutine
}

type wheel struct {
	interval time.Duration
	ticker   *time.Ticker

	pos    int
	slots  []*list.List
	timers map[*wheelTimer]int // timer → slot index

	eventCh chan *event
}

var (
	// defaultWheel: 1s tick × 600 slots ≈ 10 minutes per full circle before multi-circle.
	defaultWheel = newWheel(time.Second, 600)
)

func newWheel(interval time.Duration, slotNum int) *wheel {
	if interval <= 0 {
		interval = time.Second
	}
	if slotNum <= 0 {
		slotNum = 600
	}
	w := &wheel{
		interval: interval,
		slots:    make([]*list.List, slotNum),
		timers:   make(map[*wheelTimer]int),
		eventCh:  make(chan *event, 1024),
	}
	for i := 0; i < slotNum; i++ {
		w.slots[i] = list.New()
	}
	w.ticker = time.NewTicker(w.interval)
	go w.loop()
	return w
}

func (w *wheel) addTimer(t *wheelTimer, d time.Duration) {
	w.eventCh <- &event{typ: add, wt: t, dur: d}
}

func (w *wheel) removeTimer(t *wheelTimer) {
	w.eventCh <- &event{typ: remove, wt: t}
}

func (w *wheel) resetTimer(t *wheelTimer, d time.Duration) {
	w.eventCh <- &event{typ: reset, wt: t, dur: d}
}

func (w *wheel) loop() {
	for {
		select {
		case <-ictx.Context.Done():
			w.ticker.Stop()
			return

		case <-w.ticker.C:
			w.tickerHandler()

		case evt := <-w.eventCh:
			switch evt.typ {
			case add:
				w.add(evt.wt, evt.dur)
			case remove:
				w.remove(evt.wt)
			case reset:
				w.remove(evt.wt)
				if evt.dur < 0 {
					w.fire(evt.wt)
				} else {
					w.add(evt.wt, evt.dur)
				}
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

		w.fire(wt)

		next := el.Next()
		slot.Remove(el)
		el = next
		delete(w.timers, wt)
	}
}

func (w *wheel) fire(t *wheelTimer) {
	if t.f != nil {
		f := t.f
		go f()
		return
	}
	if t.c == nil {
		return
	}
	select {
	case t.c <- time.Now():
	default:
	}
}

// add schedules t to fire after d. Must run on the wheel loop goroutine.
func (w *wheel) add(t *wheelTimer, d time.Duration) {
	// If already present (should not happen for well-formed use), drop old slot entry.
	if _, ok := w.timers[t]; ok {
		w.removeFromSlot(t)
	}

	t.dur = d
	if d < 0 {
		w.fire(t)
		return
	}

	// Use integer division on duration units to avoid float; keep math in int64
	// so large delays do not overflow int on 32-bit platforms mid-calculation.
	interval := int64(w.interval)
	if interval <= 0 {
		interval = int64(time.Second)
	}
	multi := int64(d) / interval
	if multi < 0 {
		multi = 0
	}
	n := int64(len(w.slots))
	pos := (int64(w.pos) + multi) % n
	t.circle = int(multi / n)
	w.slots[pos].PushBack(t)
	w.timers[t] = int(pos)
}

// remove unschedules t and drains any pending channel notification.
func (w *wheel) remove(t *wheelTimer) {
	w.removeFromSlot(t)
	w.drain(t)
}

func (w *wheel) removeFromSlot(t *wheelTimer) {
	pos, ok := w.timers[t]
	if !ok {
		return
	}
	slot := w.slots[pos]
	for el := slot.Front(); el != nil; el = el.Next() {
		if el.Value.(*wheelTimer) == t {
			slot.Remove(el)
			break
		}
	}
	delete(w.timers, t)
}

func (w *wheel) drain(t *wheelTimer) {
	if t.c == nil {
		return
	}
	select {
	case <-t.c:
	default:
	}
}
