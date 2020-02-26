package logging

import (
	"fmt"
	"time"

	"github.com/ava-labs/gecko/utils/logging"
)

// Log is a wrapper struct for the shell output log
type Log struct {
	log	*logging.Log
}

// Config is a struct representation of the `log` field in the config file
type Config struct {
	Terminal, LogFile, Dir	string
}

// New ...
func New(config Config) (*Log, error) {
	terminalLvl, err := ToLevel(config.Terminal)
	if err != nil {
		fmt.Printf("invalid terminal log level '%s', defaulting to %s\n", config.Terminal, terminalLvl.String())
	}
	logFileLvl, err := ToLevel(config.LogFile)
	if err != nil {
		fmt.Printf("invalid logfile log level '%s', defaulting to %s\n", config.LogFile, logFileLvl.String())
	}
	logCfg := logging.Config{
		RotationInterval:  24 * time.Hour,
		FileSize:          1 << 23, // 8 MB
		RotationSize:      7,
		FlushSize:         1,
		DisableDisplaying: true,
		DisplayLevel:      terminalLvl,
		LogLevel:          logFileLvl,
		Directory:         config.Dir,
	}
	log, err := logging.New(logCfg)
	if err != nil {
		return nil, err
	}
	return &Log{log: log}, nil
}

// Fatal ...
func (l *Log) Fatal(format string, args ...interface{}) { l.log.Fatal(format, args...) }

// Error ...
func (l *Log) Error(format string, args ...interface{}) { l.log.Error(format, args...) }

// Warn ...
func (l *Log) Warn(format string, args ...interface{}) { l.log.Warn(format, args...) }

// Info ...
func (l *Log) Info(format string, args ...interface{}) { l.log.Info(format, args...) }

// Debug ...
func (l *Log) Debug(format string, args ...interface{}) { l.log.Debug(format, args...) }

// All ...
func (l *Log) All(format string, args ...interface{}) { l.log.All(format, args...) }

// SetLevel ...
func (l *Log) SetLevel(out Output, lvl Level) {
	switch out {
	case Terminal:
		l.log.SetDisplayLevel(lvl)
	case LogFile:
		l.log.SetLogLevel(lvl)
	case All:
		l.log.SetDisplayLevel(lvl)
		l.log.SetLogLevel(lvl)
	default:
		// do nothing
	}
}