package valuepb

import "testing"

func TestListRoundTrips(t *testing.T) {
	if got := BoolList(BoolListValue([]bool{true, false})); len(got) != 2 || !got[0] || got[1] {
		t.Fatalf("bool list: %v", got)
	}
	if got := Int64List(Int64ListValue([]int64{1, -2})); len(got) != 2 || got[0] != 1 || got[1] != -2 {
		t.Fatalf("int64 list: %v", got)
	}
	if got := Uint64List(Uint64ListValue([]uint64{3, 4})); len(got) != 2 || got[0] != 3 {
		t.Fatalf("uint64 list: %v", got)
	}
	if got := StringList(StringListValue([]string{"a", "b"})); len(got) != 2 || got[0] != "a" {
		t.Fatalf("string list: %v", got)
	}
	if got := Float64List(Float64ListValue([]float64{1.25})); len(got) != 1 || got[0] != 1.25 {
		t.Fatalf("float64 list: %v", got)
	}
}

func TestListEmptyOnWrongKind(t *testing.T) {
	v := StringValue("nope")
	if len(StringList(v)) != 0 {
		t.Fatal("expected empty list for non-list value")
	}
}
