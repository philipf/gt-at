package at

import "time"

// WeekNo calculates the week number of the year based on a provided date.
// The weeks are considered to start on Sundays.
func WeekNo(t time.Time) int {
	// Start by getting January 1 of the current year
	jan1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())

	// Find the next Sunday from January 1
	offset := (7 - int(jan1.Weekday())) % 7
	firstSunday := jan1.AddDate(0, 0, offset)

	// If the date is before the first Sunday, consider it as week 1
	if t.Before(firstSunday) {
		return 1
	}

	// Calculate how many days have passed since the first Sunday
	daysPassed := t.Sub(firstSunday).Hours() / 24

	// Get the week number
	return int(daysPassed/7) + 2
}

// SundayOfTheWeek returns the date of the Sunday of the week based on a provided date.
func SundayOfTheWeek(t time.Time) time.Time {
	// Subtract the weekday number from the given date.
	// Since Sunday = 0 in time.Weekday, it gives the exact offset we need.
	offset := int(t.Weekday())
	return t.AddDate(0, 0, -offset)
}

// InferYear infers the most likely year for a given month, based on a reference date
// and a window in months.
func InferYear(month time.Month, windowInMonths int, dt time.Time) int {
	rangeStart := dt.AddDate(0, -windowInMonths, 0)
	rangeEnd := dt.AddDate(0, windowInMonths, 0)

	// Handle scenario where both range start and end are within the same year
	if rangeStart.Year() == rangeEnd.Year() {
		return rangeStart.Year()
	}

	// Check if the given month is within the range of the start year
	if month >= rangeStart.Month() && month <= time.December {
		return rangeStart.Year()
	}

	// Otherwise, consider it to be in the range of the end year
	return rangeEnd.Year()
}

// Date is a utility function to create a date with time set to midnight and the
// local timezone.
func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
