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

	// Version
	Version bool

	// TX fees
	TxFee uint

	// IP
	PublicIP              string
	DynamicUpdateDuration string
	DynamicPublicIP       string

	// Network ID
	NetworkID string

	// Crypto
	SignatureVerificationEnabled bool

	// APIs
	APIAdminEnabled    bool
	APIIPCsEnabled     bool
	APIKeystoreEnabled bool
	APIMetricsEnabled  bool
	APIHealthEnabled   bool
	APIInfoEnabled     bool

	// HTTP
	HTTPHost        string
	HTTPPort        uint
	HTTPTLSEnabled  bool
	HTTPTLSCertFile string
	HTTPTLSKeyFile  string

	// Bootstrapping
	BootstrapIPs                     string
	BootstrapIDs                     string
	BootstrapBeaconConnectionTimeout string

	// Database
	DBType string
	DBDir  string

	// Build
	BuildDir string

	// Plugins
	PluginDir string

	// Logging
	LogLevel            string
	LogDir              string
	LogDisplayLevel     string
	LogDisplayHighlight string

	// Consensus
	SnowAvalancheBatchSize      int
	SnowAvalancheNumParents     int
	SnowSampleSize              int
	SnowQuorumSize              int
	SnowVirtuousCommitThreshold int
	SnowRogueCommitThreshold    int
	SnowConcurrentRepolls       int
	MinDelegatorStake           int
	ConsensusShutdownTimeout    string
	ConsensusGossipFrequency    string
	MinDelegationFee            int
	MinValidatorStake           int
	MaxStakeDuration            string
	MaxValidatorStake           int

	// Staking
	StakingEnabled        bool
	StakeMintingPeriod    string
	StakingPort           uint
	StakingDisabledWeight int
	StakingTLSKeyFile     string
	StakingTLSCertFile    string

	// Auth
	APIAuthRequired        bool
	APIAuthPasswordFileKey string
	MinStakeDuration       string

	// Whitelisted Subnets
	WhitelistedSubnets string

	// Config
	ConfigFile string

	// IPCS
	IPCSChainIDs string
	IPCSPath     string

	// File Descriptor Limit
	FDLimit int

	// Benchlist
	BenchlistFailThreshold      int
	BenchlistMinFailingDuration string
	BenchlistPeerSummaryEnabled bool
	BenchlistDuration           string

	// Network
	NetworkInitialTimeout                   string
	NetworkMinimumTimeout                   string
	NetworkMaximumTimeout                   string
	NetworkHealthMaxSendFailRateKey         float64
	NetworkHealthMaxPortionSendQueueFillKey float64
	NetworkHealthMaxTimeSinceMsgSentKey     string
	NetworkHealthMaxTimeSinceMsgReceivedKey string
	NetworkHealthMinConnPeers               int
	NetworkTimeoutCoefficient               int
	NetworkTimeoutHalflife                  string
	NetworkCompressionEnabled               bool

	// Peer List Gossiping
	NetworkPeerListGossipFrequency string
	NetworkPeerListGossipSize      int
	NetworkPeerListSize            int

	// Uptime Requirement
	UptimeRequirement float64

	// Retry
	RetryBootstrapWarnFrequency int
	RetryBootstrap              bool

	// Health
	HealthCheckAveragerHalflifeKey string
	HealthCheckFreqKey             string

	// Router
	RouterHealthMaxOutstandingRequestsKey int
	RouterHealthMaxDropRateKey            float64

	IndexEnabled bool

	PluginModeEnabled bool

	MeterVMsEnabled bool
}

