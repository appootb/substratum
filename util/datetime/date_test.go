package datetime

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateJSONRoundTrip(t *testing.T) {
	ts := time.Date(2024, 5, 6, 12, 0, 0, 0, time.Local)
	d := NewDate(ts)
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	var out Date
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "2024-05-06" {
		t.Fatalf("got %s", out.String())
	}
	if !out.Time().Equal(time.Date(2024, 5, 6, 0, 0, 0, 0, time.Local)) {
		t.Fatalf("time=%v", out.Time())
	}
}

func TestDateEmptyAndInvalid(t *testing.T) {
	var d Date
	if err := d.UnmarshalJSON([]byte(`""`)); err != nil {
		t.Fatal(err)
	}
	if err := d.UnmarshalJSON([]byte(`"not-a-date"`)); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestDateCustomLayout(t *testing.T) {
	ts := time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local)
	d := NewDate(ts, "02/01/2006")
	if d.String() != "02/01/2024" {
		t.Fatalf("got %s", d.String())
	}
}

func TestDateNewTruncatesTimeOfDay(t *testing.T) {
	ts := time.Date(2024, 5, 6, 15, 4, 5, 0, time.Local)
	d := NewDate(ts)
	if d.Time().Hour() != 0 || d.Time().Minute() != 0 || d.Time().Second() != 0 {
		t.Fatalf("expected midnight, got %v", d.Time())
	}
	if d.String() != "2024-05-06" {
		t.Fatal(d.String())
	}
}

func TestDateMarshalURLMatchesString(t *testing.T) {
	ts := time.Date(2024, 5, 6, 12, 30, 0, 0, time.Local)
	d := NewDate(ts)
	s, err := d.MarshalURL()
	if err != nil {
		t.Fatal(err)
	}
	if s != d.String() || s != "2024-05-06" {
		t.Fatalf("MarshalURL=%q String=%q", s, d.String())
	}
}

func TestDateUnmarshalURL(t *testing.T) {
	var d Date
	if err := d.UnmarshalURL("2024-08-01"); err != nil {
		t.Fatal(err)
	}
	if d.String() != "2024-08-01" {
		t.Fatal(d.String())
	}
	if err := d.UnmarshalURL("bad"); err == nil {
		t.Fatal("expected error")
	}
}

func TestDateTimeJSONRoundTrip(t *testing.T) {
	ts := time.Date(2024, 5, 6, 7, 8, 9, 0, time.Local)
	dt := NewDateTime(ts)
	b, err := json.Marshal(dt)
	if err != nil {
		t.Fatal(err)
	}
	var out DateTime
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "2024-05-06 07:08:09" {
		t.Fatalf("got %s", out.String())
	}
}

func TestDateTimeEmptyAndInvalid(t *testing.T) {
	var dt DateTime
	if err := dt.UnmarshalJSON([]byte(`""`)); err != nil {
		t.Fatal(err)
	}
	if err := dt.UnmarshalJSON([]byte(`"bad"`)); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestDateTimeMarshalURLMatchesString(t *testing.T) {
	ts := time.Date(2024, 5, 6, 7, 8, 9, 0, time.Local)
	dt := NewDateTime(ts)
	s, err := dt.MarshalURL()
	if err != nil {
		t.Fatal(err)
	}
	if s != dt.String() {
		t.Fatalf("MarshalURL=%q String=%q", s, dt.String())
	}
}

func TestZeroValueStringUsesDefaultLayout(t *testing.T) {
	var d Date
	// zero time formatted with default layout
	if d.String() == "" {
		// zero time still formats; ensure no panic and default layout applied
		_ = d.String()
	}
	// After setting layout empty explicitly
	d.layout = ""
	s := d.String()
	if s == "" && !d.payload.IsZero() {
		t.Fatal("expected default layout formatting")
	}
	// zero payload with empty layout should still use DefaultDateLayout
	var d2 Date
	got := d2.String()
	want := d2.payload.Format(DefaultDateLayout)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
