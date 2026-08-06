// Package datetime provides calendar helpers and date/time types with JSON codecs.
package datetime

import "time"

const (
	// Day is 24 hours (duration helper only; calendar boundaries use AddDate).
	Day = time.Hour * 24
	// Week is 7 days as a fixed duration (calendar week ends use AddDate).
	Week = Day * 7
)

// BeginOfThisWeek returns 00:00:00 of this Monday in ts's location.
func BeginOfThisWeek(ts time.Time) time.Time {
	offset := int(time.Monday - ts.Weekday())
	if offset > 0 {
		offset = -6
	}
	ts = ts.AddDate(0, 0, offset)
	return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
}

// EndOfThisWeek returns 23:59:59 of this Sunday (calendar week, DST-safe).
func EndOfThisWeek(ts time.Time) time.Time {
	return BeginOfThisWeek(ts).AddDate(0, 0, 7).Add(-time.Second)
}

// BeginOfToday returns 00:00:00 of the calendar day of ts.
func BeginOfToday(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
}

// EndOfToday returns 23:59:59 of the calendar day of ts (DST-safe via AddDate).
func EndOfToday(ts time.Time) time.Time {
	return BeginOfToday(ts).AddDate(0, 0, 1).Add(-time.Second)
}

// BeginOfTheLastMonth returns 00:00:00 on the first day of the previous month.
func BeginOfTheLastMonth(ts time.Time) time.Time {
	year := ts.Year()
	month := ts.Month() - 1
	if month == 0 {
		month = time.December
		year--
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, ts.Location())
}

// EndOfTheLastMonth returns 23:59:59 of the last day of the previous month.
func EndOfTheLastMonth(ts time.Time) time.Time {
	return BeginOfThisMonth(ts).Add(-time.Second)
}

// BeginOfThisMonth returns 00:00:00 on the first day of the month of ts.
func BeginOfThisMonth(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), 1, 0, 0, 0, 0, ts.Location())
}

// EndOfThisMonth returns 23:59:59 of the last day of the month of ts.
func EndOfThisMonth(ts time.Time) time.Time {
	return BeginOfTheNextMonth(ts).Add(-time.Second)
}

// BeginOfTheNextMonth returns 00:00:00 on the first day of the next month.
func BeginOfTheNextMonth(ts time.Time) time.Time {
	year := ts.Year()
	month := ts.Month() + 1
	if month > time.December {
		month -= time.December
		year++
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, ts.Location())
}

// EndOfTheNextMonth returns 23:59:59 of the last day of the next month.
func EndOfTheNextMonth(ts time.Time) time.Time {
	year := ts.Year()
	month := ts.Month() + 2
	if month > time.December {
		month -= time.December
		year++
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, ts.Location()).Add(-time.Second)
}
