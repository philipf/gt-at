package at

import (
	"encoding/json"
	"time"
)

type Credentials struct {
	Username string
	Password string
}

type RequestEntry struct {
	Id        int       `json:"id"`
	IsTicket  bool      `json:"isTicket"`
	Date      time.Time `json:"date"`
	StartTime string    `json:"startTime"`
	Duration  float32   `json:"duration"`
	Summary   string    `json:"summary"`
	Project   string    `json:"project"`
}

func UnmarshalToRequestEntries(data []byte, dateFormat string) ([]RequestEntry, error) {
	var r []RequestEntry
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalToTimeEntries(data []byte, dateFormat string) (TimeEntries, error) {
	r, err := UnmarshalToRequestEntries(data, dateFormat)

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
