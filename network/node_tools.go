package network

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/kennygrant/sanitize"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
	"golang.org/x/crypto/ssh"
)

const (
	datadir string = "./stash"
)

// HostAuth represents a full set of host SSH credentials
type HostAuth struct {
	User, IP string
	Auth     ssh.AuthMethod
}

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

// InitAuth creates a mapping of host names to SSH authentications
// Requires user input at terminal
func InitAuth(deploys []DeployConfig) (map[string]HostAuth, error) {
	m := make(map[string]HostAuth)

	var authAll bool
	if len(deploys) > 1 {
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
	
	for _, deploy := range deploys {
		auth := authFunc(&deploy.IP)
		if auth == nil {
			return nil, fmt.Errorf("Authentication cancelled")
		}
		m[deploy.IP] = HostAuth{deploy.User, deploy.IP, auth}
	}
	return m, nil
}

// Deploy deploys nodes to hosts as specified in `config`
func Deploy(deploys []DeployConfig, isPrompt bool) error {
	log := cfg.Config.Log
	const cfp string = "./startnode.sh"

	authMap, err := InitAuth(deploys)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(deploys))
	for _, deploy := range deploys {
		go func(deploy DeployConfig) {
			defer wg.Done()
			hostAuth := authMap[deploy.IP]
			user, ip, auth := hostAuth.User, hostAuth.IP, hostAuth.Auth

			if err := InitHost(user, ip, auth); err != nil {
				log.Error("%s: %s", hostAuth.IP, err.Error())
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
			for _, n := range deploy.Nodes {
				if err := configureCLIFiles(&n.Flags, datadir, client); err != nil {
					log.Error("%s: %s", ip, err.Error())
					return
				}
				basename := sanitize.BaseName(n.Name)
				datapath := datadir + "/" + basename
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
		}(deploy)
	}
	wg.Wait()
	return nil
}

// Remove removes nodes from hosts as specified in `config`
func Remove(deploys []DeployConfig, isPrompt bool) error {
	log := cfg.Config.Log

	authMap, err := InitAuth(deploys)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(deploys))
	for _, deploy := range deploys {
		go func(deploy DeployConfig) {
			defer wg.Done()
			hostAuth := authMap[deploy.IP]
			user, ip, auth := hostAuth.User, hostAuth.IP, hostAuth.Auth

			client, err := NewSSH(user, ip, auth)
			if err != nil {
				log.Error("%s: %s", ip, err.Error())
				return
			}
			defer client.Close()

			var cmds []string
			for _, n := range deploy.Nodes {
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
		}(deploy)
	}
	wg.Wait()
	return nil
}

func configureCLIFiles(flags *node.Flags, datadir string, client *SSHClient) error {
	wd, _ := os.Getwd()
	if fp := flags.HTTPTLSCertFile; fp != "" {
		if string(fp[0]) != "/" {
			fp = fmt.Sprintf("%s/%s", wd, fp)
		}
		cfp := datadir + "/" + filepath.Base(fp)
		if err := client.CopyFile(fp, cfp); err != nil {
			return fmt.Errorf("%s: %s", err.Error(), fp)
		}
		flags.HTTPTLSCertFile = cfp
	}
	if fp := flags.HTTPTLSKeyFile; fp != "" {
		if string(fp[0]) != "/" {
			fp = fmt.Sprintf("%s/%s", wd, fp)
		}
		cfp := datadir + "/" + filepath.Base(fp)
		if err := client.CopyFile(fp, cfp); err != nil {
			return fmt.Errorf("%s: %s", err.Error(), fp)
		}
		flags.HTTPTLSKeyFile = cfp
	}
	if fp := flags.StakingTLSCertFile; fp != "" {
		if string(fp[0]) != "/" {
			fp = fmt.Sprintf("%s/%s", wd, fp)
		}
		cfp := datadir + "/" + filepath.Base(fp)
		if err := client.CopyFile(fp, cfp); err != nil {
			return fmt.Errorf("%s: %s", err.Error(), fp)
		}
		flags.StakingTLSCertFile = cfp
	}
	if fp := flags.StakingTLSKeyFile; fp != "" {
		if string(fp[0]) != "/" {
			fp = fmt.Sprintf("%s/%s", wd, fp)
		}
		cfp := datadir + "/" + filepath.Base(fp)
		if err := client.CopyFile(fp, cfp); err != nil {
			return fmt.Errorf("%s: %s", err.Error(), fp)
		}
		flags.StakingTLSKeyFile = cfp
	}
	return nil
}
