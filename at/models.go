package at

import (
	"sort"
	"strconv"
	"time"
)

type TimeEntry struct {
	Id           int
	IsTicket     bool // if not a ticket, it's a task
	Date         time.Time
	DateStr      string
	StartTimeStr string
	Duration     float32 // in hours
	Summary      string
	Project      string

	// Derived properties
	Exists             bool
	Submitted          bool
	Error              error
	DurationHours      int
	DurationMinutes    float32
	DurationHoursStr   string
	DurationMinutesStr string
	WeekNo             int
	WeekPeerLocator    interface{}
}

// NewEntry constructs a TimeEntry and calculates its derived properties
func NewEntry(id int,
	isTicket bool,
	date time.Time,
	startTimeStr string,
	duration float32,
	summary string,
	project string,
	dateFormat string) *TimeEntry {

	e := &TimeEntry{
		Id:           id,
		IsTicket:     isTicket,
		Date:         date,
		DateStr:      date.Format(dateFormat),
		StartTimeStr: startTimeStr,
		Duration:     duration,
		Summary:      summary,
		Project:      project,
	}

	e.calculateDerived()

	return e
}

// calculateDerived computes the derived properties of the TimeEntry
func (te *TimeEntry) calculateDerived() {
	te.DurationHours = int(te.Duration)
	te.DurationMinutes = (te.Duration - float32(te.DurationHours)) * 60
	te.DurationMinutes = float32(int(te.DurationMinutes + 0.5)) // Round DurationMinutes to nearest minute
	te.DurationHoursStr = strconv.Itoa(te.DurationHours)
	te.DurationMinutesStr = strconv.Itoa(int(te.DurationMinutes))

	te.WeekNo = WeekNo(te.Date)
}

// SetError sets an error for the TimeEntry
func (te *TimeEntry) SetError(err error) {
	te.Error = err
}

type TimeEntries []*TimeEntry

func (t TimeEntries) Len() int           { return len(t) }
func (t TimeEntries) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TimeEntries) Less(i, j int) bool { return t[i].Date.Before(t[j].Date) }

// SortByDate sorts the entries by their Date
func (t TimeEntries) SortByDate() {
	sort.Sort(t)
}

// DistinctIds returns a list of distinct Ids from the TimeEntries
func (a TimeEntries) DistinctIds() []int {
	seen := make(map[int]bool)
	var result []int

	for _, entry := range a {
		if !seen[entry.Id] {
			seen[entry.Id] = true
			result = append(result, entry.Id)
		}
	}

	return result
}

// DistinctWeekNos returns a slice of distinct week numbers
func (a TimeEntries) DistinctWeekNos() []int {
	seen := make(map[int]bool)
	var result []int

	for _, entry := range a {
		if !seen[entry.WeekNo] {
			seen[entry.WeekNo] = true
			result = append(result, entry.WeekNo)
		}
	}

	return result
}

// SplitEntries splits the TimeEntries into two lists based on their IsTicket flag, the first list contains tickets, the second contains tasks
func (a TimeEntries) SplitEntries() (TimeEntries, TimeEntries) {
	tickets := make(TimeEntries, 0)
	tasks := make(TimeEntries, 0)

	for _, entry := range a {
		if entry.IsTicket {
			tickets = append(tickets, entry)
		} else {
			tasks = append(tasks, entry)
		}
	}

	return tickets, tasks
}

// ById retrieves entries based on their Id
func (a TimeEntries) ById(id int) TimeEntries {
	entries := make(TimeEntries, 0)

	for _, entry := range a {
		if entry.Id == id {
			entries = append(entries, entry)
		}
	}

	return entries
}

// GroupByWeekNo groups the TimeEntries based on their week number
func (a TimeEntries) GroupByWeekNo() map[int]TimeEntries {
	groups := make(map[int]TimeEntries)

	for _, entry := range a {
		groups[entry.WeekNo] = append(groups[entry.WeekNo], entry)
	}

	return groups
}

// ByDate retrieves all entries that match a given date
func (entries TimeEntries) ByDate(date time.Time) TimeEntries {
	result := make(TimeEntries, 0)

	for _, entry := range entries {
		if entry.Date == date {
			result = append(result, entry)
		}
	}

	return result
}
