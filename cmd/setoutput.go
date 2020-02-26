/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"fmt"

	"github.com/ava-labs/avash/utils/logging"
	gLogging "github.com/ava-labs/gecko/utils/logging"
	"github.com/spf13/cobra"
)

// SetOutputCmd sets the shell output type and verbosity
var SetOutputCmd = &cobra.Command{
	Use:		"setoutput [log output] [log level]",
	Short:		"Sets log output.",
	Long:		`Sets the log level of a specific log output type.`,
	Run:	func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			_, outErr := logging.ToOutput(args[0])
			_, lvlErr := gLogging.ToLevel(args[1])
			if  outErr == nil && lvlErr == nil {
				fmt.Println("unimplemented")
			} else {
				cmd.Help()
			}
		} else {
			cmd.Help()
		}
	},
}