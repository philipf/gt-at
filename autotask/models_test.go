package autotask

import (
	"testing"
	"time"
)

func TestCalculateDerivedDurations(t *testing.T) {
	tests := []struct {
		duration           float32
		expectedHours      int
		expectedMinutes    float32
		expectedHoursStr   string
		expectedMinutesStr string
	}{
		{1, 1, 0, "1", "0"},
		{2.25, 2, 15, "2", "15"},
		{1.5, 1, 30, "1", "30"},
		{1.75, 1, 45, "1", "45"},
		{3.5, 3, 30, "3", "30"},
		{4.5, 4, 30, "4", "30"},
		{7.7, 7, 42, "7", "42"},
	}

	for _, test := range tests {
		te := &TimeEntry{
			Id:              1,
			Date:            time.Now(),
			Exists:          true,
			Duration:        test.duration,
			WeekNo:          1,
			WeekPeerLocator: nil,
		}

		te.CalculateDerivedDurations()

		if te.DurationHours != test.expectedHours {
			t.Errorf("Expected DurationHours: %d, but got %d for Duration: %f", test.expectedHours, te.DurationHours, test.duration)
		}

		if te.DurationMinutes != test.expectedMinutes {
			t.Errorf("Expected DurationMinutes: %f, but got %f for Duration: %f", test.expectedMinutes, te.DurationMinutes, test.duration)
		}

		if te.DurationHoursStr != test.expectedHoursStr {
			t.Errorf("Expected DurationHoursStr: %s, but got %s for Duration: %f", test.expectedHoursStr, te.DurationHoursStr, test.duration)
		}

		if te.DurationMinutesStr != test.expectedMinutesStr {
			t.Errorf("Expected DurationMinutesStr: %s, but got %s for Duration: %f", test.expectedMinutesStr, te.DurationMinutesStr, test.duration)
		}
	}
}
