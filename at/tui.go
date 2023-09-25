package at

import (
	"fmt"
	"log"

	"github.com/olekukonko/tablewriter"
)

// PrintSummary prints a summary table of the time entries.
func (entries TimeEntries) PrintSummary() {
	// Initialise a table writer using the log's writer.
	entries.SortByDateAndTime()
	table := tablewriter.NewWriter(log.Writer())
	table.SetAutoWrapText(false)
	// Set the table header.
	table.SetHeader([]string{"#", "AT-ID", "T", "Date", "Start", "Hrs", "EXS", "SAV", "ERR", "Project"})

	var total float32 = 0.0
	for i, e := range entries {
		// Check if there's an error for this entry.
		var errMsg string = ""
		if e.Error != nil {
			errMsg = "Y"
		}

		// Prepare a row for the table.
		row := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%d", e.Id),
			toPS(e.IsTicket),
			e.DateStr,
			e.StartTimeStr,
			fmt.Sprintf("%.2f", e.Duration),
			toYN(e.Exists),
			toYN(e.Submitted),
			errMsg,
			trim(e.Project, 45),
		}
		total += e.Duration
		table.Append(row)
	}

	// Set the table footer to show the total duration.
	table.SetFooter([]string{"", "", "", "Total", fmt.Sprintf("%.2f", total), "", "", "", "", "EOF"})
	table.Render()
}

// toPS converts a boolean indicating if an entry is a ticket to either "S" or "P".
func toPS(isTicket bool) string {
	return map[bool]string{true: "S", false: "P"}[isTicket]
}

// toYN converts a boolean to "Y" or "N".
func toYN(b bool) string {
	return map[bool]string{true: "Y", false: "N"}[b]
}

// trim truncates a string to a specified length.
func trim(s string, l int) string {
	asRunes := []rune(s)
	if len(asRunes) <= l {
		return s
	}
	return string(asRunes[:l])
}
