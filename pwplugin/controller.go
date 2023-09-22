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

	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{
		//BaseURL: autotask.BaseURL,
	})

	if err != nil {
		return fmt.Errorf("could not create context: %v", err)
	}

	page, err := ctx.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	err = gotoAutoTask(page, creds.Username)
	if err != nil {
		return fmt.Errorf("could not goto autotask: %v", err)
	}

	loginToEntra(page, creds.Username, creds.Password)

	log.Println("Log in progress, MFA might be required, waiting for AT Landing Page to load")

	err = page.WaitForURL("*"+autotask.URI_LANDING_SUFFIX, playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(120 * 1000),
	})

	if err != nil {
		return fmt.Errorf("could not wait for url: %v", err)
	}

	log.Println("logged in")

	autotask.BaseURL = getBaseURL(page.URL())

	logEntries(entries, dryRun, page, userDisplayName)
	entries.PrintSummary()

	logout(page)

	log.Println("End of logTimes")

	return nil
}

func gotoAutoTask(page playwright.Page, username string) error {
	_, err := page.Goto(autotask.URI_AUTOTASK)

	if err != nil {
		return err
	}

	err = page.WaitForURL("*Authentication.mvc*")
	if err != nil {
		return err
	}

	usernameInput := page.GetByRole("textbox")
	usernameInput.Fill(username)
	usernameInput.Press("Enter")

	return nil
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
}

func getBaseURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("could not parse url: %v\n", err)
		return ""
	}

	baseURL := u.Scheme + "://" + u.Host
	return baseURL
}

func loginToEntra(page playwright.Page, user, password string) {
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
}

func logout(page playwright.Page) {
	log.Println("Logging out")
	page.Goto(fmt.Sprintf(autotask.URI_LANDING, autotask.BaseURL))

	page.WaitForURL("*"+autotask.URI_LANDING_SUFFIX, playwright.PageWaitForURLOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
		Timeout:   playwright.Float(5 * 1000),
	})

	log.Println("Landing page loaded")

	err := page.Locator("[data-eii='05008GVH']").Hover() //hover to enable elements below
	if err != nil {
		log.Printf("could not hover over profile: %v\n", err)
	}

	err = page.Locator("[data-eii='0100014V']").Click() // profile logout
	if err != nil {
		log.Printf("could not click profile logout: %v\n", err)
	}

	err = page.Locator("[data-eii='03000043']").WaitFor() // radio button" No, leave my Status as In
	if err != nil {
		log.Printf("could not find radio button: %v\n", err)
	}

	err = page.Locator("[data-eii='03000043']").Click()
	if err != nil {
		log.Printf("could not click radio button: %v\n", err)
	}

	err = page.Locator("[data-eii='05008CnG']").Click() // Ok button on logout
	if err != nil {
		log.Printf("could not click logout button: %v\n", err)
	}

	log.Println("Waiting for logout to complete")
	page.WaitForURL("*Authentication.mvc*")

	log.Println("Logged out")
}
