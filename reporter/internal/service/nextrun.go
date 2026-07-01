package service

import (
	"time"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

func NextRun(kind, timeOfDay string, dayOfWeek, dayOfMonth int, runAt *time.Time, from time.Time, timezone string) (*time.Time, bool) {
	loc := loadLocation(timezone)
	from = from.UTC()

	switch kind {
	case defs.KindOnce:
		return nextOnce(runAt)
	case defs.KindDaily:
		return nextDaily(timeOfDay, from, loc)
	case defs.KindWeekly:
		next := nextWeekday(from, timeOfDay, dayOfWeek, loc)
		return &next, true
	case defs.KindMonthly:
		next := nextMonthday(from, timeOfDay, dayOfMonth, loc)
		return &next, true
	}

	return nil, false
}

func loadLocation(name string) *time.Location {
	if name == "" {
		return time.UTC
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}

	return loc
}

func nextOnce(runAt *time.Time) (*time.Time, bool) {
	if runAt == nil {
		return nil, false
	}

	runAtUTC := runAt.UTC()
	return &runAtUTC, true
}

func nextDaily(timeOfDay string, from time.Time, loc *time.Location) (*time.Time, bool) {
	next := atTime(from, timeOfDay, loc)
	if !next.After(from) {
		next = next.AddDate(0, 0, 1)
	}

	utc := next.UTC()
	return &utc, true
}

func parseTimeOfDay(value string) (int, int) {
	t, err := time.Parse("15:04", value)
	if err != nil {
		return 0, 0
	}
	return t.Hour(), t.Minute()
}

func atTime(base time.Time, timeOfDay string, loc *time.Location) time.Time {
	hour, minute := parseTimeOfDay(timeOfDay)
	local := base.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), hour, minute, 0, 0, loc)
}

func nextWeekday(from time.Time, timeOfDay string, dayOfWeek int, loc *time.Location) time.Time {
	target := time.Weekday(dayOfWeek % 7)
	candidate := atTime(from, timeOfDay, loc)

	for i := 0; i < 8; i++ {
		if candidate.Weekday() == target && candidate.After(from) {
			return candidate.UTC()
		}
		candidate = candidate.AddDate(0, 0, 1)
	}

	return candidate.UTC()
}

func nextMonthday(from time.Time, timeOfDay string, dayOfMonth int, loc *time.Location) time.Time {
	if dayOfMonth < 1 {
		dayOfMonth = 1
	}

	hour, minute := parseTimeOfDay(timeOfDay)
	local := from.In(loc)
	year, month := local.Year(), local.Month()

	for i := 0; i < 24; i++ {
		day := clampDay(year, month, dayOfMonth)
		candidate := time.Date(year, month, day, hour, minute, 0, 0, loc)
		if candidate.After(from) {
			return candidate.UTC()
		}

		year, month = nextMonth(year, month)
	}

	return atTime(from, timeOfDay, loc).UTC()
}

func nextMonth(year int, month time.Month) (int, time.Month) {
	month++
	if month > time.December {
		return year + 1, time.January
	}
	return year, month
}

func clampDay(year int, month time.Month, day int) int {
	last := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > last {
		return last
	}
	return day
}
