package common

import (
	"fmt"
	"log"

	"github.com/philipf/gt-at/autotask"
	"github.com/playwright-community/playwright-go"
)

func LogIntoAutoTask(page playwright.Page, user, password string) error {
	_, err := page.Goto("")
	// _, err := page.Goto(autotask.URI_LOGIN)
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

	err = page.WaitForURL(autotask.URI_LANDING, playwright.PageWaitForURLOptions{
		Timeout: playwright.Float(120 * 1000),
	})

	if err != nil {
		return fmt.Errorf("could not wait for url: %v", err)
	}

	log.Println("logged in")

	return nil
}

func Logout(page playwright.Page) {
	log.Println("Logging out")
	page.Locator("data-eii='0100014V'").Click() // profile logout
	//page.GetByText("No, leave my Status as In").WaitFor()
	page.Locator("data-eii='03000043'").WaitFor() // radio button" No, leave my Status as In
	page.Locator("data-eii='03000043'").Click()
	page.Locator("data-eii='05008CnG'").Click() // Ok button on logout
	log.Println("Logged out")
}
