/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"os"

	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/spf13/cobra"
)

// ExitCmd represents the exit command
var ExitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit the shell.",
	Long:  `Exit the shell, attempting to gracefully stop all processes first.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := pmgr.ProcManager.StopAllProcesses()
		if err == nil {
			os.Exit(0)
		} else {
			panic("Unable to stop process " + name + ". Exitted anyway. Error: " + err.Error())
		}

	},
}
