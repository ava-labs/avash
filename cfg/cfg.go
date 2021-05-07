// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

// Package cfg manages the configuration file for avash
package cfg

import (
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ava-labs/avash/utils/logging"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Configuration is a shell-usable wrapper of the config file
type Configuration struct {
	AvalancheLocation, DataDir string
	Log                        logging.Log
}

type configFile struct {
	AvalancheLocation, DataDir string
	Log                        configFileLog
}

type configFileLog struct {
	Terminal, LogFile, Dir string
}

// Config is a global instance of the shell configuration
var Config Configuration

// DefaultCfgName is the default config filename
const DefaultCfgName = ".avash.yaml"

// DefaultCfgNameShort is the default config filename with yml extension
const DefaultCfgNameShort = ".avash.yml"

// InitConfig initializes the config for commands to reference
func InitConfig(cfgpath string) {
	cfgname := DefaultCfgName
	if cfgpath != "" {
		cfgpath, cfgname = filepath.Split(cfgpath)
		viper.AddConfigPath(cfgpath)
	}
	if !strings.HasSuffix(cfgname, ".yaml") && !strings.HasSuffix(cfgname, ".yml") {
		fmt.Println("Config filename must end with extension '.yaml' or '.yml'")
		os.Exit(1)
	}
	viper.SetConfigName(cfgname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/avash/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Invalid config path: %s%s\n", cfgpath, cfgname)
			os.Exit(1)
		}

		// try finding filename with yml extension
		viper.SetConfigName(DefaultCfgNameShort)
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Printf("Config file not found: %s%s\n", cfgpath, cfgname)
				home, _ := homedir.Dir()
				os.OpenFile(home+"/"+cfgname, os.O_RDONLY|os.O_CREATE, 0644)
				fmt.Printf("Created empty config file: %s/%s\n", home, cfgname)
				viper.SetConfigFile(home + "/" + cfgname)
			}
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

	// Set default `avalancheLocation` if missing
	if config.AvalancheLocation == "" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}
		config.AvalancheLocation = path.Join(gopath, "src", "github.com", "ava-labs", "avalanchego", "build", "avalanchego")
	}
	if _, err := os.Stat(config.AvalancheLocation); err != nil {
		fmt.Printf("Invalid avalanchego binary location: %s\n", config.AvalancheLocation)
		fmt.Println("Make sure your $GOPATH is set or provide a configuration file with a valid `avalancheLocation` value. See README.md for more details.")
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
		AvalancheLocation: config.AvalancheLocation,
		DataDir:           config.DataDir,
		Log:               *log,
	}
	Config.Log.Info("Config file set: %s", viper.ConfigFileUsed())
	Config.Log.Info("Avash successfully configured.")
}

func makeLogConfig(config configFileLog, dataDir string) logging.Config {
	terminalLvl, err := logging.ToLevel(config.Terminal)
	if err != nil && config.Terminal != "" {
		fmt.Printf("Invalid terminal log level '%s', defaulting to %s\n", config.Terminal, terminalLvl.String())
	}
	logFileLvl, err := logging.ToLevel(config.LogFile)
	if err != nil && config.LogFile != "" {
		fmt.Printf("Invalid logfile log level '%s', defaulting to %s\n", config.LogFile, logFileLvl.String())
	}
	if config.Dir == "" {
		defaultLogDir := dataDir + "/logs"
		config.Dir = defaultLogDir
	}
	return logging.Config{
		RotationInterval:            24 * time.Hour,
		FileSize:                    1 << 23, // 8 MB
		RotationSize:                7,
		FlushSize:                   1,
		DisableContextualDisplaying: true,
		DisplayLevel:                terminalLvl,
		LogLevel:                    logFileLvl,
		Directory:                   config.Dir,
	}
}
