/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

// Package cmd implements cobra commander
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"github.com/spf13/cobra"
)

type historyrecord struct {
	cmd   *cobra.Command
	flags []string
}

// Shell is a hlper struct for storing history and the instance of the shell prompt
type Shell struct {
	history []historyrecord // array of maps, map keys are "command" and "stdout", "stderr"
	rl      *readline.Instance
	root    *cobra.Command
}

func (sh *Shell) addHistory(cmd *cobra.Command, flags []string) {
	hr := historyrecord{
		cmd:   cmd,
		flags: flags,
	}
	sh.history = append(sh.history, hr)
}

func pcFromCommands(parent readline.PrefixCompleterInterface, c *cobra.Command) {
	pc := readline.PcItem(c.Use)
	parent.SetChildren(append(parent.GetChildren(), pc))
	for _, child := range c.Commands() {
		pcFromCommands(pc, child)
	}
}

// ShellLoop is an execution loop for the terminal application
func (sh *Shell) ShellLoop() {
	completer := readline.NewPrefixCompleter()
	for _, child := range sh.root.Commands() {
		pcFromCommands(completer, child)
	}
	rln, err := readline.NewEx(&readline.Config{
		Prompt:         "avash> ",
		UniqueEditLine: false,
	})
	sh.rl = rln
	if err != nil {
		panic(err)
	}
	defer sh.rl.Close()

	for {
		ln, err := sh.rl.Readline()
		if err != nil {
			continue
		}
		cmd, flags, err := sh.root.Find(strings.Fields(ln))
		if err != nil {
			sh.rl.Terminal.Write([]byte(err.Error()))
		}
		sh.addHistory(cmd, flags)
		cmd.ParseFlags(flags)
		cmd.Run(cmd, flags)

	}
}

// AvaShell is the shell for our little client
var AvaShell *Shell

func init() {
	AvaShell = new(Shell)
}

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "avash",
	Short: "A shell environment for one more more AVA nodes",
	Long:  "A shell environment for launching and interacting with multiple AVA nodes.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println()
		}
		AvaShell.root = cmd
		AvaShell.ShellLoop()
	},
}

// Execute runs the root command for avash
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
