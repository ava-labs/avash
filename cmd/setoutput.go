/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var SetOutputCmd = &cobra.Command{
	Use:	"setoutput [output type] [verbosity]",
	Short:	"Sets terminal output.",
	Long:	`Sets the location and verbosity of print output from the terminal.`,
	Run:	func(cmd *cobra.Command, args []string) {
		fmt.Printf("unimplemented\n")
	},
}