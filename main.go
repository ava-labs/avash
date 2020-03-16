package main

import (
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/cmd"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var confFlag string

func main() {
	// allow config file path to be set by user
	pflag.String("conf", ".avash.yaml", "Config file path")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cfg.InitConfig()
	cmd.RootCmd.AddCommand(cmd.AVAWalletCmd)
	cmd.RootCmd.AddCommand(cmd.ExitCmd)
	cmd.RootCmd.AddCommand(cmd.ProcmanagerCmd)
	cmd.RootCmd.AddCommand(cmd.RunScriptCmd)
	cmd.RootCmd.AddCommand(cmd.SetOutputCmd)
	cmd.RootCmd.AddCommand(cmd.StartnodeCmd)
	cmd.RootCmd.AddCommand(cmd.VarStoreCmd)
	cmd.Execute()
}
