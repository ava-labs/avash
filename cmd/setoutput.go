/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"fmt"
	// "github.com/ava-labs/avash/utils"
	"github.com/spf13/cobra"
)

// SetOutputCmd sets the shell output type and verbosity
var SetOutputCmd = &cobra.Command{
	Use:		"setoutput [output type] [verbosity]",
	Short:		"Sets shell output.",
	Long:		`Sets the type and verbosity of print output from the shell.`,
	Run:	func(cmd *cobra.Command, args []string) {
		fmt.Printf("unimplemented\n")
	},
}