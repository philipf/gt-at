package pwplugin

import (
	"fmt"
	"log"
	"regexp"

	"github.com/philipf/gt-at/at"
	"github.com/playwright-community/playwright-go"
)

// logout logs the user out of the application and waits for the authentication page to appear.
func logout(page playwright.Page) {
	log.Println("Logging out")

	// Navigate to the landing page
	_, err := page.Goto(fmt.Sprintf(at.URI_LANDING, at.BaseURL))
	if err != nil {
		log.Printf("could not goto landing page: %v\n", err)
	}

	// Wait for the landing page to fully load
	var urlRegEx = regexp.MustCompile(".*" + at.URI_LANDING_SUFFIX)

	err = page.WaitForURL(urlRegEx, playwright.PageWaitForURLOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
		Timeout:   playwright.Float(5 * 1000),
	})

	if err != nil {
		log.Printf("could not wait for landing page: %v\n", err)
	}

	log.Println("Landing page loaded")

	// Hover over the profile section to make sub-elements accessible
	err = page.Locator("[data-eii='05008GVH']").Hover()
	if err != nil {
		log.Printf("could not hover over profile: %v\n", err)
	}

	// Click the logout button in the profile section
	err = page.Locator("[data-eii='0100014V']").Click()
	if err != nil {
		log.Printf("could not click profile logout: %v\n", err)
	}

	// Locate the close button for the status out dialog
	setStatusOutLocatorCloseButton := page.Locator("div.Dialog1 div.DialogTitleBarIcon")
	setStatusOutLocatorCloseButton.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(2 * 1000),
	})
	if err != nil {
		log.Printf("could not find status out (WaitFor): %v\n", err)
	}

	// Check if the status out dialog is visible
	setStatusOut, err := setStatusOutLocatorCloseButton.IsVisible()
	if err != nil {
		log.Printf("could not find status out (IsVisible): %v\n", err)
	}

	// If the status out dialog is visible, click the close button
	if setStatusOut {
		err = setStatusOutLocatorCloseButton.Click()
		if err != nil {
			log.Printf("could not click status out: %v\n", err)
		}
	}

	// Wait for the logout process to complete and the authentication page to appear
	log.Println("Waiting for logout to complete")
	urlRegEx = regexp.MustCompile(".*Authenticate")
	err = page.WaitForURL(urlRegEx)
	if err != nil {
		log.Printf("could not wait for authentication page: %v\n", err)
	}

	log.Println("Logged out")
}
