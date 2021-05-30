package node

import (
	"os"
	"path"
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
	P2PTLSEnabled                bool

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
	BootstrapIPs string
	BootstrapIDs string

	// Database
	DBEnabled bool
	DBDir     string

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
	SnowEpochFirstTransition    int
	SnowEpochDuration           string
	SnowConcurrentRepolls       int
	MinDelegatorStake           int
	ConsensusShutdownTimeout    string
	ConsensusGossipFrequency    string
	MinDelegationFee            int
	MinValidatorStake           int
	MaxStakeDuration            string
	MaxValidatorStake           int
	CreationTxFee               int

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

	// Connection
	ConnMeterMaxConns      int
	ConnMeterResetDuration string

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

	// Message Handling
	MaxNonStakerPendingMsgs int

	// Network Timeout
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

	// Restart on Disconnect
	RestartOnDisconnected      bool
	DisconnectedCheckFrequency string
	DisconnectedRestartTimeout string

	// Uptime Requirement
	UptimeRequirement float64

	// Retry
	RetryBootstrapMaxAttempts int
	RetryBootstrap            bool

	// Health
	HealthCheckAveragerHalflifeKey string
	HealthCheckFreqKey             string

	// Router
	RouterHealthMaxOutstandingRequestsKey int
	RouterHealthMaxDropRateKey            float64

	IndexEnabled bool
}

