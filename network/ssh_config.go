package network

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"github.com/ava-labs/avash/node"
	"gopkg.in/yaml.v2"
)

// RawConfig is a raw network configuration
type RawConfig struct {
	Hosts []struct{
		Name, User, IP string
	}
	Nodes []struct{
		Class string
		Flags node.FlagsYAML
	}
	Deploys []struct{
		Host  string
		Nodes []struct{
			Name, Class string
			Flags       node.FlagsYAML
		}
	}
}

// HostConfig is a host configuration
type HostConfig struct {
	User, IP string
}

// NodeConfig is a node configuration
type NodeConfig struct {
	Name  string
	Flags node.Flags
}

// DeployConfig is a deploy instruction for a particular host
type DeployConfig struct {
	User, IP string
	Nodes    []NodeConfig
}

// InitConfig returns a network configuration from `cfgpath`
func InitConfig(cfgpath string) ([]DeployConfig, error) {
	bytes, err := ioutil.ReadFile(cfgpath)
	if err != nil {
		return nil, err
	}
	var cfg RawConfig
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	if err := validateConfig(cfg, cfgpath); err != nil {
		return nil, err
	}
	deployCfg := buildDeploy(cfg)
	return deployCfg, nil
}

func validateConfig(cfg RawConfig, cfgpath string) error {
	if len(cfg.Hosts) == 0 {
		return fmt.Errorf("Config must contain at least one host: %s", cfgpath)
	}
	isHost := make(map[string]bool)
	isIP := make(map[string]bool)
	for _, host := range cfg.Hosts {
		if host.Name == "" {
			return fmt.Errorf("%s: host missing name", cfgpath)
		}
		if isHost[host.Name] {
			return fmt.Errorf("%s: duplicate host name: %s", cfgpath, host.Name)
		}
		if host.User == "" {
			return fmt.Errorf("%s: host missing user: %s", cfgpath, host.Name)
		}
		if host.IP == "" {
			return fmt.Errorf("%s: host missing IP address: %s", cfgpath, host.Name)
		}
		if isIP[host.IP] {
			return fmt.Errorf("%s: duplicate host IP address: %s", cfgpath, host.IP)
		}
		isHost[host.Name] = true
	}
	if len(cfg.Nodes) == 0 {
		return fmt.Errorf("%s: config must contain at least one node definition", cfgpath)
	}
	isNode := make(map[string]bool)
	for _, n := range cfg.Nodes {
		if n.Class == "" {
			return fmt.Errorf("%s: node definition missing class name", cfgpath)
		}
		if isNode[n.Class] {
			return fmt.Errorf("%s: duplicate node class name: %s", cfgpath, n.Class)
		}
		isNode[n.Class] = true
	}
	if len(cfg.Deploys) == 0 {
		return fmt.Errorf("%s: config must contain at least one deploy", cfgpath)
	}
	isDeployHost := make(map[string]bool)
	for _, deploy := range cfg.Deploys {
		if deploy.Host == "" {
			return fmt.Errorf("%s: deploy missing host name", cfgpath)
		}
		if !isHost[deploy.Host] {
			return fmt.Errorf("%s: deploy with undefined host name: %s", cfgpath, deploy.Host)
		}
		if isDeployHost[deploy.Host] {
			return fmt.Errorf("%s: duplicate deploy host target: %s", cfgpath, deploy.Host)
		}
		if len(deploy.Nodes) == 0 {
			return fmt.Errorf("%s: deploy must target at least one node", cfgpath)
		}
		isDeployHost[deploy.Host] = true
		isDeployNode := make(map[string]bool)
		for _, n := range deploy.Nodes {
			if n.Name == "" {
				return fmt.Errorf("%s: deploy node missing name", cfgpath)
			}
			if isDeployNode[n.Name] {
				return fmt.Errorf("%s: duplicate deploy node name for host '%s': %s", cfgpath, deploy.Host, n.Name)
			}
			if n.Class == "" {
				return fmt.Errorf("%s: deploy node missing class name", cfgpath)
			}
			if !isNode[n.Class] {
				return fmt.Errorf("%s: deploy node with undefined class name: %s", cfgpath, n.Class)
			}
			isDeployNode[n.Name] = true
		}
	}
	return nil
}

func buildDeploy(config RawConfig) []DeployConfig {
	hostMap := make(map[string]HostConfig)
	for _, host := range config.Hosts {
		hostMap[host.Name] = HostConfig{host.User, host.IP}
	}
	nodeMap := make(map[string]node.FlagsYAML)
	for _, n := range config.Nodes {
		nodeMap[n.Class] = n.Flags
	}
	var deploys []DeployConfig
	for _, deploy := range config.Deploys {
		host := hostMap[deploy.Host]
		var nodes []NodeConfig
		for _, n := range deploy.Nodes {
			flags := nodeMap[n.Class]
			// Override node class flags with deployment flags
			overrideFlags(&flags, n.Flags)
			nodes = append(nodes, NodeConfig{
				Name: n.Name,
				Flags: node.ConvertYAML(flags),
			})
		}
		deploys = append(deploys, DeployConfig{
			User: host.User,
			IP: host.IP,
			Nodes: nodes,
		})
	}
	return deploys
}

func overrideFlags(origFlags *node.FlagsYAML, overFlags node.FlagsYAML) {
	orig := reflect.Indirect(reflect.ValueOf(origFlags))
	over := reflect.ValueOf(overFlags)
	for i := 0; i < orig.NumField(); i++ {
		if orig.Field(i).IsNil() {
			orig.Field(i).Set(over.Field(i))
		}
	}
}