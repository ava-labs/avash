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

// NetworkSSHCommand represents the network ssh command
var NetworkSSHCommand = &cobra.Command{
	Use: "ssh",
	Short: "Tools for interacting with remote hosts via SSH.",
	Long:  `Tools for interacting with remote hosts via SSH.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// SSHDeployCommand deploys a node through an SSH client
var SSHDeployCommand = &cobra.Command{
	Use: "deploy [node name] [SSH username] [IP address]",
	Short: "Deploys a remotely running node.",
	Long:  `Deploys a remotely running node to a specified host.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := network.NewSSH(args[1], args[2])
		if err != nil {
			cfg.Config.Log.Error(err.Error())
			return
		}
		defer client.Close()
		client.TestOutput()
	},
}

func init() {
	NetworkSSHCommand.AddCommand(SSHDeployCommand)
	NetworkCommand.AddCommand(NetworkSSHCommand)
}