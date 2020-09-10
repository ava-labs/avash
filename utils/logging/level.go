package logging

import (
	"github.com/ava-labs/avalanchego/utils/logging"
)

// Level ...
type Level = logging.Level

// ToLevel ...
func ToLevel(l string) (Level, error) {
	return logging.ToLevel(l)
}
