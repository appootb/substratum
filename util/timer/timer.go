// Package timer provides a hierarchical-style timing wheel for lightweight
// timers (After / AfterFunc / NewTimer), plus a thin Wrapper around time.Timer.
package timer

import "time"

// Timer is a cancellable delayed event.
type Timer interface {
	Done() <-chan time.Time
	Stop()
	Reset(d time.Duration)
}

// NewTimer creates a wheel-based timer that fires after d.
// If d < 0, the timer is already expired (Done is ready immediately).
func NewTimer(d time.Duration) Timer {
	wt := &wheelTimer{
		c: make(chan time.Time, 1),
		w: defaultWheel,
	}
	if d < 0 {
		wt.c <- time.Now()
	} else {
		defaultWheel.addTimer(wt, d)
	}
	return wt
}

// After waits for the duration to elapse and then sends the current time on
// the returned channel. The underlying timer cannot be stopped.
func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).Done()
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// If d < 0, f is started immediately (when non-nil).
func AfterFunc(d time.Duration, f func()) Timer {
	wt := &wheelTimer{
		w: defaultWheel,
		f: f,
	}
	if d < 0 {
		if f != nil {
			go f()
		}
		return wt
	}
	defaultWheel.addTimer(wt, d)
	return wt
}

// Wrapper adapts a standard library *time.Timer to Timer.
type Wrapper struct {
	timer *time.Timer
}

// NewWrapper wraps t.
func NewWrapper(t *time.Timer) Timer {
	return &Wrapper{timer: t}
}

// Done returns the underlying timer channel.
func (t *Wrapper) Done() <-chan time.Time {
	return t.timer.C
}

// Stop stops the timer. If the timer has already fired, a pending value on C
// is drained without blocking forever when C was already received.
func (t *Wrapper) Stop() {
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
		default:
		}
	}
}

// Reset stops the timer and schedules it to fire after d.
func (t *Wrapper) Reset(d time.Duration) {
	t.Stop()
	t.timer.Reset(d)
}

// wheelTimer is a timer scheduled on a timing wheel.
// Fields other than immutable w/c/f are owned by the wheel loop goroutine
// after the timer is first scheduled (duration is passed via events).
type wheelTimer struct {
	w      *wheel
	c      chan time.Time // nil for AfterFunc timers
	f      func()         // nil for channel timers
	dur    time.Duration  // last scheduled duration; written only on wheel goroutine
	circle int
}

// Done returns the completion channel (nil channel if this is an AfterFunc timer).
func (t *wheelTimer) Done() <-chan time.Time {
	return t.c
}

// Stop removes the timer from the wheel and drains any pending channel tick.
func (t *wheelTimer) Stop() {
	t.w.removeTimer(t)
}

// Reset reschedules the timer to fire after d (wheel-owned; no data race on dur).
func (t *wheelTimer) Reset(d time.Duration) {
	t.w.resetTimer(t, d)
}
