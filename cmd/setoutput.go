/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"fmt"
	"strings"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/utils"
	"github.com/spf13/cobra"
)

// SetOutputCmd sets the shell output type and verbosity
var SetOutputCmd = &cobra.Command{
	Use:		"setoutput [output type] [verbosity]",
	Short:		"Sets shell output.",
	Long:		`Sets the type and verbosity of print output from the shell.`,
	Run:	func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			var msg utils.Output
			outputType := strings.ToLower(args[0])
			verbosity := strings.ToLower(args[1])
			if utils.IsOutputType(outputType) && utils.IsVerbosity(verbosity) {
				cfg.Config.Output.Type = outputType
				cfg.Config.Output.Verbosity = verbosity
				msg.Norm = "shell output set"
				msg.Debug = fmt.Sprintf("shell output type set to %s, verbosity set to %s", outputType, verbosity)
				utils.PrintOutput(msg)
			} else {
				cmd.Help()
			}
		} else {
			cmd.Help()
		}
	},
}