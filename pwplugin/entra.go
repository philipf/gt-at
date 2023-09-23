package pwplugin

import "github.com/playwright-community/playwright-go"

// loginToEntra automates the login process for the Entra application.
// It takes a playwright page, user, and password as arguments.
func loginToEntra(page playwright.Page, user, password string) error {

	// Check if the user string is provided
	if user != "" {
		// If available, autofill the username and click the Next button
		err := page.Locator("#i0116").Fill(user) // Fill in the Username
		if err != nil {
			return err
		}

		err = page.Locator("#idSIButton9").Click() // Click the Next button
		if err != nil {
			return err
		}
	} else {
		// If not provided, focus on the username input for manual entry
		err := page.Locator("#i0116").Focus()
		if err != nil {
			return err
		}
	}

	// Check if the password string is provided
	if password != "" {
		// If available, autofill the password and click the Sign In button
		err := page.Locator("#i0118").Fill(password) // Fill in the Password
		if err != nil {
			return err
		}
		err = page.Locator("#idSIButton9").Click() // Click the Sign In button
		if err != nil {
			return err
		}
	} else {
		// If not provided, focus on the password input for manual entry
		err := page.Locator("#i0118").Focus()
		if err != nil {
			return err
		}
	}

	return nil
}
