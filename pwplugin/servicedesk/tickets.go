package servicedesk

import (
	"fmt"
	"log"

	"github.com/philipf/gt-at/autotask"
	"github.com/philipf/gt-at/pwplugin/common"
	"github.com/playwright-community/playwright-go"
)

func LogTimeEntries(page playwright.Page, userDisplayName string, entries autotask.TimeEntries) error {
	log.Printf("Logging entries for a total of %v tickets\n", len(entries))
	ticketIds := entries.DistinctIds()

	for _, ticketId := range ticketIds {
		err := logTimeEntriesByTicketId(page, ticketId, entries, userDisplayName)
		if err != nil {
			return fmt.Errorf("logTimeEntries: could not log time entries for ticketId: %v, error: %v", ticketId, err)
		}
	}

	return nil
}

func logTimeEntriesByTicketId(page playwright.Page, ticketId int, entries autotask.TimeEntries, userDisplayName string) error {
	_, err := page.Goto(fmt.Sprintf(autotask.URI_TICKET_DETAIL, autotask.BaseURL, ticketId))

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
		if playwright.TimeoutError.Is(err) {
			log.Println("Timeout waiting for first conversation details to load")
		} else {
			return fmt.Errorf("logTimeEntries: could not find details: %v", err)
		}
	}
	log.Println("Conversations Loaded")

	// Build an array of ticket entries for a given ticketId
	// Doing this to be a little more efficient and reduce the number of page loads
	entriesById := entries.ById(ticketId)

	err = common.MarkExisiting(page, userDisplayName, entriesById)
	if err != nil {
		return fmt.Errorf("logTimeEntries: could not mark existing entries: %v", err)
	}

	for _, te := range entriesById {
		err := logTimeEntry(page, te)
		if err != nil {
			te.SetError(err)
		}
	}

	log.Println("Done loading")

	return nil
}

func logTimeEntry(page playwright.Page, te *autotask.TimeEntry) error {
	log.Printf("Logging time entry: %+v\n", te)
	if !te.IsTicket {
		return fmt.Errorf("logTimeEntry: only ticket time entries are supported")
	}

	if te.Exists {
		log.Printf("Skipping entry as it already exists: %+v\n", te)
		return nil
	}

	page.Locator("[data-eii='000001Bb']").Click() // New Time Entry button
	timeEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active")
	err := timeEntryDialog.WaitFor()
	if err != nil {
		return fmt.Errorf("logTimeEntry: could not find dialog: %v", err)
	}

	if err := page.Locator("[data-eii='010000xs'] > input[type=text]").Fill(te.DateStr); err != nil {
		return fmt.Errorf("logTimeEntry: could not fill date: %v", err)
	}
	if err := page.Locator("[data-eii='010000xt'] > input[type=text]").Fill(te.StartTimeStr); err != nil {
		return fmt.Errorf("logTimeEntry: could not fill start time: %v", err)
	}
	// if err := page.Locator("[data-eii='010000xu'] > input[type=text]").Fill(te.EndTimeStr); err != nil {
	// 	return fmt.Errorf("logTimeEntry: could not fill end time: %v", err)
	//}

	inputs := page.Locator("[data-eii='000001GH'] input[type='text']") // Duration

	if err := inputs.First().Fill(te.DurationHoursStr); err != nil {
		return fmt.Errorf("logTimeEntry: could not fill hours: %v", err)
	}

	if err := inputs.Nth(1).Fill(te.DurationMinutesStr); err != nil {
		return fmt.Errorf("logTimeEntry: could not fill minutes: %v", err)
	}

	summaryNotes := page.Locator("[data-eii='000001GK']  > div.Content2 > div.InputWrapper2 > div.ContentEditable2.Small") // Summary Notes
	if err := summaryNotes.Fill(te.Summary); err != nil {
		return fmt.Errorf("logTimeEntry: could not fill summary: %v", err)
	}

	page.WaitForTimeout(1000)

	saveButton := page.Locator("[data-eii='010000xo']") // Save button

	if err != nil {
		return fmt.Errorf("logTimeEntry: could not wait for dialog to close: %v", err)
	}

	err = saveButton.Click()
	if err != nil {
		return fmt.Errorf("logTimeEntry: could not click save button: %v", err)
	} else {
		log.Println("clicked save button")
	}

	err = timeEntryDialog.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateDetached,
	})

	if err != nil {
		return fmt.Errorf("logTimeEntry: could not wait for dialog to close: %v", err)
	}

	te.Submitted = true

	return nil
}
