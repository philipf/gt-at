package pwplugin

import (
	"fmt"
	"log"

	"github.com/olekukonko/tablewriter"
	"github.com/philipf/gt-at/autotask"
	"github.com/philipf/gt-at/pwplugin/common"
	"github.com/philipf/gt-at/pwplugin/projects"
	"github.com/philipf/gt-at/pwplugin/servicedesk"
	"github.com/playwright-community/playwright-go"
)

func NewAutoTaskPlaywright() autotask.AutoTasker {
	return &autoTaskPlaywright{}
}

type autoTaskPlaywright struct{}

func (atp *autoTaskPlaywright) LogTimes(
	//baseURL string,
	entries autotask.TimeEntries,
	creds autotask.Credentials,
	userDisplayName string,
	browserType string,
	headless, dryRun bool) error {

	log.Printf("Logging entries for a total of %v time entries\n", len(entries))

	err, browser := common.InitPlaywright(false, browserType, headless)
	defer browser.Close()

	if err != nil {
		return fmt.Errorf("could not init playwright: %v", err)
	}

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		//BaseURL: &baseURL,
	})

	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	err = common.LogIntoAutoTask(page, creds.Username, creds.Password)
	if err != nil {
		return fmt.Errorf("could not log into autotask: %v", err)
	}

	tickets, tasks := entries.SplitEntries()

	if !dryRun {
		err = servicedesk.LogTimeEntries(page, userDisplayName, tickets)

		if err != nil {
			return fmt.Errorf("could not log tickets: %v", err)
		}

		err = projects.LogTimeEntries(page, userDisplayName, tasks)

		if err != nil {
			return fmt.Errorf("could not log tasks: %v", err)
		}

	} else {
		log.Println("Dry run, skipping logTimeEntry")
	}

	prettyPrint(entries)

	common.Logout(page)

	log.Println("End of logTimes")

	return nil
}

// receiver function to convert bool to Y/N
func ToYN(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}

func prettyPrint(entries autotask.TimeEntries) {
	table := tablewriter.NewWriter(log.Writer())
	table.SetHeader([]string{"Id", "IsTicket", "Date", "Start", "Time", "Exists", "Saved", "Error"})

	for _, entry := range entries {
		var errMsg string
		if entry.Error != nil {
			errMsg = "Y"
		} else {
			errMsg = ""
		}

		row := []string{
			fmt.Sprintf("%d", entry.Id),
			ToYN(entry.IsTicket),
			entry.DateStr,
			entry.StartTimeStr,
			fmt.Sprintf("%.2f", entry.Duration),
			ToYN(entry.Exists),
			ToYN(entry.Submitted),
			errMsg,
		}
		table.Append(row)
	}

	table.Render() // Sends output to stdout
}
