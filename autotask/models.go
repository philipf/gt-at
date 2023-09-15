package autotask

import "sort"

type Credentials struct {
	Username string
	Password string
}

type TimeEntry struct {
	Id              int
	IsTicket        bool // if not a ticket, it's a task
	Date            string
	StartTime       string
	EndTime         string
	Summary         string
	Exists          bool
	Submitted       bool
	Error           error
	Duration        float32 // in hours
	WeekNo          int
	WeekPeerLocator interface{}
}

func (te *TimeEntry) SetError(err error) {
	te.Error = err
}

type TimeEntries []*TimeEntry

func (t TimeEntries) Len() int           { return len(t) }
func (t TimeEntries) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TimeEntries) Less(i, j int) bool { return t[i].Id < t[j].Id }

func (t TimeEntries) SortById() {
	sort.Sort(t)
}

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

// Split into two lists, one for tickets and one for tasks
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

// Get entries by Id
func (a TimeEntries) ById(id int) TimeEntries {
	entries := make(TimeEntries, 0)

	for _, entry := range a {
		if entry.Id == id {
			entries = append(entries, entry)
		}
	}

	return entries
}
