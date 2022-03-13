package cache

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkLRU_Rand(b *testing.B) {
	c := New(LRU, WithSize(8192))
	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			c.Set(trace[i], trace[i], DurationPersistence)
		} else {
			_, ok := c.Get(trace[i])
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkLRU_Freq(b *testing.B) {
	c := New(LRU, WithSize(8192))

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(trace[i], trace[i], DurationPersistence)
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		_, ok := c.Get(trace[i])
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func TestLRU(t *testing.T) {
	c := New(LRU, WithSize(128))

	for i := 0; i < 256; i++ {
		c.Set(i, i, DurationPersistence)
	}
	if c.Len() != 128 {
		t.Fatalf("bad len: %v", c.Len())
	}

	for i, k := range c.Keys() {
		if v, ok := c.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := c.Get(i)
		if ok {
			t.Fatal("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := c.Get(i)
		if !ok {
			t.Fatal("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		c.Del(i)
		_, ok := c.Get(i)
		if ok {
			t.Fatal("should be deleted")
		}
	}

	c.Get(192) // expect 192 to be last key in l.Keys()

	for i, k := range c.Keys() {
		if (i < 63 && k != i+193) || (i == 63 && k != 192) {
			t.Fatalf("out of order key: %v", k)
		}
	}

	c.Purge()
	if c.Len() != 0 {
		t.Fatalf("bad len: %v", c.Len())
	}
	if _, ok := c.Get(200); ok {
		t.Fatal("should contain nothing")
	}
}

func TestLRUContains(t *testing.T) {
	c := New(LRU, WithSize(2))

	c.Set(1, 1, DurationPersistence)
	c.Set(2, 2, DurationPersistence)
	if !c.Contain(1) {
		t.Fatal("1 should be contained")
	}

	c.Set(3, 3, DurationPersistence)
	if c.Contain(1) {
		t.Fatal("Contains should not have updated recent-ness of 1")
	}
}

func TestLRUPeek(t *testing.T) {
	c := New(LRU, WithSize(2))

	c.Set(1, 1, DurationPersistence)
	c.Set(2, 2, DurationPersistence)
	if v, ok := c.Peek(1); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	c.Set(3, 3, DurationPersistence)
	if c.Contain(1) {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

func TestLRUTimeout(t *testing.T) {
	c := New(LRU, WithSize(2))

	c.Set(1, 1, DurationPersistence)
	c.Set(2, 2, time.Millisecond*100)
	if !c.Contain(1) || !c.Contain(2) {
		t.Fatal("1 and 2 should be contained")
	}

	time.Sleep(time.Millisecond * 200)
	_, ok := c.Get(2)
	if ok {
		t.Fatal("2 should ge expired")
	}

	if c.Len() != 1 {
		t.Fatalf("bad len: %v", c.Len())
	}

	c.Set(3, 3, DurationPersistence)
	if !c.Contain(1) {
		t.Fatal("Contains should not have updated recent-ness of 2")
	}
}

func TestLRULoad(t *testing.T) {
	count := 0
	fn := func(key interface{}) (interface{}, time.Duration, error) {
		count++
		return 1, DurationPersistence, nil
	}

	c := New(LRU, WithSize(2))

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			v, err := c.GetOrLoad(1, fn)
			if err != nil || v.(int) != 1 {
				panic("1 should be contained")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if count != 1 {
		t.Fatal("should be loaded only once", count)
	}
}

func TestGetOrLoad(t *testing.T) {
	c := New(LRU)
	c.Set("BenchmarkGetOrLoad", "val", time.Second)
	time.Sleep(time.Second)

	count := 0
	run := make(chan bool)
	ready := sync.WaitGroup{}
	done := sync.WaitGroup{}
	//
	begin := time.Now()
	for i := 0; i < 200; i++ {
		ready.Add(1)
		done.Add(1)
		go func() {
			ready.Done()
			<-run

			_, _ = c.GetOrLoad("BenchmarkGetOrLoad", func(_ interface{}) (interface{}, time.Duration, error) {
				time.Sleep(time.Second)
				count++
				return "val", DurationMemoryLock, nil
			})
			done.Done()
		}()
	}

	ready.Wait()
	close(run)
	done.Wait()
	//
	t.Log("Duration: ", time.Now().Sub(begin))
	t.Log("Times: ", count)
	//
	if _, ok := c.Peek("BenchmarkGetOrLoad", true); ok {
		t.Fatal("should be empty")
	}
}
