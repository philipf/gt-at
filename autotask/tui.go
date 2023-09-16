package autotask

import (
	"fmt"
	"log"

	"github.com/olekukonko/tablewriter"
)

func (entries TimeEntries) PrintSummary() {
	table := tablewriter.NewWriter(log.Writer())
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"#", "AT-ID", "T", "Date", "Start", "Hrs", "EXS", "SAV", "ERR", "Project"})

	var total float32 = 0.0
	for i, e := range entries {

		var errMsg string = ""
		if e.Error != nil {
			errMsg = "Y"
		}

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

	table.SetFooter([]string{"", "", "", "Total", fmt.Sprintf("%.2f", total), "", "", "", "", "EOF"})
	table.Render()
}

func toPS(isTicket bool) string {
	return map[bool]string{true: "S", false: "P"}[isTicket]
}

func toYN(b bool) string {
	return map[bool]string{true: "Y", false: "N"}[b]
}

func trim(s string, l int) string {
	asRunes := []rune(s)
	if len(asRunes) <= l {
		return s
	}
	return string(asRunes[:l])
}
