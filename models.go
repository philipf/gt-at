package autotask

type Credentials struct {
	Username string
	Password string
}

type TimeEntry struct {
	TicketId  int
	Date      string
	StartTime string
	EndTime   string
	Summary   string
}
