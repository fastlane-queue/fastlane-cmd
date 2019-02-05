package cmd

import (
	"github.com/spf13/cobra"
)

// targetListCommand lists all the currently configured targets
var targetListCommand = &cobra.Command{
	Use:   "list",
	Short: "lists all the configured targets",
	Long:  `target list returns all the configured targets for the current user`,
	Run: func(cmd *cobra.Command, args []string) {
		Initialize()

		PrintTitle("Available Targets", true)

		PrintTargets()
	},
}

func init() {
	targetCommand.AddCommand(targetListCommand)
}
