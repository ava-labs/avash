package cmd

import (
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/network"
	"github.com/spf13/cobra"
)

// NetworkCommand represents the network command
var NetworkCommand = &cobra.Command{
	Use:   "network",
	Short: "Tools for interacting with remote hosts.",
	Long:  `Tools for interacting with remote hosts.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// SSHDeployCommand deploys a network config through an SSH client
var SSHDeployCommand = &cobra.Command{
	Use: "deploy [config file]",
	Short: "Deploys a remote network of nodes.",
	Long:  `Deploys a remote network of nodes from the provided config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		netCfg, err := network.InitConfig(args[0])
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Deployment starting... (this process typically takes 3-6 minutes depending on host)")
		if err := network.Deploy(netCfg, false); err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("All hosts finished.")
	},
}

// SSHRemoveCommand removes a network config through an SSH client
var SSHRemoveCommand = &cobra.Command{
	Use: "remove [config file]",
	Short: "Removes a remote network of nodes.",
	Long:  `Removes a remote network of nodes from the provided config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		netCfg, err := network.InitConfig(args[0])
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Removal starting...")
		if err := network.Remove(netCfg, false); err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("All hosts finished.")
	},
}

func init() {
	NetworkCommand.AddCommand(SSHDeployCommand)
	NetworkCommand.AddCommand(SSHRemoveCommand)
}