// FlagsYAML mimics Flags but uses pointers for proper YAML interpretation
// Note: FlagsYAML and Flags must always be identical in fields, otherwise parsing will break
type FlagsYAML struct {
	ClientLocation                          *string  `yaml:"-"`
	Meta                                    *string  `yaml:"-"`
	DataDir                                 *string  `yaml:"-"`
	AssertionsEnabled                       *bool    `yaml:"assertions-enabled,omitempty"`
	Version                                 *bool    `yaml:"version,omitempty"`
	TxFee                                   *uint    `yaml:"tx-fee,omitempty"`
	PublicIP                                *string  `yaml:"public-ip,omitempty"`
	DynamicPublicIP                         *string  `yaml:"dynamic-public-ip,omitempty"`
	NetworkID                               *string  `yaml:"network-id,omitempty"`
	SignatureVerificationEnabled            *bool    `yaml:"signature-verification-enabled,omitempty"`
	APIAdminEnabled                         *bool    `yaml:"api-admin-enabled,omitempty"`
	APIIPCsEnabled                          *bool    `yaml:"api-ipcs-enabled,omitempty"`
	APIKeystoreEnabled                      *bool    `yaml:"api-keystore-enabled,omitempty"`
	APIMetricsEnabled                       *bool    `yaml:"api-metrics-enabled,omitempty"`
	HTTPHost                                *string  `yaml:"http-host,omitempty"`
	HTTPPort                                *uint    `yaml:"http-port,omitempty"`
	HTTPTLSEnabled                          *bool    `yaml:"http-tls-enabled,omitempty"`
	HTTPTLSCertFile                         *string  `yaml:"http-tls-cert-file,omitempty"`
	HTTPTLSKeyFile                          *string  `yaml:"http-tls-key-file,omitempty"`
	BootstrapIPs                            *string  `yaml:"bootstrap-ips,omitempty"`
	BootstrapIDs                            *string  `yaml:"bootstrap-ids,omitempty"`
	BootstrapBeaconConnectionTimeout        *string  `yaml:"bootstrap-beacon-connection-timeout,omitempty"`
	DBType                                  *string  `yaml:"db-type,omitempty"`
	DBDir                                   *string  `yaml:"db-dir,omitempty"`
	BuildDir                                *string  `yaml:"build-dir,omitempty"`
	PluginDir                               *string  `yaml:"plugin-dir,omitempty"`
	LogLevel                                *string  `yaml:"log-level,omitempty"`
	LogDir                                  *string  `yaml:"log-dir,omitempty"`
	LogDisplayLevel                         *string  `yaml:"log-display-level,omitempty"`
	LogDisplayHighlight                     *string  `yaml:"log-display-highlight,omitempty"`
	SnowAvalancheBatchSize                  *int     `yaml:"snow-avalanche-batch-size,omitempty"`
	SnowAvalancheNumParents                 *int     `yaml:"snow-avalanche-num-parents,omitempty"`
	SnowSampleSize                          *int     `yaml:"snow-sample-size,omitempty"`
	SnowQuorumSize                          *int     `yaml:"snow-quorum-size,omitempty"`
	SnowVirtuousCommitThreshold             *int     `yaml:"snow-virtuous-commit-threshold,omitempty"`
	SnowRogueCommitThreshold                *int     `yaml:"snow-rogue-commit-threshold,omitempty"`
	SnowConcurrentRepolls                   *int     `yaml:"snow-concurrent-repolls,omitempty"`
	MinDelegatorStake                       *int     `yaml:"min-delegator-stake,omitempty"`
	ConsensusShutdownTimeout                *string  `yaml:"consensus-shutdown-timeout,omitempty"`
	ConsensusGossipFrequency                *string  `yaml:"consensus-gossip-frequency,omitempty"`
	MinDelegationFee                        *int     `yaml:"min-delegation-fee,omitempty"`
	MinValidatorStake                       *int     `yaml:"min-validator-stake,omitempty"`
	MaxStakeDuration                        *string  `yaml:"max-stake-duration,omitempty"`
	MaxValidatorStake                       *int     `yaml:"max-validator-stake,omitempty"`
	StakeMintingPeriod                      *string  `yaml:"stake-minting-period,omitempty"`
	NetworkInitialTimeout                   *string  `yaml:"network-initial-timeout,omitempty"`
	NetworkMinimumTimeout                   *string  `yaml:"network-minimum-timeout,omitempty"`
	NetworkMaximumTimeout                   *string  `yaml:"network-maximum-timeout,omitempty"`
	NetworkHealthMaxSendFailRateKey         *float64 `yaml:"network-health-max-send-fail-rate,omitempty"`
	NetworkHealthMaxPortionSendQueueFillKey *float64 `yaml:"network-health-max-portion-send-queue-full"`
	NetworkHealthMaxTimeSinceMsgSentKey     *string  `yaml:"network-health-max-time-since-msg-sent,omitempty"`
	NetworkHealthMaxTimeSinceMsgReceivedKey *string  `yaml:"network-health-max-time-since-msg-received,omitempty"`
	NetworkHealthMinConnPeers               *int     `yaml:"network-health-min-conn-peers,omitempty"`
	NetworkTimeoutCoefficient               *int     `yaml:"network-timeout-coefficient,omitempty"`
	NetworkTimeoutHalflife                  *string  `yaml:"network-timeout-halflife,omitempty"`
	NetworkPeerListGossipFrequency          *string  `yaml:"network-peer-list-gossip-frequency,omitempty"`
	NetworkPeerListGossipSize               *int     `yaml:"network-peer-list-gossip-size,omitempty"`
	NetworkPeerListSize                     *int     `yaml:"network-peer-list-size,omitempty"`
	StakingEnabled                          *bool    `yaml:"staking-enabled,omitempty"`
	StakingPort                             *uint    `yaml:"staking-port,omitempty"`
	StakingDisabledWeight                   *int     `yaml:"staking-disabled-weight,omitempty"`
	StakingTLSKeyFile                       *string  `yaml:"staking-tls-key-file,omitempty"`
	StakingTLSCertFile                      *string  `yaml:"staking-tls-cert-file,omitempty"`
	APIAuthRequired                         *bool    `yaml:"api-auth-required,omitempty"`
	APIAuthPasswordFileKey                  *string  `yaml:"api-auth-password-file,omitempty"`
	MinStakeDuration                        *string  `yaml:"min-stake-duration,omitempty"`
	WhitelistedSubnets                      *string  `yaml:"whitelisted-subnets,omitempty"`
	APIHealthEnabled                        *bool    `yaml:"api-health-enabled,omitempty"`
	ConfigFile                              *string  `yaml:"config-file,omitempty"`
	APIInfoEnabled                          *bool    `yaml:"api-info-enabled,omitempty"`
	NetworkCompressionEnabled               *bool    `yaml:"network-compression-enabled,omitempty"`
	IPCSChainIDs                            *string  `yaml:"ipcs-chain-ids,omitempty"`
	IPCSPath                                *string  `yaml:"ipcs-path,omitempty"`
	FDLimit                                 *int     `yaml:"fd-limit,omitempty"`
	BenchlistDuration                       *string  `yaml:"benchlist-duration,omitempty"`
	BenchlistFailThreshold                  *int     `yaml:"benchlist-fail-threshold,omitempty"`
	BenchlistMinFailingDuration             *string  `yaml:"benchlist-min-failing-duration,omitempty"`
	BenchlistPeerSummaryEnabled             *bool    `yaml:"benchlist-peer-summary-enabled,omitempty"`
	UptimeRequirement                       *float64 `yaml:"uptime-requirement,omitempty"`
	RetryBootstrapWarnFrequency             *uint    `yaml:"bootstrap-retry-warn-frequency,omitempty"`
	RetryBootstrap                          *bool    `yaml:"bootstrap-retry-enabled,omitempty"`
	HealthCheckAveragerHalflifeKey          *string  `yaml:"health-check-averager-halflife,omitempty"`
	HealthCheckFreqKey                      *string  `yaml:"health-check-frequency,omitempty"`
	RouterHealthMaxOutstandingRequestsKey   *int     `yaml:"router-health-max-outstanding-requests,omitempty"`
	RouterHealthMaxDropRateKey              *float64 `yaml:"router-health-max-drop-rate,omitempty"`
	IndexEnabled                            *bool    `yaml:"index-enabled,omitempty"`
	PluginModeEnabled                       *bool    `yaml:"plugin-mode-enabled,omitempty"`
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
		ClientLocation:                          "",
		Meta:                                    "",
		DataDir:                                 "",
		AssertionsEnabled:                       true,
		Version:                                 false,
		TxFee:                                   1000000,
		PublicIP:                                "127.0.0.1",
		DynamicUpdateDuration:                   "5m",
		DynamicPublicIP:                         "",
		NetworkID:                               "local",
		SignatureVerificationEnabled:            true,
		APIAdminEnabled:                         true,
		APIIPCsEnabled:                          true,
		APIKeystoreEnabled:                      true,
		APIMetricsEnabled:                       true,
		HTTPHost:                                "127.0.0.1",
		HTTPPort:                                9650,
		HTTPTLSEnabled:                          false,
		HTTPTLSCertFile:                         "",
		HTTPTLSKeyFile:                          "",
		BootstrapIPs:                            "",
		BootstrapIDs:                            "",
		BootstrapBeaconConnectionTimeout:        "60s",
		DBType:                                  "memdb",
		DBDir:                                   "db",
		BuildDir:                                "",
		PluginDir:                               "",
		LogLevel:                                "info",
		LogDir:                                  "logs",
		LogDisplayLevel:                         "", // defaults to the value provided to --log-level
		LogDisplayHighlight:                     "colors",
		SnowAvalancheBatchSize:                  30,
		SnowAvalancheNumParents:                 5,
		SnowSampleSize:                          20,
		SnowQuorumSize:                          16,
		SnowVirtuousCommitThreshold:             15,
		SnowRogueCommitThreshold:                20,
		SnowConcurrentRepolls:                   4,
		MinDelegatorStake:                       5000000,
		ConsensusShutdownTimeout:                "5s",
		ConsensusGossipFrequency:                "10s",
		MinDelegationFee:                        20000,
		MinValidatorStake:                       5000000,
		MaxStakeDuration:                        "8760h",
		MaxValidatorStake:                       3000000000000000,
		StakeMintingPeriod:                      "8760h",
		NetworkInitialTimeout:                   "5s",
		NetworkMinimumTimeout:                   "5s",
		NetworkMaximumTimeout:                   "10s",
		NetworkHealthMaxSendFailRateKey:         0.9,
		NetworkHealthMaxPortionSendQueueFillKey: 0.9,
		NetworkHealthMaxTimeSinceMsgSentKey:     "1m",
		NetworkHealthMaxTimeSinceMsgReceivedKey: "1m",
		NetworkHealthMinConnPeers:               1,
		NetworkTimeoutCoefficient:               2,
		NetworkTimeoutHalflife:                  "5m",
		NetworkPeerListGossipFrequency:          "1m",
		NetworkPeerListGossipSize:               50,
		NetworkPeerListSize:                     20,
		StakingEnabled:                          false,
		StakingPort:                             9651,
		StakingDisabledWeight:                   1,
		StakingTLSKeyFile:                       "",
		StakingTLSCertFile:                      "",
		APIAuthRequired:                         false,
		APIAuthPasswordFileKey:                  "",
		MinStakeDuration:                        "336h",
		APIHealthEnabled:                        true,
		ConfigFile:                              "",
		WhitelistedSubnets:                      "",
		APIInfoEnabled:                          true,
		NetworkCompressionEnabled:               true,
		IPCSChainIDs:                            "",
		IPCSPath:                                "/tmp",
		FDLimit:                                 32768,
		BenchlistDuration:                       "1h",
		BenchlistFailThreshold:                  10,
		BenchlistMinFailingDuration:             "5m",
		BenchlistPeerSummaryEnabled:             false,
		UptimeRequirement:                       0.6,
		RetryBootstrapWarnFrequency:             50,
		RetryBootstrap:                          true,
		HealthCheckAveragerHalflifeKey:          "10s",
		HealthCheckFreqKey:                      "30s",
		RouterHealthMaxOutstandingRequestsKey:   1024,
		RouterHealthMaxDropRateKey:              1,
		IndexEnabled:                            false,
		PluginModeEnabled:                       false,
		MeterVMsEnabled:                         false,
	}
}
