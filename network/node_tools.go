package network

import (
	"fmt"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey"
	"github.com/kennygrant/sanitize"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
	"golang.org/x/crypto/ssh"
)

const (
	datadir string = "./stash"
)

// InitHost initializes a host environment to be able to run nodes.
func InitHost(user, ip string, auth ssh.AuthMethod) error {
	const cfp string = "./init.sh"
	cmds := []string{
		"chmod 777 " + cfp,
		cfp,
	}

	client, err := NewSSH(user, ip, auth)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.CopyFile("network/init.sh", cfp); err != nil {
		return err
	}
	defer client.RemovePath(cfp)

	if err := client.Run(cmds); err != nil {
		return err
	}
	return nil
}

// InitAuth creates a mapping of IPs to SSH authentications
// Requires user input at terminal
func InitAuth(hosts []HostConfig) (map[string]ssh.AuthMethod, error) {
	m := make(map[string]ssh.AuthMethod)

	var authAll bool
	if len(hosts) > 1 {
		var authPrompt = &survey.Confirm{
			Message: "Would you like to use the same SSH credentials for all hosts?:",
		}
		survey.AskOne(authPrompt, &authAll)
	} else {
		authAll = false
	}

	var authFunc func(*string) ssh.AuthMethod
	if authAll {
		auth := PromptAuth(nil)
		authFunc = func(_ *string) ssh.AuthMethod {
			return auth
		}
	} else {
		authFunc = PromptAuth
	}
	
	for _, host := range hosts {
		auth := authFunc(&host.IP)
		if auth == nil {
			return nil, fmt.Errorf("Authentication cancelled")
		}
		m[host.IP] = auth
	}
	return m, nil
}

// Deploy deploys nodes to hosts as specified in `config`
func Deploy(config *Config, isPrompt bool) error {
	log := cfg.Config.Log
	const cfp string = "./startnode.sh"

	authMap, err := InitAuth(config.Hosts)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(config.Hosts))
	for _, host := range config.Hosts {
		go func(user, ip string, nodes []NodeConfig) {
			defer wg.Done()
			auth := authMap[ip]

			if err := InitHost(user, ip, auth); err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}

			client, err := NewSSH(user, ip, auth)
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
				flags, _ := node.FlagsToArgs(n.Flags, datapath, true)
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

	authMap, err := InitAuth(config.Hosts)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(config.Hosts))
	for _, host := range config.Hosts {
		go func(user, ip string, nodes []NodeConfig) {
			defer wg.Done()

			client, err := NewSSH(user, ip, authMap[ip])
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