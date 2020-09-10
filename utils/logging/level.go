package logging

import (
	"github.com/ava-labs/avalanche-go/utils/logging"
)

// Level ...
type Level = logging.Level

// ToLevel ...
func ToLevel(l string) (Level, error) {
	return logging.ToLevel(l)
}
