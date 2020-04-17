package node

import (
	"reflect"
)

// Flags represents available CLI flags when starting a node
type Flags struct {
	// Avash metadata
    ClientLocation               string
    Meta                         string
	DataDir                      string
	
	// Assertions
	AssertionsEnabled            bool

	// TX fees
	AvaTxFee                     uint

	// IP
	PublicIP                     string

	// Network ID
	NetworkID                    string

	// Throughput
	XputServerPort               uint
	XputServerEnabled            bool

	// Crypto
	SignatureVerificationEnabled bool

	// APIs
	APIAdminEnabled              bool
	APIIPCsEnabled               bool
	APIKeystoreEnabled           bool
	APIMetricsEnabled            bool

	// HTTP
	HTTPPort                     uint
	HTTPTLSEnabled               bool
	HTTPTLSCertFile              string
	HTTPTLSKeyFile               string

	// Bootstrapping
	BootstrapIPs                 string
	BootstrapIDs                 string

	// Database
	DBEnabled                    bool
	DBDir                        string

	// Logging
	LogLevel                     string
	LogDir                       string

	// Consensus
	SnowAvalancheBatchSize       int
	SnowAvalancheNumParents      int
	SnowSampleSize               int
	SnowQuorumSize               int
	SnowVirtuousCommitThreshold  int
	SnowRogueCommitThreshold     int

	// Staking
	StakingTLSEnabled            bool
	StakingPort                  uint
	StakingTLSKeyFile            string
	StakingTLSCertFile           string
}

// SetDefaults sets any zero-value field to its default value
func (flags *Flags) SetDefaults() {
	f := reflect.Indirect(reflect.ValueOf(flags))
	z := reflect.ValueOf(Flags{})
	d := reflect.ValueOf(DefaultFlags())
	for i := 0; i < f.NumField(); i++ {
		if f.Field(i).Interface() == z.Field(i).Interface() {
			f.Field(i).Set(d.Field(i))
		}
	}
}

// DefaultFlags returns Avash-specific default node flags
func DefaultFlags() Flags {
	return Flags{
        ClientLocation:               "",
        Meta:                         "",
        DataDir:                      "",
		AssertionsEnabled:            true,
		AvaTxFee:                     0,
		PublicIP:                     "127.0.0.1",
		NetworkID:                    "local",
		XputServerPort:               9652,
		XputServerEnabled:            false,
		SignatureVerificationEnabled: true,
		APIAdminEnabled:              true,
		APIIPCsEnabled:			      true,
		APIKeystoreEnabled:           true,
		APIMetricsEnabled:            true,
		HTTPPort:                     9650,
		HTTPTLSEnabled:               false,
		HTTPTLSCertFile:              "",
		HTTPTLSKeyFile:               "",
		BootstrapIPs:                 "",
		BootstrapIDs:                 "",
		DBEnabled:                    true,
		DBDir:                        "db1",
		LogLevel:                     "info",
		LogDir:                       "logs",
		SnowAvalancheBatchSize:       30,
		SnowAvalancheNumParents:      5,
		SnowSampleSize:               2,
		SnowQuorumSize:               2,
		SnowVirtuousCommitThreshold:  5,
		SnowRogueCommitThreshold:     10,
		StakingTLSEnabled:            false,
		StakingPort:                  9651,
        StakingTLSKeyFile:            "",
		StakingTLSCertFile:           "",
	}
}
