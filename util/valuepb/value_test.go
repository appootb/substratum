package valuepb

import (
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestNewInt64Uint64(t *testing.T) {
	v, err := New(int64(42))
	if err != nil {
		t.Fatal(err)
	}
	if Int64(v) != 42 {
		t.Fatalf("Int64=%d", Int64(v))
	}

	v, err = New(uint64(99))
	if err != nil {
		t.Fatal(err)
	}
	if Uint64(v) != 99 {
		t.Fatalf("Uint64=%d", Uint64(v))
	}
}

func TestNewDefault(t *testing.T) {
	v, err := New("hello")
	if err != nil {
		t.Fatal(err)
	}
	if String(v) != "hello" {
		t.Fatalf("String=%q", String(v))
	}
}

func TestScalarRoundTrip(t *testing.T) {
	if !Bool(BoolValue(true)) {
		t.Fatal("bool true")
	}
	if Bool(BoolValue(false)) {
		t.Fatal("bool false")
	}
	if Int64(Int64Value(-7)) != -7 {
		t.Fatal("int64")
	}
	if Uint64(Uint64Value(8)) != 8 {
		t.Fatal("uint64")
	}
	if Float64(Float64Value(1.5)) != 1.5 {
		t.Fatal("float64")
	}
	if String(StringValue("x")) != "x" {
		t.Fatal("string")
	}
}

func TestNilValueAccessors(t *testing.T) {
	var v *structpb.Value
	if Int64(v) != 0 || Uint64(v) != 0 || Float64(v) != 0 || String(v) != "" {
		t.Fatal("nil value should yield zero values")
	}
}

func TestPtrHelpers(t *testing.T) {
	if *BoolPtr(BoolValue(true)) != true {
		t.Fatal("BoolPtr")
	}
	if *Int64Ptr(Int64Value(1)) != 1 {
		t.Fatal("Int64Ptr")
	}
	if *StringPtr(StringValue("a")) != "a" {
		t.Fatal("StringPtr")
	}
}
