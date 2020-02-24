package utils

// Output contains messages for all Verbosity types
type Output struct {
	Debug	string
	Norm	string
}

// Verbosity represents valid shell output verbosities
type Verbosity int

// Verbosity
const (
	Debug	Verbosity = iota
	Norm
)

var verbosities = []string{"debug", "normal"}

func (v Verbosity) String() string {
	return verbosities[v]
}

// IsVerbosity returns true if `s` is a `Verbosity` type, false otherwise
func IsVerbosity(s string) bool {
	return contains(verbosities, s)
}

// OutputType represents valid shell output locations
type OutputType int

// OutputType
const (
	Terminal	OutputType = iota
)

var outputTypes = []string{"terminal"}

func (o OutputType) String() string {
	return []string{"terminal"}[o]
}

// IsOutputType returns true is `s` is an `OutputType` type, false otherwise
func IsOutputType(s string) bool {
	return contains(outputTypes, s)
}

func contains(arr []string, s string) bool {
	for _, o := range arr {
		if o == s {
			return true
		}
	}
	return false
}