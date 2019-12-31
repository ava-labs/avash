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

// Viper is a global instance of viper
var Viper *viper.Viper

// InitConfig initializes the config for commands to reference
func InitConfig() {
	cfgname := ".avash.yaml"
	Viper.SetConfigName(cfgname)
	viper.SetConfigType("yaml")
	Viper.AddConfigPath("$HOME/")
	Viper.AddConfigPath(".")
	Viper.AddConfigPath("/etc/avash/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, _ := homedir.Dir()
			os.OpenFile(home+"/"+cfgname, os.O_RDONLY|os.O_CREATE, 0644)
			Viper.SetConfigFile(home + "/" + cfgname)
			fmt.Println("SetConfig to: " + home + "/" + cfgname)
		}
	}

	if err := Viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
