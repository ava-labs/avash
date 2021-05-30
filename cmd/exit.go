// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.
package cmd

import (
	"os"

	"github.com/ava-labs/avash/cfg"
	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/spf13/cobra"
)

// ExitCmd represents the exit command
var ExitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit the shell.",
	Long:  `Exit the shell, attempting to gracefully stop all processes first.`,
	Run: func(cmd *cobra.Command, args []string) {
		pmgr.ProcManager.StopAllProcesses()
		if pmgr.ProcManager.HasRunning() {
			cfg.Config.Log.Fatal("Unable to stop all processes, exiting anyway...")
			os.Exit(1)
		}
		cfg.Config.Log.Info("Cleanup successful, exiting...")
		os.Exit(0)
	},
}
