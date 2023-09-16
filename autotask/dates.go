package autotask

import "time"

func WeekNo(t time.Time) int {
	// Start by getting January 1 of the current year
	jan1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())

	// Find the next Sunday from January 1
	offset := (7 - int(jan1.Weekday())) % 7
	firstSunday := jan1.AddDate(0, 0, offset)

	// If the date is before the first Sunday, consider it as week 0 (previous year's last week)
	if t.Before(firstSunday) {
		return 1
	}

	// Calculate how many days have passed since the first Sunday
	daysPassed := t.Sub(firstSunday).Hours() / 24

	// Get the week number
	return int(daysPassed/7) + 2
}

func SundayOfTheWeek(t time.Time) time.Time {
	// Subtract the weekday number from the given date.
	// Since Sunday = 0 in time.Weekday, it gives the exact offset we need.
	offset := int(t.Weekday())
	return t.AddDate(0, 0, -offset)
}

func FormatShortDate(t time.Time) string {
	return t.Format("Mon 01/02")
}

func InferYear(month time.Month, windowInMonths int, dt time.Time) int {
	rangeStart := dt.AddDate(0, -windowInMonths, 0)
	rangeEnd := dt.AddDate(0, windowInMonths, 0)

	if rangeStart.Year() == rangeEnd.Year() {
		return rangeStart.Year()
	}

	if month >= rangeStart.Month() && month <= time.December {
		return rangeStart.Year()
	}

	return rangeEnd.Year()
}

func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
