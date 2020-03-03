/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

// Package cfg manages the configuration file for avash
package cfg

import (
	"fmt"
	"os"
	"time"

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
	Log						configFileLog
}

type configFileLog struct {
	Terminal, LogFile, Dir	string
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

	// Set default `datadir` if missing
	if config.DataDir == "" {
		wd, _ := os.Getwd()
		defaultDataDir := wd + "/stash"
		config.DataDir = defaultDataDir
	}
	if err := os.MkdirAll(config.DataDir, os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Configure and create log
	logCfg := makeLogConfig(config.Log, config.DataDir)
	log, err := logging.New(logCfg)
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

func makeLogConfig(config configFileLog, dataDir string) logging.Config {
	terminalLvl, err := logging.ToLevel(config.Terminal)
	if err != nil && config.Terminal != "" {
		fmt.Printf("invalid terminal log level '%s', defaulting to %s\n", config.Terminal, terminalLvl.String())
	}
	logFileLvl, err := logging.ToLevel(config.LogFile)
	if err != nil && config.LogFile != "" {
		fmt.Printf("invalid logfile log level '%s', defaulting to %s\n", config.LogFile, logFileLvl.String())
	}
	if config.Dir == "" {
		defaultLogDir := dataDir + "/logs"
		config.Dir = defaultLogDir
	}
	return logging.Config{
		RotationInterval:  24 * time.Hour,
		FileSize:          1 << 23, // 8 MB
		RotationSize:      7,
		FlushSize:         1,
		DisableDisplaying: true,
		DisplayLevel:      terminalLvl,
		LogLevel:          logFileLvl,
		Directory:         config.Dir,
	}
}
