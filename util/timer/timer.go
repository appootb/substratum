package timer

import "time"

type Timer interface {
	Done() <-chan time.Time
	Stop()
	Reset(d time.Duration)
}

func NewTimer(d time.Duration) Timer {
	wt := &wheelTimer{
		c:   make(chan time.Time, 1),
		w:   defaultWheel,
		dur: d,
	}
	if d < 0 {
		wt.c <- time.Now()
	} else {
		defaultWheel.addTimer(wt)
	}
	return wt
}

func After(d time.Duration) <-chan time.Time {
	return NewTimer(d).Done()
}

func AfterFunc(d time.Duration, f func()) Timer {
	wt := &wheelTimer{
		w:   defaultWheel,
		dur: d,
		f:   f,
	}
	if d < 0 {
		go f()
	} else {
		defaultWheel.addTimer(wt)
	}
	return wt
}

// Wrapper for time.Timer.
type Wrapper struct {
	timer *time.Timer
}

func NewWrapper(t *time.Timer) Timer {
	return &Wrapper{
		timer: t,
	}
}

func (t *Wrapper) Done() <-chan time.Time {
	return t.timer.C
}

func (t *Wrapper) Stop() {
	if !t.timer.Stop() {
		<-t.timer.C
	}
}

func (t *Wrapper) Reset(d time.Duration) {
	t.Stop()
	t.timer.Reset(d)
}

type wheelTimer struct {
	w      *wheel
	c      chan time.Time
	f      func()
	dur    time.Duration
	circle int
}

func (t *wheelTimer) Done() <-chan time.Time {
	return t.c
}

func (t *wheelTimer) Stop() {
	t.w.removeTimer(t)
}

func (t *wheelTimer) Reset(d time.Duration) {
	t.dur = d
	t.w.resetTimer(t)
}
