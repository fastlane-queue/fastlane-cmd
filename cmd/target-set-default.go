package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// targetSetDefaultCommand lists all the currently configured targets
var targetSetDefaultCommand = &cobra.Command{
	Use:   "set-default",
	Short: "sets the default target",
	Long:  `target set-default specifies a target to be the default`,
	Run: func(cmd *cobra.Command, args []string) {
		InitLog()

		if len(args) < 1 {
			fmt.Println(au.Red("Invalid input. `target set-default` must be called with a target name:\n"))
			fmt.Println(au.Green("fastlane-cmd target set-default main"))
			os.Exit(1)
			return
		}

		name := args[0]
		if _, ok := Config.Hosts[name]; !ok {
			fmt.Println(au.Sprintf(au.Red("Target '%s' was not found."), au.Blue(name)))
			fmt.Println()
			PrintTitle("Available Targets", false)
			PrintTargets()
			os.Exit(1)
			return
		}

		Config.ClearDefaults()
		Config.Hosts[name].Default = true
		err := Config.Serialize()
		if err != nil {
			panic(err)
		}
		PrintTitle("", true)
		fmt.Println(au.Sprintf(au.Green("Target '%s' updated successfuly."), au.Blue(args[0])))

		fmt.Println()
		PrintTitle("Updated Targets", false)

		PrintTargets()
	},
}

func init() {
	targetCommand.AddCommand(targetSetDefaultCommand)
}
