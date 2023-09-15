package projects

import (
	"fmt"
	"log"
	"strconv"

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

	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})

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

	weekGroups := entriesById.GroupByWeekNo()

	// Loop through each week group and create a new time entry for each week
	for _, weekEntries := range weekGroups {
		err = logTimeEntriesByWeek(page, weekEntries)
		if err != nil {
			return fmt.Errorf("logTimeEntries: could not log time entries for week: %v, error: %v", weekEntries[0].WeekNo, err)
		}
	}

	log.Println("Done loading")

	return nil
}

func logTimeEntriesByWeek(page playwright.Page, weekEntries autotask.TimeEntries) error {
	peer := findWeekEntryPeer(weekEntries, weekEntries[0])

	if peer == nil {
		return newWeekEntries(page, weekEntries)

	} else {
		return editWeekEntries(page, weekEntries, peer)
	}
}

func newWeekEntries(page playwright.Page, weekEntries autotask.TimeEntries) error {
	if err := page.Locator("[data-eii='00000135']").Click(); err != nil {
		return fmt.Errorf("newWeekEntries: could not click new time entry button: %v", err)
	}
	return captureWeekEntries(page, weekEntries)
}

func editWeekEntries(page playwright.Page, weekEntries autotask.TimeEntries, peer *autotask.TimeEntry) error {
	convLocator := peer.WeekPeerLocator.(playwright.Locator)

	err := convLocator.Locator("div.FooterActions div.LinkButton2").Nth(3).Click()

	if err != nil {
		return fmt.Errorf("editWeekEntries: could not click edit button: %v", err)
	}

	return captureWeekEntries(page, weekEntries)
}

func captureWeekEntries(page playwright.Page, weekEntries autotask.TimeEntries) error {
	weekEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active")
	err := weekEntryDialog.WaitFor()
	if err != nil {
		return fmt.Errorf("newWeekEntries: could not find weekEntryDialog: %v", err)
	}

	// Click Sunday's edit button
	timeEntryDialogSelector := page.Locator("div.Body > div.Scrolling > table > tbody div.Icon").First()
	// err = timeEntryDialogSelector.WaitFor()
	// if err != nil {
	// 	return fmt.Errorf("newWeekEntries: could not find Sunday's edit button: %v", err)
	// }

	err = timeEntryDialogSelector.Click()
	if err != nil {
		return fmt.Errorf("newWeekEntries: could not click Sunday's edit button: %v", err)
	}

	nextDayButton := page.Locator("[data-eii='0100014L']") // Next Day button
	err = nextDayButton.WaitFor()
	if err != nil {
		return fmt.Errorf("newWeekEntries: could not find timeEntryDialog: %v", err)
	}

	// Capture each week day's time if it exists
	sunday := autotask.SundayOfTheWeek(weekEntries[0].Date)

	for i := 0; i < 7; i++ {
		// Find the time entry for the current day
		entries := weekEntries.ByDate(sunday.AddDate(0, 0, i))

		if len(entries) > 1 {
			for _, e := range entries {
				e.SetError(fmt.Errorf("newWeekEntries: more than one entry for a given day: %v", sunday))
			}
		} else if len(entries) == 0 {
			// No time entry for this day, skip to the next day
		} else {
			te := entries[0]
			err = captureDay(page, te)
			if err != nil {
				te.SetError(fmt.Errorf("newWeekEntries: could not capture day: %v", err))
			}
		}

		if i < 6 {
			err := navigateToNextDay(page)
			if err != nil {
				return fmt.Errorf("newWeekEntries: could not navigate to next day: %v", err)
			}
		} else {
			err := saveWeek(page)
			if err != nil {
				return fmt.Errorf("newWeekEntries: could not save week: %v", err)
			}
		}
	}

	// Mark all entries as submitted
	for _, te := range weekEntries {
		if te.Error != nil {
			te.Submitted = true
		}
	}

	return nil
}

func captureDay(page playwright.Page, te *autotask.TimeEntry) error {
	err := page.Locator("[data-eii='0100014M']").Fill(strconv.FormatFloat(float64(te.Duration), 'f', -1, 32))
	if err != nil {
		return fmt.Errorf("newWeekEntries: could not fill in duration: %v", err)
	}

	summaryNotes := page.Locator("[data-eii='0100014N']  > div.Content2 > div.InputWrapper2 > div.ContentEditable2.Small")
	err = summaryNotes.Fill(te.Summary)
	if err != nil {
		return fmt.Errorf("newWeekEntries: could not fill in summary notes: %v", err)
	}

	page.WaitForTimeout(1000)

	return nil
}

func saveWeek(page playwright.Page) error {
	// click save button
	okButton := page.Locator("[data-eii='0100014J']") // OK button to save
	err := okButton.Click()
	if err != nil {
		return fmt.Errorf("saveWeek: could not click ok button: %v", err)
	}

	weekEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active").Last()

	saveAndCloseButton := page.Locator("[data-eii='010000p7']") // Save and Close button
	err = saveAndCloseButton.Click()
	if err != nil {
		return fmt.Errorf("saveWeek: could not click save and close button: %v", err)
	}

	// wait for dialog to close
	err = weekEntryDialog.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateDetached,
	})

	if err != nil {
		return fmt.Errorf("saveWeek: could not wait for dialog to close: %v", err)
	}

	return nil
}

func navigateToNextDay(page playwright.Page) error {
	nextDayButtonSelector := "[data-eii='0100014L']"

	//weekEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active").Last()

	nextDayButton := page.Locator(nextDayButtonSelector)
	err := nextDayButton.Click()
	if err != nil {
		return fmt.Errorf("navigateToNextDay: could not click next day button: %v", err)
	}

	// err = weekEntryDialog.WaitFor(playwright.LocatorWaitForOptions{
	// 	State: playwright.WaitForSelectorStateDetached,
	// })

	// if err != nil {
	// 	return fmt.Errorf("navigateToNextDay: could not find weekEntryDialog: %v", err)
	// }

	err = page.Locator(nextDayButtonSelector).WaitFor()
	if err != nil {
		return fmt.Errorf("navigateToNextDay: could not find next day button: %v", err)
	}

	return nil

}

func findWeekEntryPeer(entriesById autotask.TimeEntries, te *autotask.TimeEntry) *autotask.TimeEntry {
	weekNo := te.WeekNo

	for _, e := range entriesById {
		if e.WeekNo == weekNo && e.Exists {
			return te
		}
	}

	return nil
}
