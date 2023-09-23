package at

import (
	"testing"
	"time"
)

func TestInferYear(t *testing.T) {
	tests := []struct {
		month          time.Month
		windowInMonths int
		dt             time.Time
		expectedYear   int
	}{
		{
			// When rangeStart and rangeEnd are in the same year.
			month:          time.March,
			windowInMonths: 4,
			dt:             time.Date(2020, time.June, 1, 0, 0, 0, 0, time.Local),
			expectedYear:   2020,
		},
		{
			// When month is within the range of rangeStart.Month() and December.
			month:          time.November,
			windowInMonths: 6,
			dt:             time.Date(2021, time.May, 1, 0, 0, 0, 0, time.Local),
			expectedYear:   2020,
		},
		{
			// When month is outside the above two conditions.
			month:          time.February,
			windowInMonths: 6,
			dt:             time.Date(2022, time.July, 1, 0, 0, 0, 0, time.Local),
			expectedYear:   2022,
		},
	}

	for _, tt := range tests {
		got := InferYear(tt.month, tt.windowInMonths, tt.dt)
		if got != tt.expectedYear {
			t.Errorf("For month %s with window of %d months around %v, expected %d but got %d",
				tt.month, tt.windowInMonths, tt.dt, tt.expectedYear, got)
		}
	}
}
