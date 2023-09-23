package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise gt-at",
	Long:  `Creates the configuration file for gt-at`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		initialiseConfigFile(getConfigFile())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func isConfigured() {
	cfgFile := getConfigFile()

	_, err := os.Stat(cfgFile)

	// Check if file doesn't exist or there's another error accessing the file
	if err != nil {
		if os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("no config file found, please run `gt-at init` first: %v", err))
		} else {
			cobra.CheckErr(fmt.Errorf("error accessing config file: %v", err))
		}
	}
}

func getConfigFile() string {
	if cfgFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			cobra.CheckErr(fmt.Errorf("could not get user home directory: %v", err))
		}

		cfgFile = path.Join(homeDir, ".gt-at.yaml")
	}

	return cfgFile
}

func initialiseConfigFile(cf string) error {
	log.Printf("Initialising config file: %v\n", cf)

	setViperDefaults()

	viper.SetConfigFile(cf)
	err := viper.ReadInConfig() // Read the existing config if available
	if err != nil {
		cobra.CheckErr(fmt.Errorf("fatal error config file: %s", err))
	}

	// Initialise Viper settings
	setViperSetting("Your first name and last name in AutoTask (e.g Philip Fourie)", settingAutoTaskDisplayName)
	setViperSetting("Autotask date format, as configured AT preferences for your Profile, it should be defined using https://pkg.go.dev/time#pkg-constants (sorry)", settingAutoTaskDateFormat)
	setViperSetting("Autotask day format, as shown in AT week entries when capturing Tasks, it should be defined using https://pkg.go.dev/time#pkg-constants (sorry)", settingAutoTaskDayFormat)
	setViperSetting("Username, this is normally your company email address", settingCredentialsUsername)
	setViperSetting("Browser type (chromium|firefox|webkit)", settingPlaywrightBrowser)

	err = viper.WriteConfigAs(cf)
	if err != nil {
		return fmt.Errorf("cannot write config file in user's home directory:  [%v]", err)
	}

	return nil
}

func setViperSetting(question, setting string) {
	v, err := prompt(question, viper.GetString(setting))
	if err != nil {
		cobra.CheckErr(err)
	}

	viper.Set(setting, v)
}

func setViperDefaults() {
	viper.SetDefault(settingAutoTaskDisplayName, "")
	viper.SetDefault(settingAutoTaskDateFormat, "2006/01/02")
	viper.SetDefault(settingAutoTaskDayFormat, "Mon 01/02")

	viper.SetDefault(settingCredentialsUsername, "")

	viper.SetDefault(settingPlaywrightBrowser, "chromium")
	viper.SetDefault(settingPlaywrightHeadless, false)
}

const (
	settingAutoTaskDisplayName = "autotask.display-name"
	settingAutoTaskDateFormat  = "autotask.formats.date"
	settingAutoTaskDayFormat   = "autotask.formats.day"
	settingCredentialsUsername = "credentials.username"
	settingPlaywrightBrowser   = "playwright.browser-type"
	settingPlaywrightHeadless  = "playwright.headless"
)

func prompt(question, defaultValue string) (string, error) {
	if defaultValue == "" {
		fmt.Printf("%s:", question)
	} else {
		fmt.Printf("%s [%s]:", question, defaultValue)
	}

	v, err := readLine()
	if err != nil {
		return "", err
	}

	if v == "" {
		return defaultValue, nil
	}

	return v, nil
}

func readLine() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		return input, nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
