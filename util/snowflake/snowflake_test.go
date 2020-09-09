package snowflake

import (
	"testing"
	"time"
)

func TestSnowflake_Next(t *testing.T) {
	sf := New()
	vals := make(map[uint64]bool)

	for i := 0; i < 1000; i++ {
		id := sf.Next()
		if _, ok := vals[id]; ok {
			t.Fatal("id exist", id)
		}
		vals[id] = true
	}
}

func TestSnowflake_Timestamp(t *testing.T) {
	sf := New()

	id := sf.Next()
	ts := sf.Timestamp(id)
	diff := time.Now().Sub(ts)
	if diff < 0 {
		t.Fatal("id after now")
	}
	if diff > time.Millisecond {
		t.Fatal("id before 10ms")
	}
	t.Log("diff", diff)
}
