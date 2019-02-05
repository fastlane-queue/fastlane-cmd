package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/fastlane-queue/fastlane-cmd/config"
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

// RootCmd is the root command for fastlane-cmd CLI application
var RootCmd = &cobra.Command{
	Use:   "fastlane-cmd",
	Short: "fastlane-cmd runs jobs in fastlane",
	Long:  `Use fastlane-cmd to easily run jobs in fastlane.`,
}

// Execute runs RootCmd to initialize fastlane-cmd CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// PrintTitle for command line
func PrintTitle(title string) {
	fmt.Printf("fastlane-cmd v%s\n", Version)
	fmt.Println()

	sep := strings.Repeat("-", utf8.RuneCountInString(title))
	fmt.Println(title)
	fmt.Println(sep)
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
