/*
Copyright © 2023 Philip Fourie
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/philipf/gt-at/autotask"
	"github.com/philipf/gt-at/pwplugin"
	"github.com/spf13/cobra"
)

var jsonfile string
var username string

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

		err := load(jsonfile, autotask.LoadOptions{
			DryRun:          false,
			UserDisplayName: username,
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

	// prettry print the entries
	for _, e := range entries {
		fmt.Println(e)
	}

	autoTasker := pwplugin.NewAutoTaskPlaywright()
	return autoTasker.LogTimes(entries, opts.Credentials, opts.UserDisplayName, "chromium", false, opts.DryRun)
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
	rootCmd.Flags().StringVarP(&jsonfile, "filename", "f", "time.json", "name of json file that will be imported")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
