package pwplugin

import (
	"os"
	"testing"

	"github.com/philipf/gt-at/autotask"
)

func TestLogTimes(t *testing.T) {
	at := NewAutoTaskPlaywright()

	es := []*autotask.TimeEntry{
		// 		{
		// 			//TicketId:  279750,
		// 			Id:        278364,
		// 			IsTicket:  true,
		// 			Date:      "2023/09/13", // format to user locale
		// 			StartTime: "10:29",
		// 			EndTime:   "11:10",
		// 			Duration:  0.75,
		// 			Summary: `Start   End    Time   Notes
		// 10:29 - 11:10  00:40  10:30 stand-up
		// Duration: 0.75`,
		// 		},

		{
			Id:           278364,
			IsTicket:     false,
			DateStr:      "2023/09/14", // format to user locale
			StartTimeStr: "10:30",
			EndTimeStr:   "11:00",
			Duration:     0.5,
			Summary: `Start   End    Time   Notes
10:29 - 11:10  00:40  10:30 Stand-up
Duration: 0.75`,
		},

		{
			Id:           278364,
			IsTicket:     false,
			DateStr:      "2023/09/11", // format to user locale
			StartTimeStr: "15:54",
			EndTimeStr:   "16:13",
			Duration:     0.25,
			Summary: `Start End Time Notes  
15:54 - 16:13 00:18 Catch-up /w Jon, issues about SW caused by User Assigned Managed Identity and CMK  
Duration: 0.25`,
		},
	}

	creds := autotask.Credentials{
		Username: os.Getenv("AUTOTASK_USERNAME"),
		Password: os.Getenv("AUTOTASK_PASSWORD"),
	}

	err := at.LogTimes(es, creds, "Philip Fourie", "chromium", false, false)

	if err != nil {
		t.Errorf("could not log times: %v", err)
	}
}
