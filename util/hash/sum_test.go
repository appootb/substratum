package hash

import "testing"

func TestSumInts(t *testing.T) {
	if Sum(42) != 42 {
		t.Fatal(Sum(42))
	}
	if Sum(int64(7)) != 7 {
		t.Fatal(Sum(int64(7)))
	}
	if Sum(uint32(9)) != 9 {
		t.Fatal(Sum(uint32(9)))
	}
	if Sum(uint8(3)) != 3 {
		t.Fatalf("uint8: got %d", Sum(uint8(3)))
	}
}

func TestSumStringNumeric(t *testing.T) {
	if Sum("123") != 123 {
		t.Fatalf("got %d", Sum("123"))
	}
}

func TestSumStringHashStable(t *testing.T) {
	a := Sum("abc")
	b := Sum("abc")
	if a != b {
		t.Fatal("checksum not stable")
	}
	if a == 0 {
		t.Fatal("unexpected zero checksum for non-empty string")
	}
}

func TestSumPointer(t *testing.T) {
	v := 99
	if Sum(&v) != 99 {
		t.Fatalf("got %d", Sum(&v))
	}
}
