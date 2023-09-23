package at

import (
	"encoding/json"
	"time"
)

// Credentials holds the authentication details, currently unused in this snippet.
type Credentials struct {
	Username string
	Password string
}

// RequestEntry represents a single entry as received in a JSON request.
type RequestEntry struct {
	Id        int       `json:"id"`
	IsTicket  bool      `json:"isTicket"`
	Date      time.Time `json:"date"`
	StartTime string    `json:"startTime"`
	Duration  float32   `json:"duration"`
	Summary   string    `json:"summary"`
	Project   string    `json:"project"`
}

// UnmarshalToRequestEntries converts JSON data into a slice of RequestEntry.
func UnmarshalToRequestEntries(data []byte) ([]RequestEntry, error) {
	var r []RequestEntry
	err := json.Unmarshal(data, &r)
	return r, err
}

// UnmarshalToTimeEntries converts JSON data into a TimeEntries.
func UnmarshalToTimeEntries(data []byte, dateFormat string) (TimeEntries, error) {
	r, err := UnmarshalToRequestEntries(data)

	if err != nil {
		return nil, err
	}

	var entries TimeEntries

	for _, e := range r {
		te := NewEntry(e.Id, e.IsTicket, e.Date, e.StartTime, e.Duration, e.Summary, e.Project, dateFormat)
		entries = append(entries, te)
	}

	return entries, nil
}
