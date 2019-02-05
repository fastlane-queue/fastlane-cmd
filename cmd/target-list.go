package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// targetListCommand lists all the currently configured targets
var targetListCommand = &cobra.Command{
	Use:   "list",
	Short: "lists all the configured targets",
	Long:  `target list returns all the configured targets for the current user`,
	Run: func(cmd *cobra.Command, args []string) {
		InitLog()

		PrintTitle("Available Targets")

		var hosts = Config.List()
		table := tablewriter.NewWriter(os.Stdout)
		for _, host := range hosts {
			table.Append([]string{host.Name, host.URL})
		}
		table.SetHeader([]string{"Target Name", "Base URL"})
		table.Render()
	},
}

func init() {
	// enqueueCommand.Flags().StringVarP(&startHost, "bind", "b", "0.0.0.0", "Host to bind wall to")
	// enqueueCommand.Flags().IntVarP(&startPort, "port", "p", 3001, "Port to bind wall http server to")

	targetCommand.AddCommand(targetListCommand)
}
