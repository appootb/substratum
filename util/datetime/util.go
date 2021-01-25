package datetime

import "time"

const (
	Day  = time.Hour * 24
	Week = Day * 7
)

func BeginOfThisWeek(ts time.Time) time.Time {
	offset := int(time.Monday - ts.Weekday())
	if offset > 0 {
		offset = -6
	}
	ts = ts.AddDate(0, 0, offset)
	return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
}

func EndOfThisWeek(ts time.Time) time.Time {
	return BeginOfThisWeek(ts.Add(Week)).Add(-time.Second)
}

func BeginOfToday(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
}

func EndOfToday(ts time.Time) time.Time {
	return BeginOfToday(ts.Add(Day)).Add(-time.Second)
}

func BeginOfTheLastMonth(ts time.Time) time.Time {
	year := ts.Year()
	month := (ts.Month() - 1) % time.December
	if month == 0 {
		month = time.December
		year--
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, ts.Location())
}

func EndOfTheLastMonth(ts time.Time) time.Time {
	return BeginOfThisMonth(ts).Add(-time.Second)
}

func BeginOfThisMonth(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), 1, 0, 0, 0, 0, ts.Location())
}

func EndOfThisMonth(ts time.Time) time.Time {
	return BeginOfTheNextMonth(ts).Add(-time.Second)
}

func BeginOfTheNextMonth(ts time.Time) time.Time {
	year := ts.Year()
	month := (ts.Month() + 1) % time.December
	if month == time.January {
		year++
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, ts.Location())
}

func EndOfTheNextMonth(ts time.Time) time.Time {
	year := ts.Year()
	month := (ts.Month() + 2) % time.December
	if month == time.January || month == time.February {
		year++
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, ts.Location()).Add(-time.Second)
}
