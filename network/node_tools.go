package network

import (
	"fmt"
	"strings"
	"sync"
	"github.com/kennygrant/sanitize"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
)

const (
	datadir string = "./stash"
)

// InitHost initializes a host environment to be able to run nodes.
// Returns a connection to the host machine.
func InitHost(user, ip string, isPrompt bool) (*SSHClient, error) {
	const cfp string = "./init.sh"
	cmds := []string{
		"chmod 777 " + cfp,
		cfp,
	}

	client, err := NewSSH(user, ip, isPrompt)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	if err := client.CopyFile("network/init.sh", cfp); err != nil {
		return nil, err
	}
	defer client.RemovePath(cfp)

	if err := client.Run(cmds); err != nil {
		return nil, err
	}
	conn, err := client.Clone()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Deploy deploys nodes to hosts as specified in `config`
func Deploy(config *Config, isPrompt bool) error {
	log := cfg.Config.Log
	const cfp string = "./startnode.sh"

	var wg sync.WaitGroup
	wg.Add(len(config.Hosts))
	for _, host := range config.Hosts {
		go func(user, ip string, nodes []NodeConfig) {
			defer wg.Done()

			client, err := InitHost(user, ip, isPrompt)
			if err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}
			defer client.Close()

			if err := client.CopyFile("network/startnode.sh", cfp); err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}
			defer client.RemovePath(cfp)

			cmds := []string{
				fmt.Sprintf("chmod 777 %s", cfp),
			}
			for _, n := range nodes {
				basename := sanitize.BaseName(n.Name)
				datapath := datadir + "/" + basename
				n.Flags.SetDefaults()
				flags, _ := node.FlagsToArgs(n.Flags, datapath)
				args := strings.Join(flags, " ")
				cmd := fmt.Sprintf("%s --name=%s %s", cfp, n.Name, args)
				cmds = append(cmds, cmd)
			}

			if err := client.Run(cmds); err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}
			log.Info("%s: successfully deployed", ip)
		}(host.User, host.IP, host.Nodes)
	}
	wg.Wait()
	return nil
}

// Remove removes nodes from hosts as specified in `config`
func Remove(config *Config, isPrompt bool) error {
	log := cfg.Config.Log

	var wg sync.WaitGroup
	wg.Add(len(config.Hosts))
	for _, host := range config.Hosts {
		go func(user, ip string, nodes []NodeConfig) {
			defer wg.Done()

			client, err := NewSSH(user, ip, isPrompt)
			if err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}
			defer client.Close()

			var cmds []string
			for _, n := range nodes {
				tmpCmds := []string{
					fmt.Sprintf("docker stop %s", n.Name),
					fmt.Sprintf("docker rm %s", n.Name),
				}
				cmds = append(cmds, tmpCmds...)
			}

			if err := client.Run(cmds); err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}
			log.Info("%s: successfully removed", ip)
		}(host.User, host.IP, host.Nodes)
	}
	wg.Wait()
	return nil
}