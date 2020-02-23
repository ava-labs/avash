package main

import (
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/cmd"
	// "github.com/spf13/viper"
)

func main() {
	// cfg.Viper = viper.New()
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
