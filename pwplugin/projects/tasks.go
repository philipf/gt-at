package projects

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/philipf/gt-at/at"
	"github.com/philipf/gt-at/pwplugin/common"
	"github.com/playwright-community/playwright-go"
)

func Capture(page playwright.Page, userDisplayName string, entries at.TimeEntries, dateFormat, dayFormat string) error {
	log.Printf("Capture entries for a total of %v tasks\n", len(entries))

	taskIds := entries.DistinctIds()

	for _, id := range taskIds {
		err := captureByTaskId(page, id, entries, userDisplayName, dateFormat, dayFormat)
		if err != nil {
			fmt.Printf("Capture: could not log time entries for taskId: %v, error: %v\n", id, err)
		}
	}

	return nil
}

func captureByTaskId(page playwright.Page, taskId int, entries at.TimeEntries, userDisplayName, dateFormat, dayFormat string) error {
	_, err := page.Goto(fmt.Sprintf(at.URI_TASK_DETAIL, at.BaseURL, taskId))

	if err != nil {
		log.Fatalf("captureByTaskId: could not goto taskDetailUri: %v", err)
	}

	log.Println("Waiting for first conversation details to load")

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
			return fmt.Errorf("captureByTaskId: could not find details: %v", err)
		}
	}

	log.Println("Conversations Loaded")

	// Build an array of ticket entries for a given taskId
	// Doing this to be a little more efficient and reduce the number of page loads
	entriesById := entries.ById(taskId)

	err = common.MarkExisiting(page, userDisplayName, entriesById, dateFormat)
	if err != nil {
		return fmt.Errorf("captureByTaskId: could not mark existing entries: %v", err)
	}

	weekGroups := entriesById.GroupByWeekNo()

	// Loop through each week group and create a new time entry for each week
	for _, weekEntries := range weekGroups {
		err = captureByWeek(page, weekEntries, dayFormat)
		if err != nil {
			return fmt.Errorf("captureByTaskId: could not log time entries for week: %v, error: %v", weekEntries[0].WeekNo, err)
		}
	}

	log.Println("Done loading")

	return nil
}

func captureByWeek(page playwright.Page, weekEntries at.TimeEntries, dayFormat string) error {
	peer := findWeekEntryPeer(weekEntries, weekEntries[0])

	if peer == nil {
		return captureNewWeek(page, weekEntries, dayFormat)

	} else {
		return captureExistingWeek(page, weekEntries, peer)
	}
}

func captureNewWeek(page playwright.Page, weekEntries at.TimeEntries, dayFormat string) error {
	if err := page.Locator("[data-eii='00000135']").Click(); err != nil {
		return fmt.Errorf("newWeekEntries: could not click new time entry button: %v", err)
	}

	navigateToWeek(page, weekEntries[0].Date, dayFormat)

	return captureWeek(page, weekEntries)
}

func navigateToWeek(page playwright.Page, entryTime time.Time, dayFormat string) error {
	entryWeekStart := at.SundayOfTheWeek(entryTime)

	for i := 0; i <= 3; i++ {

		s, err := page.Locator("body > div.Dialog1.Dialog2.Normal.Active tr.Heading > td.TextCell div.Label").First().TextContent()

		if err != nil {
			return fmt.Errorf("navigateToWeek: could not find pageDateLabel: %v", err)
		}

		firstParse, _ := time.Parse(dayFormat, s)
		if err != nil {
			return fmt.Errorf("navigateToWeek: could not parse pageWeekStart: %v", err)
		}

		inferredYear := at.InferYear(firstParse.Month(), 3, time.Now())
		pageWeekStart := time.Date(inferredYear, firstParse.Month(), firstParse.Day(), 0, 0, 0, 0, time.Local)

		if pageWeekStart.Year() == entryWeekStart.Year() &&
			pageWeekStart.Month() == entryWeekStart.Month() &&
			pageWeekStart.Day() == entryWeekStart.Day() {
			return nil

		} else if entryWeekStart.Before(pageWeekStart) {
			err = gotoPrevWeek(page)
			if err != nil {
				return fmt.Errorf("navigateToNextDay: could not find load indicator: %v", err)
			}
		} else {
			err = gotoNextWeek(page)
			if err != nil {
				return fmt.Errorf("navigateToNextDay: could not find load indicator: %v", err)
			}
		}
	}

	return errors.New("navigateToWeek: could not find week")
}

