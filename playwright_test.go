package autotask

import (
	"os"
	"testing"
)

func TestLogTimes(t *testing.T) {
	at := NewAutoTaskPlaywright()

	es := []*TimeEntry{
		{
			TicketId:  279750,
			Date:      "2023/09/13", // format to user locale
			StartTime: "10:29",
			EndTime:   "11:10",
			Summary: `Start   End    Time   Notes
10:29 - 11:10  00:40  10:30 stand-up
Duration: 0.75`,
		},

		{
			TicketId:  279750,
			Date:      "2023/09/14", // format to user locale
			StartTime: "10:30",
			EndTime:   "11:00",
			Summary: `Start   End    Time   Notes
10:29 - 11:10  00:40  10:30 stand-up
Duration: 0.75`,
		},
	}

	creds := Credentials{
		Username: os.Getenv("AUTOTASK_USERNAME"),
		Password: os.Getenv("AUTOTASK_PASSWORD"),
	}

	err := at.LogTimes(es, creds, "Philip Fourie", "chromium", false, false)

	if err != nil {
		t.Errorf("could not log times: %v", err)
	}
}
