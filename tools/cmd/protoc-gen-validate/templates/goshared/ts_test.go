package goshared

import (
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestampLitUsesUnixSeconds(t *testing.T) {
	var fns goSharedFuncs
	ts := &timestamppb.Timestamp{Seconds: 1_700_000_000, Nanos: 123}
	got := fns.tsLit(ts)
	want := "time.Unix(1700000000, 123)"
	if got != want {
		t.Fatalf("tsLit=%q want %q", got, want)
	}
}
