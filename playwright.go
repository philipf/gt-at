package autotask

import (
	"fmt"
	"log"

	"github.com/playwright-community/playwright-go"
)

func NewAutoTaskPlaywright() AutoTasker {
	return &autoTaskPlaywright{}
}

type autoTaskPlaywright struct{}

func (atp *autoTaskPlaywright) LogTimes(
	entries []TimeEntry,
	creds Credentials,
	browserType string,
	headless, dryRun bool) error {

	err, page, browser := initPlaywright(false, browserType, headless)
	defer browser.Close()

	if err != nil {
		return fmt.Errorf("could not init playwright: %v", err)
	}

	err = logIntoAutoTask(page, creds.Username, creds.Password)
	if err != nil {
		return fmt.Errorf("could not log into autotask: %v", err)
	}

	for _, te := range entries {
		if dryRun {
			log.Printf("dry run: %+v\n", te)
			return nil
		}

		log.Printf("Capture: %+v\n", te)
		err = logTimeEntry(page, te)

		if err != nil {
			return fmt.Errorf("could not log time entry: %+v, Reason: %v", te, err)
		}
	}

	return nil
}

func logIntoAutoTask(page playwright.Page, user, password string) error {
	_, err := page.Goto(URI_LOGIN)
	if err != nil {
		return fmt.Errorf("could not goto: %v", err)
	}

	if user != "" {
		// auto fill if available
		page.Locator("#i0116").Fill(user)    // Username
		page.Locator("#idSIButton9").Click() // Click Next button
	}

	if password != "" {
		// auto fill if available
		page.Locator("#i0118").Fill(password) // Password
		page.Locator("#idSIButton9").Click()  // Click Sign In button
	}

	log.Println("waiting for log in to complete")

	err = page.WaitForURL(URI_LANDING, createWaitForeverOpts())
	if err != nil {
		return fmt.Errorf("could not wait for url: %v", err)
	}

	log.Println("logged in")

	return nil
}

func logTimeEntry(page playwright.Page, te TimeEntry) error {

	_, err := page.Goto(fmt.Sprintf(URI_TICKET_DETAIL, te.TicketId))
	if err != nil {
		log.Fatalf("could not goto ticketDetailUri: %v", err)
	}

	page.Locator("[data-eii='000001Bb']").Click()                                                                          // New Time Entry button
	page.Locator("[data-eii='010000xs'] > input[type=text]").Fill(te.Date)                                                 // Date field
	page.Locator("[data-eii='010000xt'] > input[type=text]").Fill(te.StartTime)                                            // Start Time
	page.Locator("[data-eii='010000xu'] > input[type=text]").Fill(te.EndTime)                                              // End Time
	summaryNotes := page.Locator("[data-eii='000001GK']  > div.Content2 > div.InputWrapper2 > div.ContentEditable2.Small") // Summary Notes
	summaryNotes.Fill(te.Summary)
	summaryNotes.Type("\t") // required to get it to save 1/2 :(

	saveButton := page.Locator("[data-eii='010000xo']") // Save button
	saveButton.Hover()                                  // required to get it to save 2/2 :(

	// page.OnResponse(func(response playwright.Response) {
	// 	log.Printf("Response: %v\n", response)
	// })

	err = saveButton.Click()
	if err != nil {
		return fmt.Errorf("could not click save button: %v", err)
	} else {
		log.Println("clicked save button")
	}

	page.WaitForURL("**TicketDetail.mvc?workspace=False**")

	log.Println("waiting for save to complete")

	//bg := page.Locator("#BackgroundOverlay.Active")
	//page.WaitForEvent()

	//err = page.WaitForLoadState()
	// if bg != nil {
	// 	return fmt.Errorf("could not find bg")
	// }

	// waitOpts := playwright.LocatorWaitForOptions{
	// 	State: playwright.WaitForSelectorStateHidden,
	// }
	// err = bg.WaitFor(waitOpts)
	// if err != nil {
	// 	return fmt.Errorf("could not wait for bg: %v", err)
	// }

	log.Println("Done loading")

	// if err := page.Locator("#BackgroundOverlay:not(.Active)").WaitFor(); err != nil {
	// 	return fmt.Errorf("error waiting for #mydiv to disappear: %v", err)
	// }

	// log.Println("waiting for save to complete")
	// time.Sleep(10 * time.Second)
	// log.Println("save complete")

	return nil
}

func initPlaywright(install bool, useBrowserType string, headless bool) (error, playwright.Page, playwright.Browser) {

	log.Println("Initiating playwright")

	runOpts := playwright.RunOptions{
		Browsers: []string{useBrowserType},
		Verbose:  true,
	}

	if install {
		log.Println("Installing playwright")
		err := playwright.Install(&runOpts)
		if err != nil {
			return err, nil, nil
		}
		log.Println("Installed playwright")
	}

	pw, err := playwright.Run(&runOpts)
	if err != nil {
		return fmt.Errorf("could not start playwright: %v", err), nil, nil
	}

	browserOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	}

	var browserType playwright.BrowserType

	if useBrowserType == "chromium" {
		browserType = pw.Chromium
	} else if useBrowserType == "firefox" {
		browserType = pw.Firefox
	} else {
		browserType = pw.WebKit
	}

	browser, err := browserType.Launch(browserOpts)

	if err != nil {
		return fmt.Errorf("could not launch browser: %v", err), nil, nil
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	return err, page, browser
}

func createWaitForeverOpts() playwright.PageWaitForURLOptions {
	timeout := 0.0
	waitOpts := playwright.PageWaitForURLOptions{
		Timeout: &timeout,
	}
	return waitOpts
}
