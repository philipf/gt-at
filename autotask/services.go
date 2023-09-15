package autotask

const (
	URI_BASE          = "https://www.autotask.net/"
	URI_LOGIN         = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=24a25cde-1336-48f9-8bc7-84a76079487d&response_type=code&scope=openid email profile&redirect_uri=https://ww6.autotask.net/Mvc/Framework/SingleSignOn.mvc/OAuth2Callback&state=ci%3d12714%26returnUrl%3d"
	URI_TICKET_DETAIL = "https://ww6.autotask.net/Mvc/ServiceDesk/TicketDetail.mvc?ticketID=%d"
	URI_TASK_DETAIL   = "https://ww6.autotask.net/Mvc/Projects/TaskDetail.mvc?taskID=%d"
	URI_LANDING       = "https://ww6.autotask.net/Mvc/Framework/Navigation.mvc/Landing"
)

type AutoTasker interface {
	LogTimes(entries TimeEntries,
		credentials Credentials,
		userLongName string,
		browserType string,
		headless,
		dryRun bool) error
}
