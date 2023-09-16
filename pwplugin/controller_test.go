package pwplugin

import (
	"os"
	"testing"

	"github.com/philipf/gt-at/autotask"
)

func TestLogTimes(t *testing.T) {
	at := NewAutoTaskPlaywright()

	// es := []*autotask.TimeEntry{
	// 	autotask.NewEntry(279750, true, autotask.Date(2023, 9, 13), "10:29", 0.75,
	// 		"Start   End    Time   Notes\n10:29 - 11:10  00:40  10:30 stand-up\nDuration: 0.75"),

	// 	autotask.NewEntry(278364, false, autotask.Date(2023, 9, 13), "10:30", 0.5,
	// 		"Start   End    Time   Notes\n10:29 - 11:10  00:40  10:30 Stand-up\nDuration: 0.75"),

	// 	autotask.NewEntry(278364, false, autotask.Date(2023, 9, 17), "15:54", 0.25,
	// 		"Start   End    Time   Notes\n15:54 - 16:13 00:18 Catch-up /w Jon, issues about SW caused by User Assigned Managed Identity and CMK\nDuration: 0.25"),
	// }

	es := []*autotask.TimeEntry{
		autotask.NewEntry(266016, false, autotask.Date(2023, 9, 15), "10:30", 0.5,
			"Start   End    Time   Notes\n10:29 - 11:10  00:40  10:30 Stand-up\nDuration: 0.75"),
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
