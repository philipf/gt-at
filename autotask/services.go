package autotask

const (
	URI_AUTOTASK       = "https://www.autotask.net"
	URI_TICKET_DETAIL  = "%s/Mvc/ServiceDesk/TicketDetail.mvc?ticketID=%d"
	URI_TASK_DETAIL    = "%s/Mvc/Projects/TaskDetail.mvc?taskID=%d"
	URI_LANDING_SUFFIX = "/Mvc/Framework/Navigation.mvc/Landing" // used for waiting
	URI_LANDING        = "%s/" + URI_LANDING_SUFFIX              // used for navigating
)

var BaseURL string = URI_AUTOTASK

type AutoTasker interface {
	LogTimes(entries TimeEntries,
		credentials Credentials,
		userLongName string,
		browserType string,
		headless,
		dryRun bool) error
}
