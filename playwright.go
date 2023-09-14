package autotask

import (
	"fmt"
	"log"
	"strings"

	"github.com/playwright-community/playwright-go"
)

func NewAutoTaskPlaywright() AutoTasker {
	return &autoTaskPlaywright{}
}

type autoTaskPlaywright struct{}

func (atp *autoTaskPlaywright) LogTimes(
	//baseURL string,
	entries []*TimeEntry,
	creds Credentials,
	userLongName string,
	browserType string,
	headless, dryRun bool) error {

	err, browser := initPlaywright(false, browserType, headless)
	defer browser.Close()

	if err != nil {
		return fmt.Errorf("could not init playwright: %v", err)
	}

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		//BaseURL: &baseURL,
	})

	page, err := context.NewPage()
	//	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	err = logIntoAutoTask(page, creds.Username, creds.Password)
	if err != nil {
		return fmt.Errorf("could not log into autotask: %v", err)
	}

	// TODO, group entries by ticket id
	if !dryRun {
		err = logTimeEntry(page, userLongName, entries)
	} else {
		log.Println("Dry run, skipping logTimeEntry")
	}

	// Prettry print entries
	for _, te := range entries {
		log.Printf("%+v\n", te)
	}

	if err != nil {
		return fmt.Errorf("could not log time entry: %+v, Reason: %v", entries, err)
	}

	//time.Sleep(60 * time.Second)

	log.Println("End of logTimes")

	return nil
}

func entryExists(page playwright.Page, userLongName string, te *TimeEntry) (bool, error) {
	detailsSelector := page.Locator("div > .ConversationChunk > .ConversationItem .Details")
	convs, err := detailsSelector.All()

	if err != nil {
		return false, fmt.Errorf("could not find conversations: %v", err)
	}

	log.Printf("Found %v conversations\n", len(convs))

	for _, conv := range convs {
		author := conv.Locator("div > .Author div.Text2")

		authorName, err := author.TextContent()

		if err != nil {
			return false, fmt.Errorf("could not find author TextContent: %+v", err)
		}

		if authorName == userLongName {
			//log.Printf("Found author: %v\n", author)

			//re, _ := regexp.Compile("2023/09/11")

			//t := author.GetByText("2023/09/11 10:30 - 11:00 (0.5000 hours)")
			//t := author.GetByText(re)
			//t := conv.Filter().GetByText(te.Date)
			//t := author.GetByText(te.Date)

			// html, err := conv.InnerHTML()
			// log.Println(html)

			timeDetail := conv.Locator("div.Title div.Text > span")

			t, err := timeDetail.TextContent()

			if err != nil {
				return false, fmt.Errorf("could not find timeDetail TextContent: %+v", err)
			}

			if strings.HasPrefix(t, te.Date) {
				log.Printf("Found date: %v\n", te.Date)
				return true, nil
			}
		}
	}

	return false, nil
}

func logIntoAutoTask(page playwright.Page, user, password string) error {
	_, err := page.Goto(URI_LOGIN)
	if err != nil {
		return fmt.Errorf("could not goto: %v", err)
	}

	//r.

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

func logTimeEntry(page playwright.Page, userLongName string, entries []*TimeEntry) error {
	_, err := page.Goto(fmt.Sprintf(URI_TICKET_DETAIL, entries[0].TicketId))

	if err != nil {
		log.Fatalf("could not goto ticketDetailUri: %v", err)
	}

	log.Println("Waiting for first conversation details to load")

	loadOpts := playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
		//Timeout: playwright.Float(2000),
	}

	err = page.WaitForLoadState(loadOpts)
	log.Println("done loading")

	if err != nil {
		return fmt.Errorf("could not find details: %v", err)
	}

	for _, te := range entries {
		exists, err := entryExists(page, userLongName, te)

		if err != nil {
			te.SetError(fmt.Errorf("could not check if entry exists: %v", err))
			continue
		}

		te.Exists = exists
	}

	for _, te := range entries {
		if te.Exists {
			log.Printf("Skipping entry as it already exists: %+v\n", te)
			continue
		}

		if te.Error != nil {
			continue
		}

		page.Locator("[data-eii='000001Bb']").Click() // New Time Entry button
		timeEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active")
		err := timeEntryDialog.WaitFor()
		if err != nil {
			te.SetError(fmt.Errorf("could not find dialog: %v", err))
			continue
		}

		page.Locator("[data-eii='010000xs'] > input[type=text]").Fill(te.Date)                                                 // Date field
		page.Locator("[data-eii='010000xt'] > input[type=text]").Fill(te.StartTime)                                            // Start Time
		page.Locator("[data-eii='010000xu'] > input[type=text]").Fill(te.EndTime)                                              // End Time
		summaryNotes := page.Locator("[data-eii='000001GK']  > div.Content2 > div.InputWrapper2 > div.ContentEditable2.Small") // Summary Notes
		summaryNotes.Fill(te.Summary)

		page.WaitForTimeout(1000)

		saveButton := page.Locator("[data-eii='010000xo']") // Save button

		if err != nil {
			te.SetError(fmt.Errorf("could not wait for dialog to close: %v", err))
			continue
		}

		err = saveButton.Click()
		if err != nil {
			te.SetError(fmt.Errorf("could not click save button: %v", err))
			continue
		} else {
			log.Println("clicked save button")
		}

		err = timeEntryDialog.WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateDetached,
		})

		//page.WaitForURL("**TicketDetail.mvc?workspace=False**")

		te.Submitted = true
		log.Println("Done loading")
	}

	return nil
}

func initPlaywright(install bool, useBrowserType string, headless bool) (error, playwright.Browser) {

	log.Println("Initiating playwright")

	runOpts := playwright.RunOptions{
		Browsers: []string{useBrowserType},
		Verbose:  true,
	}

	if install {
		log.Println("Installing playwright")
		err := playwright.Install(&runOpts)
		if err != nil {
			return err, nil
		}
		log.Println("Installed playwright")
	}

	pw, err := playwright.Run(&runOpts)
	if err != nil {
		return fmt.Errorf("could not start playwright: %v", err), nil
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
		return fmt.Errorf("could not launch browser: %v", err), nil
	}

	return err, browser
}

func createWaitForeverOpts() playwright.PageWaitForURLOptions {
	timeout := 0.0
	waitOpts := playwright.PageWaitForURLOptions{
		Timeout: &timeout,
	}
	return waitOpts
}
