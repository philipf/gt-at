package servicedesk

import (
	"fmt"
	"log"
	"strings"

	"github.com/philipf/gt-at/at"
	"github.com/philipf/gt-at/pwplugin/common"
	"github.com/playwright-community/playwright-go"
)

func Capture(page playwright.Page, userDisplayName string, entries at.TimeEntries, dateFormat string) error {
	log.Printf("Capture entries for a total of %v tickets\n", len(entries))
	ticketIds := entries.DistinctIds()

	for _, ticketId := range ticketIds {
		err := captureByTicketId(page, ticketId, entries, userDisplayName, dateFormat)
		if err != nil {
			fmt.Printf("Capture: could not log time entries for ticketId: %v, error: %v\n", ticketId, err)
		}
	}

	return nil
}

func captureByTicketId(page playwright.Page, ticketId int, entries at.TimeEntries, userDisplayName, dateFormat string) error {
	_, err := page.Goto(fmt.Sprintf(at.URI_TICKET_DETAIL, at.BaseURL, ticketId))

	if err != nil {
		log.Fatalf("logTimeEntries: could not goto ticketDetailUri: %v", err)
	}

	log.Println("Waiting for conversation details to load")

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(5000),
	})

	if alertVisible, _ := page.Locator("#AlertDialog.Active").IsVisible(); alertVisible {
		page.Locator("#AlertDialogOkayButton").Click()
	}

	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			log.Println("Timeout waiting for first conversation details to load")
		} else {
			return fmt.Errorf("logTimeEntries: could not find details: %v", err)
		}
	}
	log.Println("Conversations Loaded")

	// Build an array of ticket entries for a given ticketId
	// Doing this to be a little more efficient and reduce the number of page loads
	entriesById := entries.ById(ticketId)

	err = common.MarkExisiting(page, userDisplayName, entriesById, dateFormat)
	if err != nil {
		return fmt.Errorf("logTimeEntries: could not mark existing entries: %v", err)
	}

	for _, te := range entriesById {
		err := captureEntry(page, te)
		if err != nil {
			te.SetError(err)
		}
	}

	log.Println("Done loading")

	return nil
}

func captureEntry(page playwright.Page, te *at.TimeEntry) error {
	log.Printf("Capture time entry: %+v\n", te)
	if !te.IsTicket {
		return fmt.Errorf("captureEntry: only ticket time entries are supported")
	}

	if te.Exists {
		log.Printf("Skipping entry as it already exists: %+v\n", te)
		return nil
	}

	page.Locator("[data-eii='000001Bb']").Click() // New Time Entry button
	timeEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active")
	err := timeEntryDialog.WaitFor()
	if err != nil {
		return fmt.Errorf("captureEntry: could not find dialog: %v", err)
	}

	if err := page.Locator("[data-eii='010000xs'] > input[type=text]").Fill(te.DateStr); err != nil {
		return fmt.Errorf("captureEntry: could not fill date: %v", err)
	}
	if err := page.Locator("[data-eii='010000xt'] > input[type=text]").Fill(te.StartTimeStr); err != nil {
		return fmt.Errorf("captureEntry: could not fill start time: %v", err)
	}

	inputs := page.Locator("[data-eii='000001GH'] input[type='text']") // Duration

	if err := inputs.First().Fill(te.DurationHoursStr); err != nil {
		return fmt.Errorf("captureEntry: could not fill hours: %v", err)
	}

	if err := inputs.Nth(1).Fill(te.DurationMinutesStr); err != nil {
		return fmt.Errorf("captureEntry: could not fill minutes: %v", err)
	}

	summaryNotes := page.Locator("[data-eii='000001GK']  > div.Content2 > div.InputWrapper2 > div.ContentEditable2.Small") // Summary Notes
	if err := summaryNotes.Fill(te.Summary); err != nil {
		return fmt.Errorf("captureEntry: could not fill summary: %v", err)
	}

	page.WaitForTimeout(1000) // Forced wait to allow the page to catch up

	saveButton := page.Locator("[data-eii='010000xo']") // Save button

	if err != nil {
		return fmt.Errorf("captureEntry: could not wait for dialog to close: %v", err)
	}

	err = saveButton.Click()
	if err != nil {
		return fmt.Errorf("captureEntry: could not click save button: %v", err)
	} else {
		log.Println("clicked save button")
	}

	err = timeEntryDialog.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateDetached,
	})

	if err != nil {
		return fmt.Errorf("captureEntry: could not wait for dialog to close: %v", err)
	}

	te.Submitted = true

	return nil
}
