/*
Copyright © 2023 Philip Fourie
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/philipf/gt-at/autotask"
	"github.com/philipf/gt-at/pwplugin"
	"github.com/spf13/cobra"
)

var jsonfile string
var username string
var importFile bool
var prettyPrint bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gt-at",
	Short: "Go Time - AutoTasker",
	Long:  `Go Time - AutoTasker is a tool to help you track your time in AutoTask.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gt-at version 0.0.1")

		// Make sure the json file exists
		if _, err := os.Stat(jsonfile); os.IsNotExist(err) {
			fmt.Println("File does not exist: " + jsonfile)
			os.Exit(1)
		}

		creds := autotask.Credentials{
			Username: os.Getenv("AUTOTASK_USERNAME"),
			Password: os.Getenv("AUTOTASK_PASSWORD"),
		}

		err := load(jsonfile, autotask.LoadOptions{
			DryRun:          false,
			UserDisplayName: username,
			Credentials:     creds,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func load(filename string, opts autotask.LoadOptions) error {
	// Read the json file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	entries, err := autotask.UnmarshalToTimeEntries(data)
	if err != nil {
		return err
	}

	if prettyPrint {
		printSummary(entries)
	}

	if importFile {
		autoTasker := pwplugin.NewAutoTaskPlaywright()
		return autoTasker.LogTimes(entries, opts.Credentials, opts.UserDisplayName, "chromium", false, opts.DryRun)
	} else {
		return nil
	}
}

func printSummary(entries autotask.TimeEntries) {
	table := tablewriter.NewWriter(log.Writer())
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"#", "ID", "T", "Date", "Hrs", "Project"})

	var total float32 = 0.0
	for i, entry := range entries {
		entryType := "P"
		if entry.IsTicket {
			entryType = "S"
		}

		row := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%d", entry.Id),
			entryType,
			entry.DateStr,
			fmt.Sprintf("%.2f", entry.Duration),
			entry.Project,
		}
		total += entry.Duration
		table.Append(row)
	}

	table.SetFooter([]string{"", "", "", "Total", fmt.Sprintf("%.2f", total), ""}) // Add Footer

	table.Render() // Send output
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gt-at.yaml)")

	rootCmd.Flags().StringVarP(&username, "username", "u", "Philip Fourie", "name of user")
	rootCmd.Flags().StringVarP(&jsonfile, "filename", "f", "c:/tmp/time.json", "name of json file that will be imported")
	rootCmd.Flags().BoolVarP(&importFile, "import", "i", false, "import using browser")
	rootCmd.Flags().BoolVarP(&prettyPrint, "print", "p", true, "print a summary table befoe importing")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
