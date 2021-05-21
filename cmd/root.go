// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

// Package cmd implements cobra commander
package cmd

import (
	"os"
	"strings"

	"github.com/ava-labs/avash/cfg"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const usageTmpl string = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

type historyrecord struct {
	cmd   *cobra.Command
	flags []string
}

// Shell is a helper struct for storing history and the instance of the shell prompt
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

func completerFromRoot(c *cobra.Command) []readline.PrefixCompleterInterface {
	var children []readline.PrefixCompleterInterface
	for _, child := range c.Commands() {
		childPC := readline.PcItem(child.Name(), completerFromRoot(child)...)
		children = append(children, childPC)
	}
	return children
}

// ShellLoop is an execution loop for the terminal application
func (sh *Shell) ShellLoop() {
	rootPC := completerFromRoot(sh.root)
	completer := readline.NewPrefixCompleter(rootPC...)
	rln, err := readline.NewEx(&readline.Config{
		Prompt:       "avash> ",
		AutoComplete: completer,
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

// AvalancheShell is the shell for our little client
var AvalancheShell *Shell
var RootCmd *cobra.Command

func init() {
	AvalancheShell = new(Shell)
	// allow config file path to be set by user
	var cfgpath string
	pflag.StringVar(&cfgpath, "config", cfg.DefaultCfgName, "Config file path")
	pflag.Parse()

	RootCmd = &cobra.Command{
		Use:   "avash",
		Short: "A shell environment for one or more Avalanche nodes",
		Long:  "A shell environment for launching and interacting with multiple Avalanche nodes.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			AvalancheShell.ShellLoop()
		},
		SilenceUsage: true,
	}

	cfg.InitConfig(cfgpath)
	RootCmd.AddCommand(AVAXWalletCmd)
	RootCmd.AddCommand(CallRPCCmd)
	RootCmd.AddCommand(ExitCmd)
	RootCmd.AddCommand(NetworkCommand)
	RootCmd.AddCommand(ProcmanagerCmd)
	RootCmd.AddCommand(RunScriptCmd)
	RootCmd.AddCommand(SetOutputCmd)
	RootCmd.AddCommand(StartnodeCmd)
	RootCmd.AddCommand(VarStoreCmd)
	RootCmd.SetUsageTemplate(usageTmpl)

	AvalancheShell.root = RootCmd

}

// Execute runs the root command for avash
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
