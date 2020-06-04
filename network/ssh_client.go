package network

import (
	"io"
	"io/ioutil"
	"net"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/ava-labs/avash/cfg"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHClient implements an SSH client
type SSHClient struct {
	*ssh.Client
	ip string
}

func promptRetry(err error) bool {
	cfg.Config.Log.Error(err.Error())
	var retry bool
	var retryPrompt = &survey.Confirm{
		Message: "Would you like to retry?:",
	}
	survey.AskOne(retryPrompt, &retry)
	return retry
}

func promptPassword() ssh.AuthMethod {
	var pw string
	pwPrompt := &survey.Password{
		Message: "Password:",
	}
	survey.AskOne(pwPrompt, &pw)
	return ssh.Password(pw)
}

func promptKeyFile() ssh.AuthMethod {
	var fp string
	fpPrompt := &survey.Input{
		Message: "Full path to key file:",
	}
	for {
		survey.AskOne(fpPrompt, &fp)
		buff, err := ioutil.ReadFile(fp)
		if err != nil {
			if promptRetry(err) {
				continue
			}
			return nil
		}
		key, err := ssh.ParsePrivateKey(buff)
		if err != nil {
			if promptRetry(err) {
				continue
			}
			return nil
		}
		return ssh.PublicKeys(key)
	}
}

func sshAgent() ssh.AuthMethod {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		cfg.Config.Log.Error(err.Error())
		return nil
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
}

// PromptAuth returns an `ssh.AuthMethod` object for establishing an `SSHClient`
// Requires user input at terminal
func PromptAuth(ip *string) ssh.AuthMethod {
	const strPassword string = "password"
	const strKeyFile string = "key file"
	const strAgent string = "ssh agent"
	var sshMethod string
	msg := "Choose a method to provide SSH credentials"
	if ip != nil {
		msg += " for " + *ip
	}
	msg += ":"
	authPrompt := &survey.Select{
		Message: msg,
		Options: []string{strPassword, strKeyFile, strAgent, "quit"},
	}
	var auth ssh.AuthMethod
	for {
		survey.AskOne(authPrompt, &sshMethod)
		switch sshMethod {
		case strPassword:
			auth = promptPassword()
		case strKeyFile:
			auth = promptKeyFile()
		case strAgent:
			auth = sshAgent()
		default:
			return nil
		}
		if auth == nil {
			continue
		}
		return auth
	}
}

func promptClient(config *ssh.ClientConfig) *SSHClient {
	var host string
	hostPrompt := &survey.Input{
		Message: "Target host IP address:",
	}
	var port string
	portPrompt := &survey.Input{
		Message: "Target port:",
		Default: "22",
	}
	for {
		survey.AskOne(hostPrompt, &host)
		survey.AskOne(portPrompt, &port)
		client, err := ssh.Dial("tcp", host + ":" + port, config)
		if err != nil {
			if promptRetry(err) {
				continue
			}
			return nil
		}
		return &SSHClient{client, host}
	}
}

// NewSSH instantiates a new SSH client
func NewSSH(user, ip string, auth ssh.AuthMethod) (*SSHClient, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ip + ":22", sshConfig)
	if err != nil {
		return nil, err
	}
	return &SSHClient{client, ip}, nil
}

// Run runs a series of commands on remote host and waits for exit
func (client *SSHClient) Run(cmds []string) error {
	log := cfg.Config.Log
	for _, cmd := range cmds {
		if err := func(cmd string) error {
			session, err := client.NewSession()
			if err != nil {
				return err
			}
			defer session.Close()
			modes := ssh.TerminalModes{
				ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
				ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}
			if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
				return err
			}

			log.Info("%s: running command: %s", client.ip, cmd)
			bytes, err := session.CombinedOutput(cmd)
			sessionOutput := string(bytes)
			if err != nil {
				log.Error("%s: %s", client.ip, sessionOutput)
				return err
			}
			if sessionOutput != "" {
				log.Debug("%s: %s", client.ip, sessionOutput)
			}
			return nil
		}(cmd); err != nil {
			return err
		}
	}
	return nil
}

// CopyFile copies the contents of `fp` to `cfp` through client
func (client *SSHClient) CopyFile(fp string, cfp string) error {
	sftpClient, err := sftp.NewClient(client.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(cfp)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	numBytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	cfg.Config.Log.Debug("%s: %d bytes copied: %s --> %s", client.ip, numBytes, fp, cfp)
	return nil
}

// RemovePath removes file or directory `path` through client
func (client *SSHClient) RemovePath(path string) error {
	sftpClient, err := sftp.NewClient(client.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	if err := sftpClient.Remove(path); err != nil {
		return err
	}
	cfg.Config.Log.Debug("%s: removed: %s", client.ip, path)
	return nil
}
