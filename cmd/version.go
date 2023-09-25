package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "0.0.10"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "development"
		}
		fmt.Println("gt-at version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
