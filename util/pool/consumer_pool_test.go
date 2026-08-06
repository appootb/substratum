package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConsumerPool_Merge(t *testing.T) {
	consumed := make(chan int, 10)
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {
		consumed <- len(v)
	}), WithConsumerMaxMerge(10), WithConsumerMaxDuration(time.Hour))

	size := 10
	for i := 0; i < size; i++ {
		if ok := pool.Add("TestConsumerPool_Merge", i); !ok {
			t.Fatal("add data failed")
		}
	}
	select {
	case l := <-consumed:
		if l != size {
			t.Fatalf("merged len=%d want %d", l, size)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for merge flush")
	}
}

func TestConsumerPool_Merge2(t *testing.T) {
	consumed := make(chan int, 10)
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {
		consumed <- len(v)
	}), WithConsumerMaxMerge(10), WithConsumerMaxDuration(time.Hour))

	size := 20
	got := 0
	for i := 0; i < size; i++ {
		if ok := pool.Add("TestConsumerPool_Merge2", i); !ok {
			t.Fatal("add data failed")
		}
	}
	deadline := time.After(2 * time.Second)
	for got < size {
		select {
		case l := <-consumed:
			got += l
		case <-deadline:
			t.Fatalf("timeout, got %d want %d", got, size)
		}
	}
}

func TestConsumerPool_Duration(t *testing.T) {
	consumed := make(chan int, 10)
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {
		consumed <- len(v)
	}), WithConsumerMaxMerge(100), WithConsumerMaxDuration(100*time.Millisecond))

	size := 9
	start := time.Now()
	for i := 0; i < size; i++ {
		if ok := pool.Add("TestConsumerPool_Duration", i); !ok {
			t.Fatal("add data failed")
		}
	}
	select {
	case l := <-consumed:
		if l != size {
			t.Fatalf("len=%d want %d", l, size)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
	if time.Since(start) < 80*time.Millisecond {
		t.Fatal("flushed too early (expected duration-based flush)")
	}
}

func TestConsumerPoolNegativeKeyNoPanic(t *testing.T) {
	var got int32
	var wg sync.WaitGroup
	wg.Add(1)
	pool := NewConsumer(3, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {
		atomic.AddInt32(&got, int32(len(v)))
		wg.Done()
	}), WithConsumerMaxMerge(1), WithConsumerMaxDuration(time.Hour))

	// Previously panicked: hash.Sum(-1)%3 == -1 → index out of range.
	if !pool.Add(int(-1), "v") {
		t.Fatal("Add(-1) failed")
	}
	wg.Wait()
	if atomic.LoadInt32(&got) != 1 {
		t.Fatal("handler not called")
	}
}

func TestConsumerPoolSlotSizeZeroFallsBack(t *testing.T) {
	var got int32
	var wg sync.WaitGroup
	wg.Add(1)
	pool := NewConsumer(0, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {
		atomic.AddInt32(&got, 1)
		wg.Done()
	}), WithConsumerMaxMerge(1), WithConsumerMaxDuration(time.Hour))

	if len(pool.slots) != DefaultConsumerSlots {
		t.Fatalf("slots=%d want %d", len(pool.slots), DefaultConsumerSlots)
	}
	if !pool.Add("k", 1) {
		t.Fatal("Add failed")
	}
	wg.Wait()
	if atomic.LoadInt32(&got) != 1 {
		t.Fatal("handler not called")
	}
}

func TestConsumerPoolNonPositiveOptionsFallBack(t *testing.T) {
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {}),
		WithConsumerChanLength(0),
		WithConsumerMaxMerge(0),
		WithConsumerMaxDuration(0),
	)
	slot := pool.slots[0]
	if slot.chanLength != DefaultConsumerQueueLength {
		t.Fatalf("chanLength=%d", slot.chanLength)
	}
	if slot.length != DefaultConsumerMaxMerge {
		t.Fatalf("length=%d", slot.length)
	}
	if slot.duration != DefaultConsumerMaxDuration {
		t.Fatalf("duration=%v", slot.duration)
	}
}

func TestConsumerSlotIndexNonNegative(t *testing.T) {
	const n = 3
	for _, key := range []interface{}{int(-1), int64(-100), int(-3), "abc", uint64(9), 0} {
		idx := slotIndex(key, n)
		if idx < 0 || idx >= n {
			t.Fatalf("slotIndex(%v,%d)=%d out of range", key, n, idx)
		}
	}
	if slotIndex("x", 0) != 0 {
		t.Fatal("n<=0 should return 0")
	}
}

func TestConsumerPoolHandlerPanicIsolated(t *testing.T) {
	var okCount int32
	var wg sync.WaitGroup
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, k interface{}, v []interface{}) {
		if v[0].(int) == 0 {
			panic("boom")
		}
		atomic.AddInt32(&okCount, 1)
		wg.Done()
	}), WithConsumerMaxMerge(1), WithConsumerMaxDuration(time.Hour))

	wg.Add(1)
	if !pool.Add("k", 0) {
		t.Fatal("add panic batch")
	}
	time.Sleep(20 * time.Millisecond)
	if !pool.Add("k", 1) {
		t.Fatal("add after panic")
	}
	wg.Wait()
	if atomic.LoadInt32(&okCount) != 1 {
		t.Fatal("subsequent batch should still run")
	}
}
