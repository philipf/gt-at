package pwplugin

import (
	"fmt"
	"log"

	"github.com/philipf/gt-at/at"
	"github.com/philipf/gt-at/pwplugin/common"
	"github.com/philipf/gt-at/pwplugin/projects"
	"github.com/philipf/gt-at/pwplugin/servicedesk"
	"github.com/playwright-community/playwright-go"
)

// NewAutoTaskPlaywright initializes and returns an instance of the autoTaskPlaywright.
func NewAutoTaskPlaywright() at.AutoTasker {
	return &autoTaskPlaywright{}
}

type autoTaskPlaywright struct{}

// CaptureTimes captures time entries in AutoTask using playwright.
func (atp *autoTaskPlaywright) CaptureTimes(entries at.TimeEntries, opts at.CaptureOptions) error {
	log.Printf("Capture entries for a total of %v time entries\n", len(entries))

	// Initialize playwright
	browser, err := common.InitPlaywright(false, opts.BrowserType, opts.Headless)
	if err != nil {
		return fmt.Errorf("could not init playwright: %v", err)
	}

	defer browser.Close()

	// Create new browser context
	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{})
	if err != nil {
		return fmt.Errorf("could not create context: %v", err)
	}

	// Open a new page in the browser
	page, err := ctx.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %v", err)
	}

	// Navigate to AutoTask
	err = gotoAutoTask(page, opts.Credentials.Username)
	if err != nil {
		return fmt.Errorf("could not goto autotask: %v", err)
	}

	// Log in to Entra
	err = loginToEntra(page, opts.Credentials.Username, opts.Credentials.Password)
	if err != nil {
		return fmt.Errorf("could not login to entra: %v", err)
	}

	// Wait for landing page after logging in (MFA might be required)
	log.Println("Login progress, MFA might be required, waiting for AT Landing Page to load")
	err = page.WaitForURL("*"+at.URI_LANDING_SUFFIX, playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(120 * 1000),
	})
	if err != nil {
		return fmt.Errorf("could not wait for url: %v", err)
	}

	log.Println("Logged in")
	at.BaseURL = at.GetBaseURL(page.URL())

	// If timesheet is already submitted, skip capturing entries
	if isSubmitted(page) {
		log.Println("Timesheet already submitted, skipping")
		return nil
	}

	// Capture the entries and then log out
	captureEntries(entries, opts.DryRun, page, opts.UserDisplayName, opts.DateFormat, opts.DayFormat)
	logout(page)

	log.Println("End of CaptureTimes")

	return nil
}

// isSubmitted checks if the timesheet is already submitted.
func isSubmitted(page playwright.Page) bool {
	count, err := page.GetByText("Recall (Un-submit)").Count()

	if err != nil {
		log.Printf("could not get count: %v\n", err)
		return false
	}

	isVisble, err := page.GetByText("Recall (Un-submit)").IsVisible()

	if err != nil {
		log.Printf("could not visible: %v\n", err)
		return false
	}

	if count == 0 || !isVisble {
		return false
	}

	return true
}

// gotoAutoTask navigates the browser to the AutoTask URI.
func gotoAutoTask(page playwright.Page, username string) error {
	_, err := page.Goto(at.URI_AUTOTASK)
	if err != nil {
		return err
	}

	err = page.WaitForURL("*Authentication.mvc*")
	if err != nil {
		return err
	}

	// Fill in the username and proceed
	usernameInput := page.GetByRole("textbox")
	usernameInput.Fill(username)
	usernameInput.Press("Enter")

	return nil
}

// captureEntries handles the capturing of both tickets and tasks.
func captureEntries(entries at.TimeEntries,
	dryRun bool,
	page playwright.Page,
	userDisplayName, dateFormat, dayFormat string) {
	tickets, tasks := entries.SplitEntries()

	// Only proceed if it's not a dry run
	if !dryRun {
		err := servicedesk.Capture(page, userDisplayName, tickets, dateFormat)
		if err != nil {
			log.Printf("could not capture tickets: %v\n", err)
		}

		err = projects.Capture(page, userDisplayName, tasks, dateFormat, dayFormat)
		if err != nil {
			log.Printf("could not capture tasks: %v\n", err)
		}
	} else {
		log.Println("Dry run, skipping")
	}
}
