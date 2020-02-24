package utils

import (
	"fmt"
	"github.com/ava-labs/avash/cfg"
)

// PrintOutput outputs a message from an Output based on the config
func PrintOutput(s Output) {
	output := cfg.Config.Output
	msg := formatVerbosity(s, output.Verbosity)
	switch output.Type {
	case Terminal.String():
		fmt.Println(msg)
	default:
		fmt.Println(msg)
	}
}

func formatVerbosity(s Output, v string) string {
	switch v {
	case Norm.String():
		return s.Norm
	case Debug.String():
		return s.Debug
	default:
		return s.Norm
	}
}