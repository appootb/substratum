package pool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsyncPool(t *testing.T) {
	b := uint32(0)
	pool := NewAsync(AsyncFunc(func(ctx context.Context, data interface{}) {
		if i, ok := data.(uint32); ok {
			atomic.AddUint32(&b, i)
		} else {
			t.Error("data is not int type")
		}
	}))
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				if ok := pool.Add(uint32(1)); !ok {
					t.Error("add data failed")
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}
	pool.Add(uint32(1))
	time.Sleep(2 * time.Second)
	if b != 10001 {
		t.Error("result not match", b, 10001)
	}
}
