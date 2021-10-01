package datetime

import (
	"testing"
	"time"
)

var (
	_ts = time.Date(2020, 02, 02, 20, 20, 02, 02, time.Local)
)

func TestDay(t *testing.T) {
	beginOfDay := time.Date(2020, 02, 02, 0, 0, 0, 0, time.Local)
	if BeginOfToday(_ts) != beginOfDay {
		t.Fatal(beginOfDay)
	}

	endOfDay := time.Date(2020, 02, 02, 23, 59, 59, 0, time.Local)
	if EndOfToday(_ts) != endOfDay {
		t.Fatal(beginOfDay)
	}

	beginOfWeek := time.Date(2020, 1, 27, 0, 0, 0, 0, time.Local)
	if BeginOfThisWeek(_ts) != beginOfWeek {
		t.Fatal(beginOfWeek)
	}

	endOfWeek := time.Date(2020, 2, 2, 23, 59, 59, 0, time.Local)
	if EndOfThisWeek(_ts) != endOfWeek {
		t.Fatal(endOfWeek)
	}
}

func TestLastMonth(t *testing.T) {
	beginOfMonth := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	if BeginOfTheLastMonth(_ts) != beginOfMonth {
		t.Fatal(beginOfMonth)
	}

	endOfMonth := time.Date(2020, 1, 31, 23, 59, 59, 0, time.Local)
	if EndOfTheLastMonth(_ts) != endOfMonth {
		t.Fatal(endOfMonth)
	}
}

func TestThisMonth(t *testing.T) {
	beginOfMonth := time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local)
	if BeginOfThisMonth(_ts) != beginOfMonth {
		t.Fatal(beginOfMonth)
	}

	endOfMonth := time.Date(2020, 2, 29, 23, 59, 59, 0, time.Local)
	if EndOfThisMonth(_ts) != endOfMonth {
		t.Fatal(endOfMonth)
	}
}

var (
	_nextMonthTimes = [][]time.Time{
		{
			time.Date(2020, 1, 3, 20, 20, 2, 2, time.Local),
			time.Date(2020, 2, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 2, 29, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2020, 2, 2, 20, 20, 2, 2, time.Local),
			time.Date(2020, 3, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 3, 31, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2020, 9, 2, 20, 20, 2, 2, time.Local),
			time.Date(2020, 10, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 10, 31, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2020, 10, 2, 20, 20, 2, 2, time.Local),
			time.Date(2020, 11, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 11, 30, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2020, 11, 2, 20, 20, 2, 2, time.Local),
			time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local),
			time.Date(2020, 12, 31, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2020, 12, 2, 20, 20, 2, 2, time.Local),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 31, 23, 59, 59, 0, time.Local),
		},
	}
)

func TestNextMonth(t *testing.T) {
	for _, monthTimes := range _nextMonthTimes {
		ts, beginOfMonth, endOfMonth := monthTimes[0], monthTimes[1], monthTimes[2]
		if BeginOfTheNextMonth(ts) != beginOfMonth {
			t.Fatal(beginOfMonth, BeginOfTheNextMonth(ts))
		}

		if EndOfTheNextMonth(ts) != endOfMonth {
			t.Fatal(endOfMonth, EndOfTheNextMonth(ts))
		}
	}
}
