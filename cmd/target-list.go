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
		InitLog()

		PrintTitle("Available Targets", true)

		PrintTargets()
	},
}

func init() {
	// enqueueCommand.Flags().StringVarP(&startHost, "bind", "b", "0.0.0.0", "Host to bind wall to")
	// enqueueCommand.Flags().IntVarP(&startPort, "port", "p", 3001, "Port to bind wall http server to")

	targetCommand.AddCommand(targetListCommand)
}
