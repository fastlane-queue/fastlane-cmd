package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// target - command that manages all the locally configured fastlane targets
var targetCommand = &cobra.Command{
	Use:   "target",
	Short: "target - configures fastlane servers",
	Long:  `Commands to manage the locally configured fastlane targets`,
}

// PrintTargets configured in .fastlanerc
func PrintTargets() {
	var hosts = Config.List()

	if len(hosts) == 0 {
		fmt.Println()
		fmt.Println(au.Red("No targets found for the current user."))
		fmt.Println()
		command := au.Green("fastlane-cmd target set main http://fastlane-server:10000")
		fmt.Printf("Use '%s' to add a new one.\n", command)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	for _, host := range hosts {
		isDefault := ""
		if host.Default {
			isDefault = "*"
		}
		lastUsedTime := "Never"
		if host.LastUsed != 0 {
			tm := time.Unix(host.LastUsed, 0)
			lastUsedTime = tm.Format("Jan 2 2006 15:04")
		}
		table.Append([]string{host.Name, host.URL, lastUsedTime, isDefault})
	}
	table.SetHeader([]string{"Target Name", "Base URL", "Last Usage", "Default"})
	table.Render()
}

func init() {
	RootCmd.AddCommand(targetCommand)
}
