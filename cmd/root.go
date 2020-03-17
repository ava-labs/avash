/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

// Package cmd implements cobra commander
package cmd

import (
	"os"
	"strings"

	"github.com/ava-labs/avash/cfg"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
		if err := cmd.ParseFlags(flags); err != nil {
			cfg.Config.Log.Error(err.Error())
			continue
		}
		if err := cmd.ValidateArgs(flags); err != nil {
			cfg.Config.Log.Error(err.Error())
			continue
		}
		cmd.Run(cmd, flags)
	}
}

// AvaShell is the shell for our little client
var AvaShell *Shell

func init() {
	AvaShell = new(Shell)
	// allow config file path to be set by user
	pflag.String("config", ".avash.yaml", "Config file path")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cfg.InitConfig()
	RootCmd.AddCommand(AVAWalletCmd)
	RootCmd.AddCommand(ExitCmd)
	RootCmd.AddCommand(ProcmanagerCmd)
	RootCmd.AddCommand(RunScriptCmd)
	RootCmd.AddCommand(SetOutputCmd)
	RootCmd.AddCommand(StartnodeCmd)
	RootCmd.AddCommand(VarStoreCmd)
}

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:			"avash",
	Short:			"A shell environment for one more more AVA nodes",
	Long:			"A shell environment for launching and interacting with multiple AVA nodes.",
	SilenceUsage:	true,
	Args:			cobra.NoArgs,
	Run: 			func(cmd *cobra.Command, args []string) {
		AvaShell.root = cmd
		AvaShell.ShellLoop()
	},
}

// Execute runs the root command for avash
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
