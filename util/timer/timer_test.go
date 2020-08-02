package timer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	defaultWheel = newWheel(context.Background(), time.Millisecond, 1000)

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestAfterFunc(t *testing.T) {
	i := 10
	c := make(chan bool)
	var f func()
	f = func() {
		i--
		if i >= 0 {
			AfterFunc(0, f)
			time.Sleep(1 * time.Second)
		} else {
			c <- true
		}
	}

	AfterFunc(0, f)
	<-c
}

func benchmark(b *testing.B, bench func(n int)) {
	// Create equal number of garbage timers on each P before starting
	// the benchmark.
	var wg sync.WaitGroup
	garbageAll := make([][]Timer, runtime.GOMAXPROCS(0))
	for i := range garbageAll {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			garbage := make([]Timer, 1<<15)
			for j := range garbage {
				garbage[j] = AfterFunc(time.Hour, nil)
			}
			garbageAll[i] = garbage
		}(i)
	}
	wg.Wait()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bench(1000)
		}
	})
	b.StopTimer()

	for _, garbage := range garbageAll {
		for _, t := range garbage {
			t.Stop()
		}
	}
}

func BenchmarkAfterFunc(b *testing.B) {
	benchmark(b, func(n int) {
		c := make(chan bool)
		var f func()
		f = func() {
			n--
			if n >= 0 {
				AfterFunc(0, f)
			} else {
				c <- true
			}
		}

		AfterFunc(0, f)
		<-c
	})
}

func BenchmarkAfter(b *testing.B) {
	benchmark(b, func(n int) {
		for i := 0; i < n; i++ {
			<-After(1)
		}
	})
}

func BenchmarkStop(b *testing.B) {
	benchmark(b, func(n int) {
		for i := 0; i < n; i++ {
			NewTimer(1 * time.Second).Stop()
		}
	})
}

func BenchmarkSimultaneousAfterFunc(b *testing.B) {
	benchmark(b, func(n int) {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			AfterFunc(0, wg.Done)
		}
		wg.Wait()
	})
}

func BenchmarkStartStop(b *testing.B) {
	benchmark(b, func(n int) {
		timers := make([]Timer, n)
		for i := 0; i < n; i++ {
			timers[i] = AfterFunc(time.Hour, nil)
		}

		for i := 0; i < n; i++ {
			timers[i].Stop()
		}
	})
}

func BenchmarkReset(b *testing.B) {
	benchmark(b, func(n int) {
		t := NewTimer(time.Hour)
		for i := 0; i < n; i++ {
			t.Reset(time.Hour)
		}
		t.Stop()
	})
}

func TestAfter(t *testing.T) {
	const delay = 100 * time.Millisecond
	start := time.Now()
	end := <-After(delay)
	delayadj := delay
	if duration := time.Now().Sub(start); duration < delayadj {
		t.Fatalf("After(%s) slept for only %d ns", delay, duration)
	}
	if min := start.Add(delayadj); end.Before(min) {
		t.Fatalf("After(%s) expect >= %s, got %s", delay, min, end)
	}
}

func TestAfterTick(t *testing.T) {
	const Count = 10
	Delta := 100 * time.Millisecond
	if testing.Short() {
		Delta = 10 * time.Millisecond
	}
	t0 := time.Now()
	for i := 0; i < Count; i++ {
		<-After(Delta)
	}
	t1 := time.Now()
	d := t1.Sub(t0)
	target := Delta * Count
	if d < target*9/10 {
		t.Fatalf("%d ticks of %s too fast: took %s, expected %s", Count, Delta, d, target)
	}
	if !testing.Short() && d > target*30/10 {
		t.Fatalf("%d ticks of %s too slow: took %s, expected %s", Count, Delta, d, target)
	}
}

