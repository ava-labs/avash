package node

import (
	"reflect"
)

// Flags represents available CLI flags when starting a node
type Flags struct {
	// Avash metadata
	ClientLocation string
	Meta           string
	DataDir        string

	// Assertions
	AssertionsEnabled bool

	// TX fees
	AvaTxFee uint

	// IP
	PublicIP string

	// Network ID
	NetworkID string

	// Throughput
	XputServerPort    uint
	XputServerEnabled bool

	// Crypto
	SignatureVerificationEnabled bool
	P2PTLSEnabled                bool

	// APIs
	APIAdminEnabled    bool
	APIIPCsEnabled     bool
	APIKeystoreEnabled bool
	APIMetricsEnabled  bool

	// HTTP
	HTTPPort        uint
	HTTPTLSEnabled  bool
	HTTPTLSCertFile string
	HTTPTLSKeyFile  string

	// Bootstrapping
	BootstrapIPs string
	BootstrapIDs string

	// Database
	DBEnabled bool
	DBDir     string

	// Plugins
	PluginDir string

	// Logging
	LogLevel string
	LogDir   string

	// Consensus
	SnowAvalancheBatchSize      int
	SnowAvalancheNumParents     int
	SnowSampleSize              int
	SnowQuorumSize              int
	SnowVirtuousCommitThreshold int
	SnowRogueCommitThreshold    int

	// Staking
	StakingTLSEnabled  bool
	StakingPort        uint
	StakingTLSKeyFile  string
	StakingTLSCertFile string
}

// FlagsYAML mimics Flags but uses pointers for proper YAML interpretation
// Note: FlagsYAML and Flags must always be identical in fields, otherwise parsing will break
type FlagsYAML struct {
	ClientLocation               *string `yaml:"-"`
	Meta                         *string `yaml:"-"`
	DataDir                      *string `yaml:"-"`
	AssertionsEnabled            *bool   `yaml:"assertions-enabled,omitempty"`
	AvaTxFee                     *uint   `yaml:"ava-tx-fee,omitempty"`
	PublicIP                     *string `yaml:"public-ip,omitempty"`
	NetworkID                    *string `yaml:"network-id,omitempty"`
	XputServerPort               *uint   `yaml:"xput-server-port,omitempty"`
	XputServerEnabled            *bool   `yaml:"xput-server-enabled,omitempty"`
	SignatureVerificationEnabled *bool   `yaml:"signature-verification-enabled,omitempty"`
	APIAdminEnabled              *bool   `yaml:"api-admin-enabled,omitempty"`
	APIIPCsEnabled               *bool   `yaml:"api-ipcs-enabled,omitempty"`
	APIKeystoreEnabled           *bool   `yaml:"api-keystore-enabled,omitempty"`
	APIMetricsEnabled            *bool   `yaml:"api-metrics-enabled,omitempty"`
	HTTPPort                     *uint   `yaml:"http-port,omitempty"`
	HTTPTLSEnabled               *bool   `yaml:"http-tls-enabled,omitempty"`
	HTTPTLSCertFile              *string `yaml:"http-tls-cert-file,omitempty"`
	HTTPTLSKeyFile               *string `yaml:"http-tls-key-file,omitempty"`
	BootstrapIPs                 *string `yaml:"bootstrap-ips,omitempty"`
	BootstrapIDs                 *string `yaml:"bootstrap-ids,omitempty"`
	DBEnabled                    *bool   `yaml:"db-enabled,omitempty"`
	DBDir                        *string `yaml:"db-dir,omitempty"`
	PluginDir                    *string `yaml:"plugin-dir,omitempty"`
	LogLevel                     *string `yaml:"log-level,omitempty"`
	LogDir                       *string `yaml:"log-dir,omitempty"`
	SnowAvalancheBatchSize       *int    `yaml:"snow-avalanche-batch-size,omitempty"`
	SnowAvalancheNumParents      *int    `yaml:"snow-avalanche-num-parents,omitempty"`
	SnowSampleSize               *int    `yaml:"snow-sample-size,omitempty"`
	SnowQuorumSize               *int    `yaml:"snow-quorum-size,omitempty"`
	SnowVirtuousCommitThreshold  *int    `yaml:"snow-virtuous-commit-threshold,omitempty"`
	SnowRogueCommitThreshold     *int    `yaml:"snow-rogue-commit-threshold,omitempty"`
	StakingTLSEnabled            *bool   `yaml:"staking-tls-enabled,omitempty"`
	StakingPort                  *uint   `yaml:"staking-port,omitempty"`
	StakingTLSKeyFile            *string `yaml:"staking-tls-key-file,omitempty"`
	StakingTLSCertFile           *string `yaml:"staking-tls-cert-file,omitempty"`
}

// SetDefaults sets any zero-value field to its default value
func (flags *Flags) SetDefaults() {
	f := reflect.Indirect(reflect.ValueOf(flags))
	d := reflect.ValueOf(DefaultFlags())
	for i := 0; i < f.NumField(); i++ {
		if f.Field(i).IsZero() {
			f.Field(i).Set(d.Field(i))
		}
	}
}

// ConvertYAML converts a FlagsYAML struct into a Flags struct
func ConvertYAML(flags FlagsYAML) Flags {
	var result Flags
	res := reflect.Indirect(reflect.ValueOf(&result))
	f := reflect.ValueOf(flags)
	d := reflect.ValueOf(DefaultFlags())
	for i := 0; i < res.NumField(); i++ {
		if f.Field(i).IsNil() {
			res.Field(i).Set(d.Field(i))
		} else {
			res.Field(i).Set(f.Field(i).Elem())
		}
	}
	return result
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
		APIIPCsEnabled:               true,
		APIKeystoreEnabled:           true,
		APIMetricsEnabled:            true,
		HTTPPort:                     9650,
		HTTPTLSEnabled:               false,
		HTTPTLSCertFile:              "",
		HTTPTLSKeyFile:               "",
		BootstrapIPs:                 "",
		BootstrapIDs:                 "",
		DBEnabled:                    true,
		DBDir:                        "db",
		PluginDir:                    "$GOPATH/src/github.com/ava-labs/gecko/build/plugins",
		LogLevel:                     "info",
		LogDir:                       "logs",
		SnowAvalancheBatchSize:       30,
		SnowAvalancheNumParents:      5,
		SnowSampleSize:               2,
		SnowQuorumSize:               2,
		SnowVirtuousCommitThreshold:  5,
		SnowRogueCommitThreshold:     10,
		P2PTLSEnabled:                true,
		StakingTLSEnabled:            false,
		StakingPort:                  9651,
		StakingTLSKeyFile:            "",
		StakingTLSCertFile:           "",
	}
}
