package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blang/semver"
	"github.com/fastlane-queue/fastlane-cmd/config"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Version of this package
var Version = "0.0.1"

// Config is the local user configurations for fastlane-cmd
var Config *config.Config

// Log is the Default Logger
var Log logrus.FieldLogger

// Verbose determines how verbose fastlane-cmd will run under
var Verbose int

// JSONLogFormat indicates that logs should be JSON
var JSONLogFormat bool

// SkipAutoUpdate indicates that auto-update should not be used
var SkipAutoUpdate bool

// NoColors indicates that logs should be JSON
var NoColors bool

var au aurora.Aurora

// RootCmd is the root command for fastlane-cmd CLI application
var RootCmd = &cobra.Command{
	Use:   "fastlane-cmd",
	Short: "fastlane-cmd runs jobs in fastlane",
	Long:  `Use fastlane-cmd to easily run jobs in fastlane.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Use == "fastlane-cmd" {
			os.Exit(0)
			return
		}

		if err := cmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func checkUpdate() {
	since := time.Now().Unix() - Config.LastUpdateCheck
	if since > 0 && since < 300 {
		return
	}

	lastRelease := getLastRelease()
	if lastRelease == nil {
		return
	}
	curr, err := semver.Make(Version)
	if err != nil {
		return
	}
	if lastRelease.Version.GT(curr) {
		fmt.Println(au.Sprintf(
			"Current version (%s) is obsolete.\nThere is a new version "+
				"available (%s).\nDo you want to update? [y/n]", au.Red(Version),
			au.Bold(lastRelease.Version),
		))

		if askForConfirmation() {
			doUpdate(lastRelease)
		}
	}

	Config.LastUpdateCheck = time.Now().Unix()
	err = Config.Serialize()
	if err != nil {
		return
	}
}

func doUpdate(release *Release) {
	fmt.Println("updating...")
}

// Execute runs RootCmd to initialize fastlane-cmd CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func Initialize() {
	InitLog()
	au = aurora.NewAurora(!NoColors)
	if !SkipAutoUpdate {
		checkUpdate()
	}
}

// PrintTitle for command line
func PrintTitle(title string, showVersion bool) {
	if showVersion {
		version := au.Green(fmt.Sprintf("v%s", Version))
		fmt.Println(au.Blue(fmt.Sprintf("fastlane-cmd %s\n", version)))
	}

	if title != "" {
		fmt.Println(au.Bold(title))
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(
		&Verbose, "verbose", "v", 0,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug",
	)

	RootCmd.PersistentFlags().BoolVarP(
		&JSONLogFormat, "json", "j", false,
		"JSON Log format (instead of text)",
	)

	RootCmd.PersistentFlags().BoolVarP(
		&NoColors, "no-colors", "c", false,
		"Don't show colored output",
	)

	RootCmd.PersistentFlags().BoolVarP(
		&SkipAutoUpdate, "skip-auto-update", "s", false,
		"Don't check if update is available",
	)

	var err error
	Config, err = config.NewConfig()
	if err != nil {
		panic(err)
	}
}

// InitLog structure
func InitLog() {
	ll := logrus.InfoLevel
	switch Verbose {
	case 0:
		ll = logrus.ErrorLevel
	case 1:
		ll = logrus.WarnLevel
	case 3:
		ll = logrus.DebugLevel
	}

	var log = logrus.New()
	if JSONLogFormat {
		log.Formatter = new(logrus.JSONFormatter)
	}
	log.Level = ll

	Log = log
}
