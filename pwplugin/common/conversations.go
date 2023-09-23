package common

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/philipf/gt-at/at"
	"github.com/playwright-community/playwright-go"
)

// MarkExisiting goes through timeEntries and marks them as existing if they are found on the page.
func MarkExisiting(page playwright.Page, userDisplayName string, timeEntries at.TimeEntries, dateFormat string) error {
	detailsSelector := page.Locator("div > .ConversationChunk > .ConversationItem .Details")
	convs, err := detailsSelector.All()

	if err != nil {
		return fmt.Errorf("markExistingEnties: could not find conversations: %v", err)
	}

	log.Printf("Found %v conversations\n", len(convs))

	for _, te := range timeEntries {
		for _, conv := range convs {
			author := conv.Locator("div > .Author div.Text2")
			authorName, err := author.TextContent()

			if err != nil {
				te.SetError(fmt.Errorf("markExistingEnties: could not find author TextContent: %+v", err))
				continue
			}

			if authorName == userDisplayName {
				timeDetail := conv.Locator("div.Title div.Text > span")

				t, err := timeDetail.TextContent()
				if err != nil {
					te.SetError(fmt.Errorf("markExistingEnties: could not find timeDetail TextContent: %+v", err))
					continue
				}

				weekNo := getConvWeekNo(t, dateFormat)

				if te.WeekNo == weekNo {
					te.WeekPeerLocator = conv
				}

				if strings.HasPrefix(t, te.DateStr) {
					log.Printf("Found date: %v\n", te.DateStr)
					te.Exists = true
					continue
				}
			}
		}
	}

	return nil
}

// getConvWeekNo parses the date string to extract its week number.
func getConvWeekNo(t, dateFormat string) int {
	dateStr := t[:len(dateFormat)]

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		// Logging instead of silently ignoring, this might provide useful debugging info.
		log.Printf("Error parsing date: %v\n", err)
		return -1
	}

	weekNo := at.WeekNo(date)
	return weekNo
}
