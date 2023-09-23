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
	// Filename for JSON import and flag to indicate if only a report is required.
	jsonfile   string
	reportOnly bool
)

// importCmd represents the import command for Cobra
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

	// Flags for the import command.
	importCmd.Flags().StringVarP(&jsonfile, "filename", "f", "/tmp/time.json", "name of json file that should be imported")
	importCmd.Flags().BoolVarP(&reportOnly, "reportOnly", "r", false, "print a summary of the time entries, but doesn't import them")
}

// load processes the file and imports it.
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

// getLoadOptions retrieves options for the load from configuration.
func getLoadOptions() at.CaptureOptions {
	// Assuming that getConfigFile() and other "setting..." constants are defined elsewhere in the code.
	viper.SetConfigFile(getConfigFile())
	err := viper.ReadInConfig()
	if err != nil {
		cobra.CheckErr(fmt.Errorf("fatal error config file: %s \n", err))
	}

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

// readFile reads and unmarshals a JSON file into time entries.
func readFile(filename, dateFormat string) (at.TimeEntries, error) {
	// Ensure the file exists.
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
