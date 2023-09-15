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
