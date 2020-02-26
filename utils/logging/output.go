package logging

import (
	"fmt"
	"strings"
)

// Output ...
type Output int

// Enum ...
const (
	Terminal Output = iota
	LogFile
)

// ToOutput ...
func ToOutput(s string) (Output, error) {
	switch strings.ToUpper(s) {
	case "TERMINAL":
		return Terminal, nil
	case "LOGFILE":
		return LogFile, nil
	default:
		return Terminal, fmt.Errorf("unknown log output: %s", s)
	}
}

func (o Output) String() string {
	switch o {
	case Terminal:
		return "TERMINAL"
	case LogFile:
		return "LOGFILE"
	default:
		return "?????"
	}
}
