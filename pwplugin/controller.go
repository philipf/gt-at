package pwplugin

import (
	"fmt"
	"log"
	"net/url"

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

	page, err := navigateToAutoTask(browser, creds)

	// err = common.LogIntoAutoTask(page, creds.Username, creds.Password)
	// if err != nil {
	// 	return fmt.Errorf("could not log into autotask: %v", err)
	// }

	//logEntries(entries, dryRun, page, userDisplayName)

	err = page.WaitForURL("*"+autotask.URI_LANDING, playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(120 * 1000),
	})

	newURL, _ := getBaseURL(page.URL())

	fmt.Println(newURL)

	common.Logout(page)

	log.Println("End of logTimes")

	return nil
}

func navigateToAutoTask(browser playwright.Browser, creds autotask.Credentials) (playwright.Page, error) {
	context, err := browser.NewContext()

	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// context, err := browser.NewContext(playwright.BrowserNewContextOptions{
	// 	//BaseURL: &baseURL,
	// })

	_, err = page.Goto(autotask.URI_AUTOTASK)

	if err != nil {
		return nil, err
	}

	page.WaitForURL("*Authentication.mvc*")
	usernameInput := page.GetByRole("textbox")
	usernameInput.Fill(creds.Username)
	usernameInput.Press("Enter")

	return page, nil
}

func logEntries(entries autotask.TimeEntries,
	dryRun bool,
	page playwright.Page,
	userDisplayName string) {
	tickets, tasks := entries.SplitEntries()

	if !dryRun {
		err := servicedesk.LogTimeEntries(page, userDisplayName, tickets)

		if err != nil {
			log.Printf("could not log tickets: %v\n", err)
		}

		err = projects.LogTimeEntries(page, userDisplayName, tasks)

		if err != nil {
			log.Printf("could not log tasks: %v\n", err)
		}

	} else {
		log.Println("Dry run, skipping logTimeEntry")
	}

	entries.PrintSummary()
}

func getBaseURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	baseURL := u.Scheme + "://" + u.Host
	//fmt.Println(baseURL)
	return baseURL, nil
}
