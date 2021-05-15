// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cmd

import (
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/utils/logging"
	"github.com/spf13/cobra"
)

// SetOutputCmd sets the shell output type and verbosity
var SetOutputCmd = &cobra.Command{
	Use:   "setoutput [log output] [log level]",
	Short: "Sets log output.",
	Long:  `Sets the log level of a specific log output type.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}
		log := cfg.Config.Log
		output, outErr := logging.ToOutput(args[0])
		level, lvlErr := logging.ToLevel(args[1])
		if outErr != nil {
			log.Error(outErr.Error())
			return
		}
		if lvlErr != nil {
			log.Error(lvlErr.Error())
			return
		}
		log.SetLevel(output, level)
		log.Info("%s log level set: %s", output.String(), level.String())
	},
}
