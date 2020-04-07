package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"github.com/AlecAivazis/survey"
	"github.com/ava-labs/avash/cfg"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHClient implements an SSH client
type SSHClient struct {
	*ssh.Client
	config *ssh.ClientConfig
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
		return &SSHClient{client, config, host}
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
	return &SSHClient{client, sshConfig, ip}, nil
}

// Run runs a series of commands on remote host and waits for exit
func (client *SSHClient) Run(cmds []string) error {
	for _, cmd := range cmds {
		session, err := client.NewSession()
		if err != nil {
			return err
		}
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			return err
		}
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		cfg.Config.Log.Debug("Running command: %s", cmd)
		if err := session.Run(cmd); err != nil {
			return err
		}

		session.Close()
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
	cfg.Config.Log.Debug("%d bytes copied: %s --> %s", numBytes, fp, cfp)
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
	cfg.Config.Log.Debug("Removed: %s", path)
	return nil
}

// Clone creates another client instance connected to the same host
func (client *SSHClient) Clone() (*SSHClient, error) {
	clone, err := ssh.Dial("tcp", client.ip + ":22", client.config)
	if err != nil {
		return nil, err
	}
	return &SSHClient{clone, client.config, client.ip}, nil
}
