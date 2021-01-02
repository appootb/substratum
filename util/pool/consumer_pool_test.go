package pool

import (
	"context"
	"testing"
	"time"
)

func TestConsumerPool_Merge(t *testing.T) {
	consumed := make(chan int, 10)
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, arg map[interface{}][]interface{}) {
		for _, v := range arg {
			consumed <- len(v)
		}
	}))
	size := 10
	start := time.Now()
	for i := 0; i < size; i++ {
		if ok := pool.Add(i, i); !ok {
			t.Error("add data failed")
		}
	}
	for {
		l := <-consumed
		size -= l
		if size == 0 {
			break
		}
	}
	if time.Now().Sub(start) > time.Second {
		t.Fatal("too long time")
	}
}

func TestConsumerPool_Merge2(t *testing.T) {
	consumed := make(chan int, 10)
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, arg map[interface{}][]interface{}) {
		for _, v := range arg {
			consumed <- len(v)
		}
	}))
	size := 20
	start := time.Now()
	for i := 0; i < size; i++ {
		if ok := pool.Add(i, i); !ok {
			t.Error("add data failed")
		}
	}
	for {
		l := <-consumed
		size -= l
		if size == 0 {
			break
		}
	}
	if time.Now().Sub(start) > time.Second {
		t.Fatal("too long time")
	}
}

func TestConsumerPool_Duration(t *testing.T) {
	consumed := make(chan int, 10)
	pool := NewConsumer(1, ConsumerFunc(func(ctx context.Context, arg map[interface{}][]interface{}) {
		for _, v := range arg {
			consumed <- len(v)
		}
	}))
	size := 9
	start := time.Now()
	for i := 0; i < size; i++ {
		if ok := pool.Add(i, i); !ok {
			t.Error("add data failed")
		}
	}
	for {
		l := <-consumed
		size -= l
		if size == 0 {
			break
		}
	}
	if time.Now().Sub(start) > time.Second {
		t.Fatal("too long time")
	}
}
