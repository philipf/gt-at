package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/philipf/gt-at/autotask"
	"github.com/philipf/gt-at/pwplugin"
	"github.com/spf13/cobra"
)

var (
	jsonfile    string
	username    string
	importFile  bool
	prettyPrint bool
	dryRun      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gt-at",
	Short: "Go Time - AutoTasker",
	Long:  `Go Time - AutoTasker is a tool to help you track your time in AutoTask.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gt-at version 0.0.3")

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
			DryRun:          dryRun,
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
	log.Printf("Reading file: %v\n", filename)

	defer log.Println("Done")

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	log.Printf("Unmarshalling json file: %v\n", filename)
	entries, err := autotask.UnmarshalToTimeEntries(data)
	if err != nil {
		return err
	}

	if prettyPrint {
		log.Printf("Printing summary of time entries\n")
		entries.PrintSummary()
	}

	if importFile {
		log.Printf("Importing time entries\n")
		autoTasker := pwplugin.NewAutoTaskPlaywright()
		return autoTasker.LogTimes(entries, opts.Credentials, opts.UserDisplayName, "chromium", false, opts.DryRun)
	} else {
		return nil
	}
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
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gt-at.yaml)")

	rootCmd.Flags().StringVarP(&username, "username", "u", "Philip Fourie", "name of the user as it is displayed in Auto Task conversations")
	rootCmd.Flags().StringVarP(&jsonfile, "filename", "f", "/tmp/time.json", "name of json file that should be imported")
	rootCmd.Flags().BoolVarP(&importFile, "import", "i", false, "import using browser")
	rootCmd.Flags().BoolVarP(&prettyPrint, "print", "p", false, "print a summary table before importing")
	rootCmd.Flags().BoolVarP(&dryRun, "dry", "d", false, "performs a dry run and does not capture any time entries (WIP)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
