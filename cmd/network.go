package cmd

import (
	"fmt"
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
		log := cfg.Config.Log
		c1, err := network.NewSSH(args[1], args[2])
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer c1.Close()
		if err := initHost(c1); err != nil {
			log.Error(err.Error())
			return
		}
		// New connection necessary to refresh user groups
		c2, err := c1.Clone()
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer c2.Close()
		if err := deploy(c2, args[0]); err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Node successfully deployed!")
	},
}

// SSHKillCommand kills a node through an SSH client
var SSHKillCommand = &cobra.Command{
	Use: "kill [node name] [SSH username] [IP address]",
	Short: "Kills a remotely running node.",
	Long:  `Kills a remotely running node on a specified host.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		client, err := network.NewSSH(args[1], args[2])
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer client.Close()
		if err := kill(client, args[0]); err != nil {
			log.Error(err.Error())
			return
		}
	},
}

func initHost(client *network.SSHClient) error {
	const cfp string = "./init.sh"
	cmds := []string{
		"chmod 777 " + cfp,
		cfp,
	}

	if err := client.CopyFile("network/init.sh", cfp); err != nil {
		return err
	}
	defer client.RemovePath(cfp)

	if err := client.Run(cmds); err != nil {
		return err
	}
	return nil
}

func deploy(client *network.SSHClient, name string) error {
	const cfp string = "./startnode.sh"
	cmds := []string{
		fmt.Sprintf("chmod 777 %s", cfp),
		fmt.Sprintf("%s --name=%s --staking-tls-enabled=false", cfp, name),
	}

	if err := client.CopyFile("network/startnode.sh", cfp); err != nil {
		return err
	}
	defer client.RemovePath(cfp)

	if err := client.Run(cmds); err != nil {
		return err
	}
	return nil
}

func kill(client *network.SSHClient, name string) error {
	cmds := []string{
		fmt.Sprintf("docker kill %s", name),
	}
	if err := client.Run(cmds); err != nil {
		return err
	}
	return nil
}

func init() {
	NetworkSSHCommand.AddCommand(SSHDeployCommand)
	NetworkSSHCommand.AddCommand(SSHKillCommand)
	NetworkCommand.AddCommand(NetworkSSHCommand)
}