func TestAfterStop(t *testing.T) {
	var errs []string
	logErrs := func() {
		for _, e := range errs {
			t.Log(e)
		}
	}

	for i := 0; i < 5; i++ {
		AfterFunc(100*time.Millisecond, func() {})
		t0 := NewTimer(50 * time.Millisecond)
		c1 := make(chan bool, 1)
		t1 := AfterFunc(150*time.Millisecond, func() { c1 <- true })
		c2 := After(200 * time.Millisecond)
		t0.Stop()
		t1.Stop()
		<-c2
		select {
		case <-t0.Done():
			errs = append(errs, "event 0 was not stopped")
			continue
		case <-c1:
			errs = append(errs, "event 1 was not stopped")
			continue
		default:
		}
		t1.Stop()

		// Test passed, so all done.
		if len(errs) > 0 {
			t.Logf("saw %d errors, ignoring to avoid flakiness", len(errs))
			logErrs()
		}

		return
	}

	t.Errorf("saw %d errors", len(errs))
	logErrs()
}

func TestAfterQueuing(t *testing.T) {
	// This test flakes out on some systems,
	// so we'll try it a few times before declaring it a failure.
	const attempts = 5
	err := errors.New("!=nil")
	for i := 0; i < attempts && err != nil; i++ {
		delta := time.Duration(20+i*50) * time.Millisecond
		if err = testAfterQueuing(delta); err != nil {
			t.Logf("attempt %v failed: %v", i, err)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

var slots = []int{5, 3, 6, 6, 6, 1, 1, 2, 7, 9, 4, 8, 0}

type afterResult struct {
	slot int
	t    time.Time
}

func await(slot int, result chan<- afterResult, ac <-chan time.Time) {
	result <- afterResult{slot, <-ac}
}

func testAfterQueuing(delta time.Duration) error {
	// make the result channel buffered because we don't want
	// to depend on channel queueing semantics that might
	// possibly change in the future.
	result := make(chan afterResult, len(slots))

	t0 := time.Now()
	for _, slot := range slots {
		go await(slot, result, After(time.Duration(slot)*delta))
	}
	var order []int
	var times []time.Time
	for range slots {
		r := <-result
		order = append(order, r.slot)
		times = append(times, r.t)
	}
	for i := range order {
		if i > 0 && order[i] < order[i-1] {
			return fmt.Errorf("After calls returned out of order: %v", order)
		}
	}
	for i, t := range times {
		dt := t.Sub(t0)
		target := time.Duration(order[i]) * delta
		if dt < target-delta/2 || dt > target+delta*10 {
			return fmt.Errorf("After(%s) arrived at %s, expected [%s,%s]", target, dt, target-delta/2, target+delta*10)
		}
	}
	return nil
}

func TestTimerStopStress(t *testing.T) {
	if testing.Short() {
		return
	}
	for i := 0; i < 100; i++ {
		go func(i int) {
			timer := AfterFunc(2*time.Second, func() {
				t.Errorf("timer %d was not stopped", i)
			})
			time.Sleep(1 * time.Second)
			timer.Stop()
		}(i)
	}
	time.Sleep(3 * time.Second)
}

func testReset(d time.Duration) error {
	t0 := NewTimer(2 * d)
	time.Sleep(d)
	t0.Reset(3 * d)
	time.Sleep(2 * d)
	select {
	case <-t0.Done():
		return errors.New("timer fired early")
	default:
	}
	time.Sleep(2 * d)
	select {
	case <-t0.Done():
	default:
		return errors.New("reset timer did not fire")
	}

	t0.Reset(50 * time.Millisecond)
	return nil
}

func TestReset(t *testing.T) {
	// We try to run this test with increasingly larger multiples
	// until one works so slow, loaded hardware isn't as flaky,
	// but without slowing down fast machines unnecessarily.
	const unit = 25 * time.Millisecond
	tries := []time.Duration{
		1 * unit,
		3 * unit,
		7 * unit,
		15 * unit,
	}
	var err error
	for _, d := range tries {
		err = testReset(d)
		if err == nil {
			t.Logf("passed using duration %v", d)
			return
		}
	}
	t.Error(err)
}

// Test that sleeping for an interval so large it overflows does not
// result in a short sleep duration.
func TestOverflowSleep(t *testing.T) {
	const big = time.Duration(int64(1<<63 - 1))
	select {
	case <-After(big):
		t.Fatalf("big timeout fired")
	case <-After(25 * time.Millisecond):
		// OK
	}
	const neg = time.Duration(-1 << 63)
	select {
	case <-After(neg):
		// OK
	case <-After(1 * time.Second):
		t.Fatalf("negative timeout didn't fire")
	}
}
