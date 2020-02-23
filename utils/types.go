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

func (v Verbosity) String() string {
	return []string{"debug", "normal"}[v]
}

// OutputType represents valid shell output locations
type OutputType int

// OutputType
const (
	Terminal	OutputType = iota
)

func (o OutputType) String() string {
	return []string{"terminal"}[o]
}