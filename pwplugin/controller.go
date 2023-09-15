package pwplugin

import (
	"fmt"
	"log"
	"time"

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

	// loop through entries and set WeekNo
	for i, te := range entries {
		dt, err := time.Parse("2006/01/02", te.DateStr)
		if err != nil {
			return fmt.Errorf("could not parse date: %v", err)
		}
		entries[i].WeekNo = autotask.WeekNo(dt)
		entries[i].Date = dt // TODO: create factory method for TimeEntry
	}

	log.Printf("Logging entries for a total of %v time entries\n", len(entries))

	err, browser := common.InitPlaywright(false, browserType, headless)
	defer browser.Close()

	if err != nil {
		return fmt.Errorf("could not init playwright: %v", err)
	}

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		//BaseURL: &baseURL,
	})

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

	for _, te := range entries {
		log.Printf("%+v\n", te)
	}

	common.Logout(page)

	log.Println("End of logTimes")

	return nil
}
