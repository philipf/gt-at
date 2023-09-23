package pwplugin

import (
	"fmt"
	"log"

	"github.com/philipf/gt-at/at"
	"github.com/playwright-community/playwright-go"
)

func logout(page playwright.Page) {
	log.Println("Logging out")
	page.Goto(fmt.Sprintf(at.URI_LANDING, at.BaseURL))

	page.WaitForURL("*"+at.URI_LANDING_SUFFIX, playwright.PageWaitForURLOptions{
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

	setStatusOutLocatorCloseButton := page.Locator("div.Dialog1 div.DialogTitleBarIcon")
	setStatusOutLocatorCloseButton.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(2 * 1000),
	})

	if err != nil {
		log.Printf("could not find status out (WaitFor): %v\n", err)
	}

	setStatusOut, err := setStatusOutLocatorCloseButton.IsVisible()
	if err != nil {
		log.Printf("could not find status out (IsVisible): %v\n", err)
	}

	if setStatusOut {
		err = setStatusOutLocatorCloseButton.Click()
		if err != nil {
			log.Printf("could not click status out: %v\n", err)
		}
	}

	log.Println("Waiting for logout to complete")
	page.WaitForURL("*Authentication.mvc*")

	log.Println("Logged out")
}
