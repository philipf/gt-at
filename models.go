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
	Exists    bool
	Submitted bool
	Error     error
}

func (te *TimeEntry) SetError(err error) {
	te.Error = err
}
