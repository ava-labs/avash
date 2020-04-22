package network

import (
	"fmt"
	"io/ioutil"
	"github.com/ava-labs/avash/node"
	"gopkg.in/yaml.v2"
)

// Config is a network configuration
type Config struct {
	Hosts []HostConfig
}

// HostConfig is a host configuration
type HostConfig struct {
	User, IP string
	Nodes    []NodeConfig
}

// NodeConfig is a node configuration
type NodeConfig struct {
	Name  string
	Flags node.Flags
}

// InitConfig returns a network configuration from `cfgpath`
func InitConfig(cfgpath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(cfgpath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Hosts) == 0 {
		return nil, fmt.Errorf("Must contain at least one host: %s", cfgpath)
	}
	for _, host := range cfg.Hosts {
		if host.User == "" {
			return nil, fmt.Errorf("Missing host name: %s", cfgpath)
		}
		if host.IP == "" {
			return nil, fmt.Errorf("Missing host IP: %s", cfgpath)
		}
		if len(host.Nodes) == 0 {
			return nil, fmt.Errorf("Must contain at least one node per host: %s", cfgpath)
		}
		for _, node := range host.Nodes {
			if node.Name == "" {
				return nil, fmt.Errorf("Missing node name: %s", cfgpath)
			}
		}
	}
	return &cfg, nil
}