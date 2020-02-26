/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

// Package cfg manages the configuration file for avash
package cfg

import (
	"fmt"
	"os"

	"github.com/ava-labs/avash/utils/logging"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Configuration is a shell-usable wrapper of the config file
type Configuration struct {
	AvaLocation, DataDir	string
	Log						logging.Log
}

type configFile struct {
	AvaLocation, DataDir	string
	Log						logging.Config
}

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

	var config configFile
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode config into struct, %v\n", err)
		os.Exit(1)
	}

	log, err := logging.New(config.Log)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Config = Configuration{
		AvaLocation:	config.AvaLocation,
		DataDir:		config.DataDir,
		Log:			*log,
	}
}