func gotoNextWeek(page playwright.Page) error {
	loadIndicator := "#LoadingIndicator.Active"
	page.Locator("body > div.Dialog1.Dialog2.Normal.Active .MoveRight").Click()
	return page.Locator(loadIndicator).WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateDetached})
}

func gotoPrevWeek(page playwright.Page) error {
	loadIndicator := "#LoadingIndicator.Active"
	page.Locator("body > div.Dialog1.Dialog2.Normal.Active .MoveLeft").Click()
	return page.Locator(loadIndicator).WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateDetached})
}

func captureExistingWeek(page playwright.Page, weekEntries at.TimeEntries, peer *at.TimeEntry) error {
	convLocator := peer.WeekPeerLocator.(playwright.Locator)

	err := convLocator.Locator("div.FooterActions div.LinkButton2").Nth(3).Click()

	if err != nil {
		return fmt.Errorf("editWeekEntries: could not click edit button: %v", err)
	}

	return captureWeek(page, weekEntries)
}

func captureWeek(page playwright.Page, weekEntries at.TimeEntries) error {
	weekEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active")
	err := weekEntryDialog.WaitFor()
	if err != nil {
		return fmt.Errorf("captureWeek: could not find weekEntryDialog: %v", err)
	}

	// Click Sunday's edit button
	timeEntryDialogSelector := page.Locator("div.Body > div.Scrolling > table > tbody div.Icon").First()

	err = timeEntryDialogSelector.Click()
	if err != nil {
		return fmt.Errorf("captureWeek: could not click Sunday's edit button: %v", err)
	}

	nextDayButton := page.Locator("[data-eii='0100014L']") // Next Day button
	err = nextDayButton.WaitFor()
	if err != nil {
		return fmt.Errorf("captureWeek: could not find timeEntryDialog: %v", err)
	}

	// Capture each week day's time if it exists
	sunday := at.SundayOfTheWeek(weekEntries[0].Date)

	entriesCaptured := 0

	for i := 0; i < 7; i++ {
		// Find the time entry for the current day
		entry := weekEntries.ByDate(sunday.AddDate(0, 0, i))

		if len(entry) > 1 {
			for _, e := range entry {
				e.SetError(fmt.Errorf("captureWeek: more than one entry for a given day: %v", sunday))
			}
		} else if len(entry) == 0 {
			// No time entry for this day, skip to the next day
		} else {
			te := entry[0]
			err = captureDay(page, te)
			if err != nil {
				te.SetError(fmt.Errorf("captureWeek: could not capture day: %v", err))
			}
			entriesCaptured++
		}

		if i >= 6 || entriesCaptured >= len(weekEntries) {
			err := saveWeek(page)
			if err != nil {
				return fmt.Errorf("captureWeek: could not save week: %v", err)
			}
			break
		} else {
			err := navigateToNextDay(page)
			if err != nil {
				return fmt.Errorf("captureWeek: could not navigate to next day: %v", err)
			}
		}
	}

	// Mark all entries as submitted
	for _, te := range weekEntries {
		if te.Error == nil {
			te.Submitted = true
		}
	}

	return nil
}

func captureDay(page playwright.Page, te *at.TimeEntry) error {
	log.Printf("Capture time entry: %+v\n", te)
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

	weekEntryDialog := page.Locator("body > div.Dialog1.Dialog2.Normal.Active").First()

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

	nextDayButton := page.Locator(nextDayButtonSelector)
	err := nextDayButton.Click()
	if err != nil {
		return fmt.Errorf("navigateToNextDay: could not click next day button: %v", err)
	}

	err = page.Locator(nextDayButtonSelector).WaitFor()
	if err != nil {
		return fmt.Errorf("navigateToNextDay: could not find next day button: %v", err)
	}

	return nil

}

func findWeekEntryPeer(entriesById at.TimeEntries, te *at.TimeEntry) *at.TimeEntry {
	weekNo := te.WeekNo

	for _, e := range entriesById {
		if e.WeekNo == weekNo && e.WeekPeerLocator != nil {
			return te
		}
	}

	return nil
}
