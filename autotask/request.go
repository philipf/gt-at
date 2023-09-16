package autotask

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
}

func UnmarshalToRequestEntries(data []byte) ([]RequestEntry, error) {
	var r []RequestEntry
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalToTimeEntries(data []byte) (TimeEntries, error) {
	r, err := UnmarshalToRequestEntries(data)

	if err != nil {
		return nil, err
	}

	// create a slice of TimeEntry from the RequestEntry slice
	var entries TimeEntries

	for _, e := range r {
		te := NewEntry(e.Id, e.IsTicket, e.Date, e.StartTime, e.Duration, e.Summary)
		entries = append(entries, te)
	}

	return entries, nil
}

type LoadOptions struct {
	Credentials     Credentials
	DryRun          bool
	UserDisplayName string
}
