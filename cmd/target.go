package cmd

import (
	"github.com/spf13/cobra"
)

// target - command that manages all the locally configured fastlane targets
var targetCommand = &cobra.Command{
	Use:   "target",
	Short: "target - configures fastlane servers",
	Long:  `Commands to manage the locally configured fastlane targets`,
}

func init() {
	RootCmd.AddCommand(targetCommand)
}
