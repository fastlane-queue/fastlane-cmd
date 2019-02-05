package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fastlane-queue/fastlane-cmd/config"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Version of this package
var Version = "0.1.0"

// Config is the local user configurations for fastlane-cmd
var Config *config.Config

// Log is the Default Logger
var Log logrus.FieldLogger

// Verbose determines how verbose fastlane-cmd will run under
var Verbose int

// JSONLogFormat indicates that logs should be JSON
var JSONLogFormat bool

// NoColors indicates that logs should be JSON
var NoColors bool

var au aurora.Aurora

// RootCmd is the root command for fastlane-cmd CLI application
var RootCmd = &cobra.Command{
	Use:   "fastlane-cmd",
	Short: "fastlane-cmd runs jobs in fastlane",
	Long:  `Use fastlane-cmd to easily run jobs in fastlane.`,
}

func getLastRelease() map[string]interface{} {
	url := "https://api.github.com/repos/fastlane-queue/fastlane-cmd/releases"
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, _ := netClient.Get(url)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(au.Red("Failed to verify if there's a new version of fastlane-cmd."), err)
			return nil
		}
		var result []interface{}
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			fmt.Println(au.Red("Failed to verify if there's a new version of fastlane-cmd."), err)
			return nil
		}
		if len(result) == 0 {
			return nil
		}
		return result[0].(map[string]interface{})
	}

	fmt.Println(
		au.Sprintf(
			au.Red("Failed to verify if there's a new version of fastlane-cmd (Status Code: %d)."), resp.StatusCode,
		),
	)
	return nil
}

func checkUpdate() {
	lastRelease := getLastRelease()
	if lastRelease == nil {
		return
	}
}

// Execute runs RootCmd to initialize fastlane-cmd CLI application
func Execute(cmd *cobra.Command) {
	au = aurora.NewAurora(!NoColors)
	checkUpdate()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// PrintTitle for command line
func PrintTitle(title string, showVersion bool) {
	if showVersion {
		version := au.Green(fmt.Sprintf("v%s", Version))
		fmt.Println(au.Blue(fmt.Sprintf("fastlane-cmd %s\n", version)))
	}

	if title != "" {
		// sep := strings.Repeat("-", utf8.RuneCountInString(title))
		fmt.Println(au.Bold(title))
		// fmt.Println(sep)
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
