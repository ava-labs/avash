/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"fmt"

	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/utils/logging"
	"github.com/spf13/cobra"
)

// SetOutputCmd sets the shell output type and verbosity
var SetOutputCmd = &cobra.Command{
	Use:		"setoutput [log output] [log level]",
	Short:		"Sets log output.",
	Long:		`Sets the log level of a specific log output type.`,
	Run:	func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			output, outErr := logging.ToOutput(args[0])
			level, lvlErr := logging.ToLevel(args[1])
			if outErr == nil {
				if lvlErr == nil {
					cfg.Config.Log.SetLevel(output, level)
					fmt.Println("set shell output")
				} else {
					fmt.Println(lvlErr)
				}
			} else {
				fmt.Println(outErr)
			}
		} else {
			cmd.Help()
		}
	},
}