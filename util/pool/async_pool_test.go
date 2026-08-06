package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsyncPool(t *testing.T) {
	var b uint32
	var wg sync.WaitGroup
	const workers = 10
	const perWorker = 100
	expected := uint32(workers*perWorker + 1)

	pool := NewAsync(AsyncFunc(func(ctx context.Context, data interface{}) {
		if i, ok := data.(uint32); ok {
			atomic.AddUint32(&b, i)
			wg.Done()
		}
	}), WithAsyncChanLength(workers*perWorker+10))

	wg.Add(int(expected))
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < perWorker; j++ {
				for !pool.Add(uint32(1)) {
					time.Sleep(time.Millisecond)
				}
			}
		}()
	}
	for !pool.Add(uint32(1)) {
		time.Sleep(time.Millisecond)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout, got %d want %d", atomic.LoadUint32(&b), expected)
	}
	if atomic.LoadUint32(&b) != expected {
		t.Fatalf("result not match got=%d want=%d", b, expected)
	}
}

func TestAsyncPoolNonPositiveOptionsFallBack(t *testing.T) {
	var got uint32
	var wg sync.WaitGroup
	wg.Add(1)
	pool := NewAsync(AsyncFunc(func(ctx context.Context, data interface{}) {
		atomic.AddUint32(&got, 1)
		wg.Done()
	}), WithAsyncConcurrency(0), WithAsyncChanLength(0))

	if pool.concurrency != DefaultAsyncConcurrency {
		t.Fatalf("concurrency=%d want %d", pool.concurrency, DefaultAsyncConcurrency)
	}
	if pool.chanLength != DefaultAsyncChanLength {
		t.Fatalf("chanLength=%d want %d", pool.chanLength, DefaultAsyncChanLength)
	}
	if !pool.Add(1) {
		t.Fatal("Add should succeed with fallback workers")
	}
	wg.Wait()
	if atomic.LoadUint32(&got) != 1 {
		t.Fatal("handler not invoked")
	}
}

func TestAsyncPoolAddDropsWhenFull(t *testing.T) {
	// Block the single worker so the buffer fills.
	block := make(chan struct{})
	var startedOnce sync.Once
	started := make(chan struct{})
	pool := NewAsync(AsyncFunc(func(ctx context.Context, data interface{}) {
		startedOnce.Do(func() { close(started) })
		<-block
	}), WithAsyncConcurrency(1), WithAsyncChanLength(1))

	if !pool.Add(1) {
		t.Fatal("first Add")
	}
	<-started
	if !pool.Add(2) {
		t.Fatal("second Add fills buffer")
	}
	if pool.Add(3) {
		t.Fatal("third Add should fail when full")
	}
	close(block)
}

func TestAsyncPoolHandlerPanicDoesNotKillWorker(t *testing.T) {
	var okCount uint32
	var wg sync.WaitGroup
	pool := NewAsync(AsyncFunc(func(ctx context.Context, data interface{}) {
		if data.(int) == 0 {
			panic("boom")
		}
		atomic.AddUint32(&okCount, 1)
		wg.Done()
	}), WithAsyncConcurrency(1), WithAsyncChanLength(8))

	wg.Add(1)
	if !pool.Add(0) {
		t.Fatal("add panic task")
	}
	// Give worker time to recover from panic.
	time.Sleep(20 * time.Millisecond)
	if !pool.Add(1) {
		t.Fatal("add after panic")
	}
	wg.Wait()
	if atomic.LoadUint32(&okCount) != 1 {
		t.Fatal("worker should continue after panic")
	}
}
