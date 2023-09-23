package pwplugin

import "github.com/playwright-community/playwright-go"

func loginToEntra(page playwright.Page, user, password string) {
	if user != "" {
		// auto fill if available
		page.Locator("#i0116").Fill(user)    // Username
		page.Locator("#idSIButton9").Click() // Click Next button
	} else {
		page.Locator("#i0116").Focus()
	}

	if password != "" {
		// auto fill if available
		page.Locator("#i0118").Fill(password) // Password
		page.Locator("#idSIButton9").Click()  // Click Sign In button
	} else {
		page.Locator("#i0118").Focus()
	}
}
