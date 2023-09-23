package common

import (
	"fmt"
	"log"

	"github.com/playwright-community/playwright-go"
)

func InitPlaywright(install bool, useBrowserType string, headless bool) (playwright.Browser, error) {

	log.Println("Initiating playwright")

	runOpts := playwright.RunOptions{
		Browsers: []string{useBrowserType},
		Verbose:  true,
	}

	if install {
		log.Println("Installing playwright")
		err := playwright.Install(&runOpts)
		if err != nil {
			return nil, err
		}
		log.Println("Installed playwright")
	}

	pw, err := playwright.Run(&runOpts)
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %v", err)
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
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}

	return browser, nil
}
