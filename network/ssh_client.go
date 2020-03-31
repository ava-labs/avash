package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"github.com/AlecAivazis/survey"
	"github.com/ava-labs/avash/cfg"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHClient implements an SSH client
type SSHClient struct {
	*ssh.Client
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
		Message: "Full path to key file (PEM):",
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

func promptAuth() ssh.AuthMethod {
	const strPassword string = "password"
	const strKeyFile string = "key file"
	const strAgent string = "ssh agent"
	var sshMethod string
	authPrompt := &survey.Select{
		Message: "Choose a method to provide SSH credentials:",
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
		return &SSHClient{client}
	}
}

// NewSSH instantiates a new SSH client
func NewSSH(username string, ip string) (*SSHClient, error) {
	auth := promptAuth()
	if auth == nil {
		return nil, fmt.Errorf("Authentication quit")
	}
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ip + ":22", sshConfig)
	if err != nil {
		return nil, err
	}
	return &SSHClient{client}, nil
}

// TestOutput logs a test message through client
func (client *SSHClient) TestOutput() {
	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return
	}
	session.Run("echo this is a test")
	b, err := ioutil.ReadAll(stdout)
	cfg.Config.Log.Info(string(b))
}