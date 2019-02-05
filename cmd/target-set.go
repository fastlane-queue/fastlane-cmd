package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var targetSetDefault bool

// targetSetCommand lists all the currently configured targets
var targetSetCommand = &cobra.Command{
	Use:   "set",
	Short: "sets the specified target",
	Long:  `target set specifies a target (new or existing) details`,
	Run: func(cmd *cobra.Command, args []string) {
		InitLog()

		if len(args) < 2 {
			fmt.Println(au.Red("Invalid input. `target set` must be called with a name and base URL, like so:\n"))
			fmt.Println(au.Green("fastlane-cmd target set main 'http://fastlane-server:10000'"))
			os.Exit(1)
			return
		}

		defaultTarget := targetSetDefault
		if len(Config.Hosts) == 0 {
			defaultTarget = true
		}

		Config.SetTarget(args[0], args[1], defaultTarget)

		err := Config.Serialize()
		if err != nil {
			panic(err)
		}
		PrintTitle("", true)
		fmt.Println(au.Sprintf(au.Green("Target '%s' updated successfuly to '%s'."), au.Blue(args[0]), au.Blue(args[1])))

		fmt.Println()
		PrintTitle("Updated Targets", false)

		PrintTargets()
	},
}

func init() {
	targetSetCommand.Flags().BoolVarP(&targetSetDefault, "default", "d", false, "Set this target to be the default target")
	targetCommand.AddCommand(targetSetCommand)
}
