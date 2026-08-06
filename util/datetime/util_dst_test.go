package datetime

import (
	"testing"
	"time"
)

func TestEndOfTodayDSTFallBack(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip(err)
	}
	// 2024-11-03 fall back: day has 25 hours; fixed 24h arithmetic was wrong.
	ts := time.Date(2024, 11, 3, 0, 30, 0, 0, loc)
	end := EndOfToday(ts)
	want := time.Date(2024, 11, 3, 23, 59, 59, 0, loc)
	if !end.Equal(want) {
		t.Fatalf("EndOfToday=%v want %v", end, want)
	}
}

func TestEndOfTodayDSTSpringForward(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip(err)
	}
	// 2024-03-10 spring forward: day has 23 hours.
	ts := time.Date(2024, 3, 10, 10, 0, 0, 0, loc)
	end := EndOfToday(ts)
	want := time.Date(2024, 3, 10, 23, 59, 59, 0, loc)
	if !end.Equal(want) {
		t.Fatalf("EndOfToday=%v want %v", end, want)
	}
}

func TestEndOfThisWeekCalendar(t *testing.T) {
	// Sunday 2024-05-05
	sun := time.Date(2024, 5, 5, 12, 0, 0, 0, time.Local)
	begin := BeginOfThisWeek(sun)
	wantBegin := time.Date(2024, 4, 29, 0, 0, 0, 0, time.Local)
	if !begin.Equal(wantBegin) {
		t.Fatalf("BeginOfThisWeek=%v want %v", begin, wantBegin)
	}
	end := EndOfThisWeek(sun)
	wantEnd := time.Date(2024, 5, 5, 23, 59, 59, 0, time.Local)
	if !end.Equal(wantEnd) {
		t.Fatalf("EndOfThisWeek=%v want %v", end, wantEnd)
	}
}

func TestEndOfTodayUsesCalendarNotFixedDuration(t *testing.T) {
	ts := time.Date(2024, 5, 6, 10, 0, 0, 0, time.Local)
	end := EndOfToday(ts)
	next := BeginOfToday(ts).AddDate(0, 0, 1)
	if end.After(next) || !end.Before(next) {
		// end must be exactly 1s before next midnight
		if !end.Equal(next.Add(-time.Second)) {
			t.Fatalf("end=%v next=%v", end, next)
		}
	}
	if !end.Equal(next.Add(-time.Second)) {
		t.Fatalf("end=%v want %v", end, next.Add(-time.Second))
	}
}
