package cmd

import (
	"github.com/spf13/cobra"
)

// enqueueCommand represents the start command
var enqueueCommand = &cobra.Command{
	Use:   "enqueue",
	Short: "enqueue a new job in fastlane",
	Long:  `Enqueues a new job in fastlane with the specified options`,
	Run: func(cmd *cobra.Command, args []string) {
		InitLog()
		Log.Info("Enqueuing job...")
		err := Config.UpdateLastUsed("main")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	// enqueueCommand.Flags().StringVarP(&startHost, "bind", "b", "0.0.0.0", "Host to bind wall to")
	// enqueueCommand.Flags().IntVarP(&startPort, "port", "p", 3001, "Port to bind wall http server to")

	RootCmd.AddCommand(enqueueCommand)
}
