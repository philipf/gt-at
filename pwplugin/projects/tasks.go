package projects

import (
	"fmt"
	"log"

	"github.com/philipf/gt-at/autotask"
	"github.com/philipf/gt-at/pwplugin/common"
	"github.com/playwright-community/playwright-go"
)

func LogTimeEntries(page playwright.Page, userDisplayName string, entries autotask.TimeEntries) error {
	log.Printf("Logging entries for a total of %v tasks\n", len(entries))

	taskIds := entries.DistinctIds()

	for _, id := range taskIds {
		err := logTimeEntriesByTaskId(page, id, entries, userDisplayName)
		if err != nil {
			return fmt.Errorf("logTimeEntries: could not log time entries for taskId: %v, error: %v", id, err)
		}
	}

	return nil
}

func logTimeEntriesByTaskId(page playwright.Page, taskId int, entries autotask.TimeEntries, userDisplayName string) error {
	_, err := page.Goto(fmt.Sprintf(autotask.URI_TASK_DETAIL, taskId))

	if err != nil {
		log.Fatalf("logTimeEntries: could not goto taskDetailUri: %v", err)
	}

	log.Println("Waiting for first conversation details to load")

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{})
	if err != nil {
		return fmt.Errorf("logTimeEntries: could not find details: %v", err)
	}
	log.Println("Conversations Loaded")

	// Build an array of ticket entries for a given taskId
	// Doing this to be a little more efficient and reduce the number of page loads
	entriesById := entries.ById(taskId)

	err = common.MarkExisiting(page, userDisplayName, entriesById)
	if err != nil {
		return fmt.Errorf("logTimeEntries: could not mark existing entries: %v", err)
	}

	for _, te := range entriesById {
		if te.Exists {
			continue
		}

		// This a new entry which needs to be captured.
		// If this is the first entry for the week then create a new time entry using the New Time Entry button
		// else is this not the first entry for the week then create a new time entry using thenEdit the time entry using the ConversationLocator

		peer := findWeekEntryPeer(entriesById, te)

		if peer == nil {
			err := newTimeEntry(page, te)
			if err != nil {
				te.SetError(err)
			}

		} else {
			err := editTimeEntry(page, te, peer)
			if err != nil {
				te.SetError(err)
			}
		}
	}

	log.Println("Done loading")

	return nil
}

func findWeekEntryPeer(entriesById autotask.TimeEntries, te *autotask.TimeEntry) *autotask.TimeEntry {
	weekNo := te.WeekNo

	for _, e := range entriesById {
		if e.WeekNo == weekNo {
			return te
		}
	}

	return nil

}

func newTimeEntry(page playwright.Page, te *autotask.TimeEntry) error {
	if te.IsTicket {
		return fmt.Errorf("logTimeEntry: only task time entries are supported")
	}

	if te.Exists {
		log.Printf("Skipping entry as it already exists: %+v\n", te)
		return nil
	}

	page.Locator("[data-eii='00000135']").Click() // New Time Entry button
	timeEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active")
	err := timeEntryDialog.WaitFor()
	if err != nil {
		return fmt.Errorf("logTimeEntry: could not find dialog: %v", err)
	}

	page.Locator("[data-eii='010000xs'] > input[type=text]").Fill(te.Date)                                                 // Date field
	page.Locator("[data-eii='010000xt'] > input[type=text]").Fill(te.StartTime)                                            // Start Time
	page.Locator("[data-eii='010000xu'] > input[type=text]").Fill(te.EndTime)                                              // End Time
	summaryNotes := page.Locator("[data-eii='000001GK']  > div.Content2 > div.InputWrapper2 > div.ContentEditable2.Small") // Summary Notes
	summaryNotes.Fill(te.Summary)

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

func editTimeEntry(page playwright.Page, te, peer *autotask.TimeEntry) error {
	panic("newTimeEntry: not implemented")
}
