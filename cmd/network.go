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

// SSHDeployCommand deploys a node through an SSH client
var SSHDeployCommand = &cobra.Command{
	Use: "deploy [node name] [SSH username] [IP address]",
	Short: "Deploys a remotely running node via SSH.",
	Long:  `Deploys a remotely running node via SSH to a specified host.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		netCfg := &network.Config{
			Hosts: []network.HostConfig{
				network.HostConfig{
					User: args[1],
					IP: args[2],
					Nodes: []string{
						args[0],
					},
				},
			},
		}
		if err := network.Deploy(netCfg, true); err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Node '%s' deployed by '%s' at %s", args[0], args[1], args[2])
	},
}

var SSHDeployAllCommand = &cobra.Command{
	Use: "deploy-all [config file]",
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
		log.Info("Deployment complete.")
	},
}

// SSHRemoveCommand removes a node through an SSH client
var SSHRemoveCommand = &cobra.Command{
	Use: "remove [node name] [SSH username] [IP address]",
	Short: "Removes a remotely running node via SSH.",
	Long:  `Removes a remotely running node via SSH on a specified host.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		netCfg := &network.Config{
			Hosts: []network.HostConfig{
				network.HostConfig{
					User: args[1],
					IP: args[2],
					Nodes: []string{
						args[0],
					},
				},
			},
		}
		if err := network.Remove(netCfg, true); err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Node '%s' removed by '%s' at %s", args[0], args[1], args[2])
	},
}

var SSHRemoveAllCommand = &cobra.Command{
	Use: "remove-all [config file]",
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
		log.Info("Removal complete.")
	},
}

func init() {
	NetworkCommand.AddCommand(SSHDeployCommand)
	NetworkCommand.AddCommand(SSHDeployAllCommand)
	NetworkCommand.AddCommand(SSHRemoveCommand)
	NetworkCommand.AddCommand(SSHRemoveAllCommand)
}