package logging

import (
	"github.com/ava-labs/avalanchego/utils/logging"
)

// Log is a wrapper struct for the shell output log
type Log struct {
	logging.Logger
}

// Config is a struct representation of the `log` field in the config file
type Config = logging.Config

// New ...
func New(config Config) (*Log, error) {
	logFactory := logging.NewFactory(config)
	log, err := logFactory.Make("avash")
	if err != nil {
		return nil, err
	}
	return &Log{log}, nil
}

// SetLevel ...
func (l *Log) SetLevel(out Output, lvl Level) {
	switch out {
	case Terminal:
		l.SetDisplayLevel(lvl)
	case LogFile:
		l.SetLogLevel(lvl)
	case All:
		l.SetDisplayLevel(lvl)
		l.SetLogLevel(lvl)
	default:
		// do nothing
	}
}
