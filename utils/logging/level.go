package logging

import (
	"github.com/ava-labs/gecko/utils/logging"
)

// Level ...
type Level = logging.Level

// ToLevel ...
func ToLevel(l string) (Level, error) {
	return logging.ToLevel(l)
}
