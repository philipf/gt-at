package at

// Constants for specific AutoTask URIs.
const (
	// Base URL for AutoTask.
	URI_AUTOTASK = "https://www.autotask.net"

	// Format string for ticket detail URL, expects the base URL and ticketID.
	URI_TICKET_DETAIL = "%s/Mvc/ServiceDesk/TicketDetail.mvc?ticketID=%d"

	// Format string for task detail URL, expects the base URL and taskID.
	URI_TASK_DETAIL = "%s/Mvc/Projects/TaskDetail.mvc?taskID=%d"

	// Suffix for the landing URL, used for waiting.
	URI_LANDING_SUFFIX = "/Mvc/Framework/Navigation.mvc/Landing"

	// Format string for the landing URL, expects the base URL.
	URI_LANDING = "%s/" + URI_LANDING_SUFFIX
)

// BaseURL is the default base URL for AutoTask operations.
var BaseURL string = URI_AUTOTASK

// CaptureOptions defines the options for the CaptureTimes method.
type CaptureOptions struct {
	Credentials     Credentials // Authentication details.
	DryRun          bool        // If true, does a dry run without actual capture.
	UserDisplayName string      // Display name of the user in AutoTask, this available under the user profile. This value is used to find time entries for the user.
	BrowserType     string      // Type of the browser to use, e.g., "chromium", "firefox" and "webkit".
	Headless        bool        // If true, browser operates in headless mode.
	DateFormat      string      // Format for date representation.
	DayFormat       string      // Format for day representation.
}

// AutoTasker is an interface for capturing time entries.
type AutoTasker interface {
	// CaptureTimes captures time entries based on the provided options.
	CaptureTimes(entries TimeEntries, opts CaptureOptions) error
}
