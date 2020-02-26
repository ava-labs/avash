package types

// Configuration is a struct representation of the configuration file
type Configuration struct {
	AvaLocation	string
	DataDir		string
	Output		OutputConfig
}

// OutputConfig is a struct representation of the `output` config field
type OutputConfig struct {
	Terminal	LogLevel
	LogFile		LogLevel
}

// Log represents shell output produced by a command
type Log struct {
	entries	[]logEntry
}

// logEntry represents a specific message in part of shell output
type logEntry struct {
	msg	string
	lvl	LogLevel
}

// Add appends an entry to `log`
func (log *Log) Add(msg string, lvl LogLevel) {
	log.entries = append(log.entries, logEntry{msg, lvl})
}

// LogLevel represents valid shell output verbosities
type LogLevel int

// LogLevel
const (
	LogOff		LogLevel = iota
	LogFatal
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogAll
)

var logLevels = []string{"off", "fatal", "error", "warn", "info", "debug", "all"}

func (lvl LogLevel) String() string {
	return logLevels[lvl]
}

// IsLogLevel returns true if `s` is a `LogLevel` type, false otherwise
func IsLogLevel(s string) bool {
	i := getIndex(logLevels, s)
	return i != -1
}

// OutputType represents valid shell output locations
type OutputType int

// OutputType
const (
	Terminal	OutputType = iota
	LogFile
	OutputAll
)

var outputTypes = []string{"terminal", "logfile", "all"}

func (o OutputType) String() string {
	return outputTypes[o]
}

// IsOutputType returns true is `s` is an `OutputType` type, false otherwise
func IsOutputType(s string) bool {
	i := getIndex(outputTypes, s)
	return i != -1
}

func getIndex(arr []string, str string) int {
	for i, s := range arr {
		if s == str {
			return i
		}
	}
	return -1
}