// FlagsYAML mimics Flags but uses pointers for proper YAML interpretation
// Note: FlagsYAML and Flags must always be identical in fields, otherwise parsing will break
type FlagsYAML struct {
	// Avash metadata
	ClientLocation *string  `yaml:"-"`
	Meta           *string  `yaml:"-"`
	DataDir        *string  `yaml:"-"`

	// Assertions
	AssertionsEnabled *bool    `yaml:"assertions-enabled,omitempty"`

	// Version
	Version *bool    `yaml:"version,omitempty"`

	// TX fees
	TxFee *uint    `yaml:"tx-fee,omitempty"`

	// IP
	PublicIP              *string  `yaml:"public-ip,omitempty"`
	DynamicUpdateDuration *string  `yaml:"dynamic-update-duration,omitempty"` //Added
	DynamicPublicIP       *string  `yaml:"dynamic-public-ip,omitempty"`

	// Network ID
	NetworkID *string  `yaml:"network-id,omitempty"`

	// Crypto
	SignatureVerificationEnabled *bool    `yaml:"signature-verification-enabled,omitempty"`
	P2PTLSEnabled                *bool    `yaml:"p2p-tls-enabled,omitempty"` //Added

	// APIs
	APIAdminEnabled    *bool    `yaml:"api-admin-enabled,omitempty"`
	APIIPCsEnabled     *bool    `yaml:"api-ipcs-enabled,omitempty"`
	APIKeystoreEnabled *bool    `yaml:"api-keystore-enabled,omitempty"`
	APIMetricsEnabled  *bool    `yaml:"api-metrics-enabled,omitempty"`
	APIHealthEnabled   *bool    `yaml:"api-health-enabled,omitempty"`
	APIInfoEnabled     *bool    `yaml:"api-info-enabled,omitempty"`

	// HTTP
	HTTPHost        *string  `yaml:"http-host,omitempty"`
	HTTPPort        *uint    `yaml:"http-port,omitempty"`
	HTTPTLSEnabled  *bool    `yaml:"http-tls-enabled,omitempty"`
	HTTPTLSCertFile *string  `yaml:"http-tls-cert-file,omitempty"`
	HTTPTLSKeyFile  *string  `yaml:"http-tls-key-file,omitempty"`

	// Bootstrapping
	BootstrapIPs *string  `yaml:"bootstrap-ips,omitempty"`
	BootstrapIDs *string  `yaml:"bootstrap-ids,omitempty"`

	// Database
	DBEnabled *bool    `yaml:"db-enabled,omitempty"`
	DBDir     *string  `yaml:"db-dir,omitempty"`

	// Plugins
	PluginDir  *string  `yaml:"plugin-dir,omitempty"`

	// Logging
	LogLevel            *string  `yaml:"log-level,omitempty"`
	LogDir              *string  `yaml:"log-dir,omitempty"`
	LogDisplayLevel     *string  `yaml:"log-display-level,omitempty"`
	LogDisplayHighlight *string  `yaml:"log-display-highlight,omitempty"`

	// Consensus
	SnowAvalancheBatchSize      *int     `yaml:"snow-avalanche-batch-size,omitempty"`
	SnowAvalancheNumParents     *int     `yaml:"snow-avalanche-num-parents,omitempty"`
	SnowSampleSize              *int     `yaml:"snow-sample-size,omitempty"`
	SnowQuorumSize              *int     `yaml:"snow-quorum-size,omitempty"`
	SnowVirtuousCommitThreshold *int     `yaml:"snow-virtuous-commit-threshold,omitempty"`
	SnowRogueCommitThreshold    *int     `yaml:"snow-rogue-commit-threshold,omitempty"`
	SnowEpochFirstTransition    *int     `yaml:"snow-epoch-first-transition,omitempty"`
	SnowEpochDuration           *string  `yaml:"snow-epoch-duration,omitempty"`
	SnowConcurrentRepolls       *int     `yaml:"snow-concurrent-repolls,omitempty"`
	MinDelegatorStake           *int     `yaml:"min-delegator-stake,omitempty"`
	ConsensusShutdownTimeout    *string  `yaml:"consensus-shutdown-timeout,omitempty"`
	ConsensusGossipFrequency    *string  `yaml:"consensus-gossip-frequency,omitempty"`
	MinDelegationFee            *int     `yaml:"min-delegation-fee,omitempty"`
	MinValidatorStake           *int     `yaml:"min-validator-stake,omitempty"`
	MaxStakeDuration            *string  `yaml:"max-stake-duration,omitempty"`
	MaxValidatorStake           *int     `yaml:"max-validator-stake,omitempty"`
	CreationTxFee               *int     `yaml:"creation-tx-fee,omitempty"`

	// Staking
	StakingEnabled        *bool    `yaml:"staking-enabled,omitempty"`
	StakeMintingPeriod    *string  `yaml:"stake-minting-period,omitempty"`
	StakingPort           *uint    `yaml:"staking-port,omitempty"`
	StakingDisabledWeight *int     `yaml:"staking-disabled-weight,omitempty"`
	StakingTLSKeyFile     *string  `yaml:"staking-tls-key-file,omitempty"`
	StakingTLSCertFile    *string  `yaml:"staking-tls-cert-file,omitempty"`

	// Auth
	APIAuthRequired        *bool    `yaml:"api-auth-required,omitempty"`
	APIAuthPasswordFileKey *string  `yaml:"api-auth-password-file,omitempty"`
	MinStakeDuration       *string  `yaml:"min-stake-duration,omitempty"`

	// Whitelisted Subnets
	WhitelistedSubnets *string  `yaml:"whitelisted-subnets,omitempty"`

	// Config
	ConfigFile *string  `yaml:"config-file,omitempty"`

	// Connection
	ConnMeterMaxConns      *int     `yaml:"conn-meter-max-conns,omitempty"`
	ConnMeterResetDuration *string  `yaml:"conn-meter-reset-duration,omitempty"`

	// IPCS
	IPCSChainIDs *string  `yaml:"ipcs-chain-ids,omitempty"`
	IPCSPath     *string  `yaml:"ipcs-path,omitempty"`

	// File Descriptor Limit
	FDLimit *int     `yaml:"fd-limit,omitempty"`

	// Benchlist
	BenchlistFailThreshold      *int     `yaml:"benchlist-fail-threshold,omitempty"`
	BenchlistMinFailingDuration *string  `yaml:"benchlist-min-failing-duration,omitempty"`
	BenchlistPeerSummaryEnabled *bool    `yaml:"benchlist-peer-summary-enabled,omitempty"`
	BenchlistDuration           *string  `yaml:"benchlist-duration,omitempty"`

	// Message Handling
	MaxNonStakerPendingMsgs *int     `yaml:"max-non-staker-pending-msgs,omitempty"`

	// Network Timeout
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

	// Restart on Disconnect
	RestartOnDisconnected      *bool    `yaml:"restart-on-disconnected,omitempty"`
	DisconnectedCheckFrequency *string  `yaml:"disconnected-check-frequency,omitempty"`
	DisconnectedRestartTimeout  *string  `yaml:"disconnected-restart-timeout,omitempty"`

	// Uptime Requirement
	UptimeRequirement *float64 `yaml:"uptime-requirement,omitempty"`

	// Retry
	RetryBootstrapMaxAttempts *uint    `yaml:"bootstrap-retry-max-attempts,omitempty"`
	RetryBootstrap            *bool    `yaml:"bootstrap-retry-enabled,omitempty"`

	// Health
	HealthCheckAveragerHalflifeKey *string  `yaml:"health-check-averager-halflife,omitempty"`
	HealthCheckFreqKey             *string  `yaml:"health-check-frequency,omitempty"`

	// Router
	RouterHealthMaxOutstandingRequestsKey *int     `yaml:"router-health-max-outstanding-requests,omitempty"`
	RouterHealthMaxDropRateKey             *float64 `yaml:"router-health-max-drop-rate,omitempty"`

	IndexEnabled *bool    `yaml:"index-enabled,omitempty"`
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
	for i := 0; i < f.NumField(); i++ {
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
		DBEnabled:                               true,
		DBDir:                                   "db",
		PluginDir:                               path.Join(os.Getenv("GOPATH"), "src", "github.com", "ava-labs", "avalanchego", "build", "plugins"),
		LogLevel:                                "info",
		LogDir:                                  "logs",
		LogDisplayLevel:                         "", // defaults to the value provided to --log-level
		LogDisplayHighlight:                     "colors",
		SnowAvalancheBatchSize:                  30,
		SnowAvalancheNumParents:                 5,
		SnowSampleSize:                          2,
		SnowQuorumSize:                          2,
		SnowVirtuousCommitThreshold:             5,
		SnowRogueCommitThreshold:                10,
		SnowEpochFirstTransition:                1609873200,
		SnowEpochDuration:                       "6h",
		SnowConcurrentRepolls:                   4,
		MinDelegatorStake:                       5000000,
		ConsensusShutdownTimeout:                "5s",
		ConsensusGossipFrequency:                "10s",
		MinDelegationFee:                        20000,
		MinValidatorStake:                       5000000,
		MaxStakeDuration:                        "8760h",
		MaxValidatorStake:                       3000000000000000,
		StakeMintingPeriod:                      "8760h",
		CreationTxFee:                           1000000,
		MaxNonStakerPendingMsgs:                 20,
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
		P2PTLSEnabled:                           true,
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
		ConnMeterMaxConns:                       5,
		ConnMeterResetDuration:                  "",
		IPCSChainIDs:                            "",
		IPCSPath:                                "/tmp",
		FDLimit:                                 32768,
		BenchlistDuration:                       "1h",
		BenchlistFailThreshold:                  10,
		BenchlistMinFailingDuration:             "5m",
		BenchlistPeerSummaryEnabled:             false,
		RestartOnDisconnected:                   true,
		DisconnectedCheckFrequency:              "10s",
		DisconnectedRestartTimeout:              "1m",
		UptimeRequirement:                       0.6,
		RetryBootstrapMaxAttempts:               50,
		RetryBootstrap:                          true,
		HealthCheckAveragerHalflifeKey:          "10s",
		HealthCheckFreqKey:                      "30s",
		RouterHealthMaxOutstandingRequestsKey:   1024,
		RouterHealthMaxDropRateKey:              1,
		IndexEnabled:                            false,
	}
}
