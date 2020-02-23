/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

// Package cfg manages the configuration file for avash
package cfg

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Configuration struct {
	AvaLocation	string
	DataDir		string
	Output		OutputConfig
}

type OutputConfig struct {
	Type		string
	Verbosity	string
}

// Viper is a global instance of viper
// var Viper *viper.Viper

// Config is a global instance of the shell configuration
var Config Configuration

// InitConfig initializes the config for commands to reference
func InitConfig() {
	cfgname := ".avash.yaml"
	viper.SetConfigName(cfgname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/avash/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, _ := homedir.Dir()
			os.OpenFile(home+"/"+cfgname, os.O_RDONLY|os.O_CREATE, 0644)
			viper.SetConfigFile(home + "/" + cfgname)
			fmt.Println("SetConfig to: " + home + "/" + cfgname)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Printf("Unable to decode config into struct, %v\n", err)
		os.Exit(1)
	}
}
