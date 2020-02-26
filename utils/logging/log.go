package logging

import (
	"fmt"
	"time"

	"github.com/ava-labs/gecko/utils/logging"
)

// Log is a wrapper struct for two distinct output logs
type Log struct {
	Terminal, LogFile	*logging.Log
}

// Config is a struct representation of the `log` field in the config file
type Config struct {
	Terminal, LogFile, Dir	string
}

// New ...
func New(config Config) (*Log, error) {
	terminalLvl, err := logging.ToLevel(config.Terminal)
	if err != nil {
		fmt.Printf("invalid terminal log level '%s', defaulting to %s\n", config.Terminal, terminalLvl.String())
	}
	terminalCfg := logging.Config{
		RotationInterval: 24 * time.Hour,
		FileSize:         1 << 23, // 8 MB
		RotationSize:     7,
		FlushSize:        1,
		DisableLogging:   true,
		Level:            logging.Info,
		Directory:        config.Dir,
	}
	terminalLog, err := logging.New(terminalCfg)
	if err != nil {
		return nil, err
	}
	logFileLvl, err := logging.ToLevel(config.LogFile)
	if err != nil {
		fmt.Printf("invalid logfile log level '%s', defaulting to %s\n", config.LogFile, logFileLvl.String())
	}
	logFileCfg := logging.Config{
		RotationInterval:  24 * time.Hour,
		FileSize:          1 << 23, // 8 MB
		RotationSize:      7,
		FlushSize:         1,
		DisableDisplaying: true,
		Level:             logging.Info,
		Directory:         config.Dir,
	}
	logFileLog, err := logging.New(logFileCfg)
	if err != nil {
		return nil, err
	}
	log := &Log{Terminal: terminalLog, LogFile: logFileLog}
	return log, nil
}

// Fatal ...
func (l *Log) Fatal(format string, args ...interface{}) {
	l.Terminal.Fatal(format, args...)
	l.LogFile.Fatal(format, args...)
}

// Error ...
func (l *Log) Error(format string, args ...interface{}) {
	l.Terminal.Error(format, args...)
	l.LogFile.Error(format, args...)
}

// Warn ...
func (l *Log) Warn(format string, args ...interface{}) {
	l.Terminal.Warn(format, args...)
	l.LogFile.Warn(format, args...)
}

// Info ...
func (l *Log) Info(format string, args ...interface{}) {
	l.Terminal.Info(format, args...)
	l.LogFile.Info(format, args...)
}

// Debug ...
func (l *Log) Debug(format string, args ...interface{}) {
	l.Terminal.Debug(format, args...)
	l.LogFile.Debug(format, args...)
}

// All ...
func (l *Log) All(format string, args ...interface{}) {
	l.Terminal.All(format, args...)
	l.LogFile.All(format, args...)
}
