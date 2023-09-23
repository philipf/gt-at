package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gt-at",
	Short: "Go Time - AutoTasker",
	Long:  `Go Time - AutoTasker is a tool to help you track your time in AutoTask.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Print an error message
		fmt.Println("Error: Invalid usage. No subcommand provided.")

		// Display the help text
		cmd.Help()

		// Exit with an error code
		os.Exit(1)

	},
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.gt-at.yaml)")
}
