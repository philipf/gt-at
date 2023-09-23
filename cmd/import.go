package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/philipf/gt-at/at"
	"github.com/philipf/gt-at/pwplugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	jsonfile   string
	reportOnly bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file of time entries into AutoTask",
	Long:  `Import a file of time entries into AutoTask`,

	Run: func(cmd *cobra.Command, args []string) {
		isConfigured()

		err := load(jsonfile)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&jsonfile, "filename", "f", "/tmp/time.json", "name of json file that should be imported")
	importCmd.Flags().BoolVarP(&reportOnly, "reportOnly", "r", false, "print a summary of the time entries, but doesn't import them")
}

func load(filename string) error {
	log.Printf("Loading file: %v\n", filename)
	defer log.Println("Done")

	opts := getLoadOptions()

	entries, err := readFile(filename, opts.DateFormat)
	if err != nil {
		return err
	}

	entries.PrintSummary()

	if reportOnly {
		return nil
	}

	log.Printf("Importing time entries\n")
	autoTasker := pwplugin.NewAutoTaskPlaywright()
	err = autoTasker.CaptureTimes(entries, opts)
	entries.PrintSummary()
	if err != nil {
		return err
	}

	return nil
}

func getLoadOptions() at.CaptureOptions {
	viper.SetConfigFile(getConfigFile())
	viper.ReadInConfig()

	opts := at.CaptureOptions{
		Credentials: at.Credentials{
			Username: viper.GetString(settingCredentialsUsername),
		},
		UserDisplayName: viper.GetString(settingAutoTaskDisplayName),
		DateFormat:      viper.GetString(settingAutoTaskDateFormat),
		DayFormat:       viper.GetString(settingAutoTaskDayFormat),
		BrowserType:     viper.GetString(settingPlaywrightBrowser),
		Headless:        viper.GetBool(settingPlaywrightHeadless),
		DryRun:          false,
	}

	return opts
}

func readFile(filename, dateFormat string) (at.TimeEntries, error) {
	// make sure the file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %v, %v", filename, err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	log.Printf("Unmarshalling json file: %v\n", filename)
	entries, err := at.UnmarshalToTimeEntries(data, dateFormat)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
