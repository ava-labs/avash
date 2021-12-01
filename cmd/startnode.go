// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kennygrant/sanitize"

	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/spf13/cobra"
)

var flags node.Flags

// StartnodeCmd represents the startnode command
var StartnodeCmd = &cobra.Command{
	Use:   "startnode [node name] args...",
	Short: "Starts a node process and gives it a name.",
	Long: `Starts an Avalanche client node using pmgo and gives it a name. Example:
	startnode MyNode1 --public-ip=127.0.0.1 --staking-port=9651 --http-port=9650 ... `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}
		log := cfg.Config.Log
		name := args[0]

		datadir := cfg.Config.DataDir
		basename := sanitize.BaseName(name)
		datapath := datadir + "/" + basename
		if basename == "" {
			log.Error("Process name can't be empty")
			return
		}

		err := validateConsensusArgs(
			flags.SnowSampleSize,
			flags.SnowQuorumSize,
			flags.SnowVirtuousCommitThreshold,
			flags.SnowRogueCommitThreshold,
		)
		if err != nil {
			log.Error(err.Error())
			return
		}

		args, md := node.FlagsToArgs(flags, sanitize.Path(datapath), false)
		defer func() {
			// Set flags to default for next `startnode` call
			flags = node.DefaultFlags()
		}()
		mdbytes, _ := json.MarshalIndent(md, " ", "    ")
		metadata := string(mdbytes)
		meta := flags.Meta
		if meta != "" {
			metadata = meta
		}
		avalancheLocation := flags.ClientLocation
		if avalancheLocation == "" {
			avalancheLocation = cfg.Config.AvalancheLocation
		}
		err = pmgr.ProcManager.AddProcess(avalancheLocation, "avalanche node", args, name, metadata, nil, nil, nil)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Created process %s.", name)
		pmgr.ProcManager.StartProcess(name)
	},
}

func validateConsensusArgs(k int, alpha int, beta1 int, beta2 int) error {
	rulesfailed := []string(nil)
	if k <= 0 {
		rulesfailed = append(rulesfailed, "k > 0")
	}
	if alpha > k {
		rulesfailed = append(rulesfailed, "alpha <= k")
	}
	if (k / 2) >= alpha {
		rulesfailed = append(rulesfailed, "alpha > floor(k/2)")
	}
	if beta1 <= 0 {
		rulesfailed = append(rulesfailed, "beta1 > 0")
	}
	if beta1 > beta2 {
		rulesfailed = append(rulesfailed, "beta2 >= beta1")
	}
	if len(rulesfailed) == 0 {
		return nil
	}
	return errors.New("Invalid consensus params: \n" + strings.Join(rulesfailed, "\n"))
}

func init() {
	flags = node.DefaultFlags()
	StartnodeCmd.Flags().StringVar(&flags.ClientLocation, "client-location", flags.ClientLocation, "Path to AVA node client, defaulting to the config file's value.")
	StartnodeCmd.Flags().StringVar(&flags.Meta, "meta", flags.Meta, "Override default metadata for the node process.")
	StartnodeCmd.Flags().StringVar(&flags.DataDir, "data-dir", flags.DataDir, "Name of directory for the data stash.")

	StartnodeCmd.Flags().BoolVar(&flags.AssertionsEnabled, "assertions-enabled", flags.AssertionsEnabled, "Turn on assertion execution.")
	StartnodeCmd.Flags().BoolVar(&flags.Version, "version", flags.Version, "If this is `true`, print the version and quit. Defaults to `false`")
	StartnodeCmd.Flags().UintVar(&flags.TxFee, "tx-fee", flags.TxFee, "Transaction fee, in $nAVAX.")

	StartnodeCmd.Flags().StringVar(&flags.PluginDir, "plugin-dir", flags.PluginDir, "Directory to search for plugins")

	StartnodeCmd.Flags().BoolVar(&flags.APIAdminEnabled, "api-admin-enabled", flags.APIAdminEnabled, "If true, this node exposes the Admin API")
	StartnodeCmd.Flags().BoolVar(&flags.APIKeystoreEnabled, "api-keystore-enabled", flags.APIKeystoreEnabled, "If true, this node exposes the Keystore API")
	StartnodeCmd.Flags().BoolVar(&flags.APIMetricsEnabled, "api-metrics-enabled", flags.APIMetricsEnabled, "If true, this node exposes the Metrics API")
	StartnodeCmd.Flags().BoolVar(&flags.APIIPCsEnabled, "api-ipcs-enabled", flags.APIIPCsEnabled, "If true, IPCs can be opened")
	StartnodeCmd.Flags().BoolVar(&flags.APIHealthEnabled, "api-health-enabled", flags.APIHealthEnabled, "If set to `true`, this node will expose the Health API. Defaults to `true`")
	StartnodeCmd.Flags().BoolVar(&flags.APIInfoEnabled, "api-info-enabled", flags.APIInfoEnabled, "If set to `true`, this node will expose the Info API. Defaults to `true`")

	StartnodeCmd.Flags().StringVar(&flags.PublicIP, "public-ip", flags.PublicIP, "Public IP of this node.")
	StartnodeCmd.Flags().StringVar(&flags.DynamicUpdateDuration, "dynamic-update-duration", flags.DynamicUpdateDuration, "The time between poll events for `--dynamic-public-ip` or NAT traversal. The recommended minimum is 1 minute. Defaults to `5m`")
	StartnodeCmd.Flags().StringVar(&flags.DynamicPublicIP, "dynamic-public-ip", flags.DynamicPublicIP, "Valid values if param is present: `opendns`, `ifconfigco` or `ifconfigme`. This overrides `--public-ip`. If set, will poll the remote service every `--dynamic-update-duration` and update the node’s public IP address.")
	StartnodeCmd.Flags().StringVar(&flags.NetworkID, "network-id", flags.NetworkID, "Network ID this node will connect to.")
	StartnodeCmd.Flags().BoolVar(&flags.SignatureVerificationEnabled, "signature-verification-enabled", flags.SignatureVerificationEnabled, "Turn on signature verification.")

	StartnodeCmd.Flags().StringVar(&flags.HTTPHost, "http-host", flags.HTTPHost, "The address that HTTP APIs listen on.")
	StartnodeCmd.Flags().UintVar(&flags.HTTPPort, "http-port", flags.HTTPPort, "Port of the HTTP server.")
	StartnodeCmd.Flags().BoolVar(&flags.HTTPTLSEnabled, "http-tls-enabled", flags.HTTPTLSEnabled, "Upgrade the HTTP server to HTTPS.")
	StartnodeCmd.Flags().StringVar(&flags.HTTPTLSCertFile, "http-tls-cert-file", flags.HTTPTLSCertFile, "TLS certificate file for the HTTPS server.")
	StartnodeCmd.Flags().StringVar(&flags.HTTPTLSKeyFile, "http-tls-key-file", flags.HTTPTLSKeyFile, "TLS private key file for the HTTPS server.")

	StartnodeCmd.Flags().StringVar(&flags.BootstrapIPs, "bootstrap-ips", flags.BootstrapIPs, "Comma separated list of bootstrap nodes to connect to. Example: 127.0.0.1:9630,127.0.0.1:9620")
	StartnodeCmd.Flags().StringVar(&flags.BootstrapIDs, "bootstrap-ids", flags.BootstrapIDs, "Comma separated list of bootstrap peer ids to connect to. Example: NodeID-JR4dVmy6ffUGAKCBDkyCbeZbyHQBeDsET,NodeID-8CrVPQZ4VSqgL8zTdvL14G8HqAfrBr4z")
	StartnodeCmd.Flags().StringVar(&flags.BootstrapBeaconConnectionTimeout, "bootstrap-beacon-connection-timeout", flags.BootstrapBeaconConnectionTimeout, "Timeout when attempting to connect to bootstrapping beacons.")

	StartnodeCmd.Flags().StringVar(&flags.DBType, "db-type", flags.DBType, "Type of the DB to use (memdb|leveldb|rocksdb)")
	StartnodeCmd.Flags().StringVar(&flags.DBDir, "db-dir", flags.DBDir, "Database directory for Avalanche state.")

	StartnodeCmd.Flags().StringVar(&flags.BuildDir, "build-dir", flags.BuildDir, "Path to the build directory.")

	StartnodeCmd.Flags().StringVar(&flags.LogLevel, "log-level", flags.LogLevel, "Specify the log level. Should be one of {verbo, debug, info, warn, error, fatal, off}")
	StartnodeCmd.Flags().StringVar(&flags.LogDir, "log-dir", flags.LogDir, "Name of directory for the node's logging.")
	StartnodeCmd.Flags().StringVar(&flags.LogDisplayLevel, "log-display-level", flags.LogDisplayLevel, "{Off, Fatal, Error, Warn, Info, Debug, Verbo}. The log level determines which events to display to the screen. If left blank, will default to the value provided to `--log-level`")
	StartnodeCmd.Flags().StringVar(&flags.LogDisplayHighlight, "log-display-highlight", flags.LogDisplayHighlight, "Whether to color/highlight display logs. Default highlights when the output is a terminal. Otherwise, should be one of {auto, plain, colors}")

	StartnodeCmd.Flags().IntVar(&flags.SnowAvalancheBatchSize, "snow-avalanche-batch-size", flags.SnowAvalancheBatchSize, "Number of operations to batch in each new vertex.")
	StartnodeCmd.Flags().IntVar(&flags.SnowAvalancheNumParents, "snow-avalanche-num-parents", flags.SnowAvalancheNumParents, "Number of vertexes for reference from each new vertex.")
	StartnodeCmd.Flags().IntVar(&flags.SnowSampleSize, "snow-sample-size", flags.SnowSampleSize, "Number of nodes to query for each network poll.")
	StartnodeCmd.Flags().IntVar(&flags.SnowQuorumSize, "snow-quorum-size", flags.SnowQuorumSize, "Alpha value to use for required number positive results.")
	StartnodeCmd.Flags().IntVar(&flags.SnowVirtuousCommitThreshold, "snow-virtuous-commit-threshold", flags.SnowVirtuousCommitThreshold, "Beta value to use for virtuous transactions.")
	StartnodeCmd.Flags().IntVar(&flags.SnowRogueCommitThreshold, "snow-rogue-commit-threshold", flags.SnowRogueCommitThreshold, "Beta value to use for rogue transactions.")
	StartnodeCmd.Flags().IntVar(&flags.SnowConcurrentRepolls, "snow-concurrent-repolls", flags.SnowConcurrentRepolls, "Snow consensus requires repolling transactions that are issued during low time of network usage. This parameter lets one define how aggressive the client will be in finalizing these pending transactions. This should only be changed after careful consideration of the tradeoffs of Snow consensus. The value must be at least `1` and at most `--snow-rogue-commit-threshold`. Defaults to `4`")
	StartnodeCmd.Flags().IntVar(&flags.MinDelegatorStake, "min-delegator-stake", flags.MinDelegatorStake, "The minimum stake, in nAVAX, that can be delegated to a validator of the Primary Network. Defaults to `25000000000` (25 AVAX) on Main Net. Defaults to `5000000` (.005 AVAX) on Test Net.")
	StartnodeCmd.Flags().StringVar(&flags.ConsensusShutdownTimeout, "consensus-shutdown-timeout", flags.ConsensusShutdownTimeout, "Timeout before killing an unresponsive chain. Defaults to `5s`")
	StartnodeCmd.Flags().StringVar(&flags.ConsensusGossipFrequency, "consensus-gossip-frequency", flags.ConsensusGossipFrequency, "Time between gossiping accepted frontiers. Defaults to `10s`")
	StartnodeCmd.Flags().IntVar(&flags.MinDelegationFee, "min-delegation-fee", flags.MinDelegationFee, "The minimum delegation fee that can be charged for delegation on the Primary Network, multiplied by `10,000` . Must be in the range `[0, 1000000]`. Defaults to `20000` (2%) on Main Net.")
	StartnodeCmd.Flags().IntVar(&flags.MinValidatorStake, "min-validator-stake", flags.MinValidatorStake, "The minimum stake, in nAVAX, required to validate the Primary Network. Defaults to `2000000000000` (2,000 AVAX) on Main Net. Defaults to `5000000` (.005 AVAX) on Test Net.")
	StartnodeCmd.Flags().StringVar(&flags.MaxStakeDuration, "max-stake-duration", flags.MaxStakeDuration, "The maximum staking duration, in seconds. Defaults to `8760h` (365 days) on Main Net.")
	StartnodeCmd.Flags().IntVar(&flags.MaxValidatorStake, "max-validator-stake", flags.MaxValidatorStake, "The maximum stake, in nAVAX, that can be placed on a validator on the primary network. Defaults to `3000000000000000` (3,000,000 AVAX) on Main Net. This includes stake provided by both the validator and by delegators to the validator.")
	StartnodeCmd.Flags().StringVar(&flags.StakeMintingPeriod, "stake-minting-period", flags.StakeMintingPeriod, "Consumption period of the staking function, in seconds. The Default on Main Net is `8760h` (365 days).")

	StartnodeCmd.Flags().BoolVar(&flags.StakingEnabled, "staking-enabled", flags.StakingEnabled, "Enable staking. If enabled, Network TLS is required.")
	StartnodeCmd.Flags().UintVar(&flags.StakingPort, "staking-port", flags.StakingPort, "Port of the consensus server.")
	StartnodeCmd.Flags().IntVar(&flags.StakingDisabledWeight, "staking-disabled-weight", flags.StakingDisabledWeight, "Weight to provide to each peer when staking is disabled. Defaults to `1`")
	StartnodeCmd.Flags().StringVar(&flags.StakingTLSCertFile, "staking-tls-cert-file", flags.StakingTLSCertFile, "TLS certificate file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.crt")
	StartnodeCmd.Flags().StringVar(&flags.StakingTLSKeyFile, "staking-tls-key-file", flags.StakingTLSKeyFile, "TLS private key file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.key")

	StartnodeCmd.Flags().BoolVar(&flags.APIAuthRequired, "api-auth-required", flags.APIAuthRequired, "If set to true, API calls require an authorization token. Defaults to `false`")
	StartnodeCmd.Flags().StringVar(&flags.APIAuthPasswordFileKey, "api-auth-password-file", flags.APIAuthPasswordFileKey, "Password file used to initially create/validate API authorization tokens. Can be changed via API call.")
	StartnodeCmd.Flags().StringVar(&flags.MinStakeDuration, "min-stake-duration", flags.MinStakeDuration, "Set the minimum staking duration. Ex: --min-stake-duration=5m")

	StartnodeCmd.Flags().StringVar(&flags.WhitelistedSubnets, "whitelisted-subnets", flags.WhitelistedSubnets, "Comma separated list of subnets that this node would validate if added to. Defaults to empty (will only validate the Primary Network)")

	StartnodeCmd.Flags().StringVar(&flags.ConfigFile, "config-file", flags.ConfigFile, "Config file specifies a JSON file to configure a node instead of specifying arguments via the command line. Command line arguments will override any options set in the config file.")

	StartnodeCmd.Flags().BoolVar(&flags.NetworkCompressionEnabled, "network-compression-enabled", flags.NetworkCompressionEnabled, "If true, compress Put, PushQuery, PeerList and Multiput messages sent to peers that support compression")

	StartnodeCmd.Flags().StringVar(&flags.IPCSChainIDs, "ipcs-chain-ids", flags.IPCSChainIDs, "Comma separated list of chain ids to connect to. There is no default value.")
	StartnodeCmd.Flags().StringVar(&flags.IPCSPath, "ipcs-path", flags.IPCSPath, "The directory (Unix) or named pipe prefix (Windows) for IPC sockets. Defaults to /tmp.")

	StartnodeCmd.Flags().IntVar(&flags.FDLimit, "fd-limit", flags.FDLimit, "Attempts to raise the process file descriptor limit to at least this value. Defaults to `32768`")

	StartnodeCmd.Flags().StringVar(&flags.BenchlistDuration, "benchlist-duration", flags.BenchlistDuration, "Amount of time a peer is benchlisted after surpassing `--benchlist-fail-threshold`. Defaults to `1h`")
	StartnodeCmd.Flags().IntVar(&flags.BenchlistFailThreshold, "benchlist-fail-threshold", flags.BenchlistFailThreshold, "Number of consecutive failed queries to a node before benching it (assuming all queries to it will fail). Defaults to `10`")
	StartnodeCmd.Flags().StringVar(&flags.BenchlistMinFailingDuration, "benchlist-min-failing-duration", flags.BenchlistMinFailingDuration, "Minimum amount of time messages to a peer must be failing before the peer is benched. Defaults to `5m`")
	StartnodeCmd.Flags().BoolVar(&flags.BenchlistPeerSummaryEnabled, "benchlist-peer-summary-enabled", flags.BenchlistPeerSummaryEnabled, "Enables peer specific query latency metrics. Defaults to `false`")

	StartnodeCmd.Flags().StringVar(&flags.NetworkInitialTimeout, "network-initial-timeout", flags.NetworkInitialTimeout, "Initial timeout value of the adaptive timeout manager, in nanoseconds. Defaults to `5s`")
	StartnodeCmd.Flags().StringVar(&flags.NetworkMinimumTimeout, "network-minimum-timeout", flags.NetworkMinimumTimeout, "Minimum timeout value of the adaptive timeout manager, in nanoseconds. Defaults to `5s`")
	StartnodeCmd.Flags().StringVar(&flags.NetworkMaximumTimeout, "network-maximum-timeout", flags.NetworkMaximumTimeout, "Maximum timeout value of the adaptive timeout manager, in nanoseconds. Defaults to `10s`")
	StartnodeCmd.Flags().Float64Var(&flags.NetworkHealthMaxSendFailRateKey, "network-health-max-send-fail-rate", flags.NetworkHealthMaxSendFailRateKey, "Network layer reports unhealthy if more than this portion of attempted message sends fail")
	StartnodeCmd.Flags().Float64Var(&flags.NetworkHealthMaxPortionSendQueueFillKey, "network-health-max-portion-send-queue-full", flags.NetworkHealthMaxPortionSendQueueFillKey, "Network layer returns unhealthy if more than this portion of the pending send queue is full")
	StartnodeCmd.Flags().StringVar(&flags.NetworkHealthMaxTimeSinceMsgSentKey, "network-health-max-time-since-msg-sent", flags.NetworkHealthMaxTimeSinceMsgSentKey, "Network layer returns unhealthy if haven't sent a message for at least this much time")
	StartnodeCmd.Flags().StringVar(&flags.NetworkHealthMaxTimeSinceMsgReceivedKey, "network-health-max-time-since-msg-received", flags.NetworkHealthMaxTimeSinceMsgReceivedKey, "Network layer returns unhealthy if haven't received a message for at least this much time")
	StartnodeCmd.Flags().IntVar(&flags.NetworkHealthMinConnPeers, "network-health-min-conn-peers", flags.NetworkHealthMinConnPeers, "Network layer returns unhealthy if connected to less than this many peers")
	StartnodeCmd.Flags().IntVar(&flags.NetworkTimeoutCoefficient, "network-timeout-coefficient", flags.NetworkTimeoutCoefficient, "Multiplied by average network response time to get the network timeout. Must be >= 1.")
	StartnodeCmd.Flags().StringVar(&flags.NetworkTimeoutHalflife, "network-timeout-halflife", flags.NetworkTimeoutHalflife, "Halflife of average network response time. Higher value --> network timeout is less volatile. Can't be 0.")
	gossipHelpMsg := fmt.Sprintf(
		"Gossip [%s] peers to [%s] peers every [%s]",
		"network-peer-list-size",
		"network-peer-list-gossip-size",
		"network-peer-list-gossip-frequency",
	)
	StartnodeCmd.Flags().StringVar(&flags.NetworkPeerListGossipFrequency, "network-peer-list-gossip-frequency", flags.NetworkPeerListGossipFrequency, gossipHelpMsg)
	StartnodeCmd.Flags().IntVar(&flags.NetworkPeerListGossipSize, "network-peer-list-gossip-size", flags.NetworkPeerListGossipSize, gossipHelpMsg)
	StartnodeCmd.Flags().IntVar(&flags.NetworkPeerListSize, "network-peer-list-size", flags.NetworkPeerListSize, gossipHelpMsg)

	StartnodeCmd.Flags().Float64Var(&flags.UptimeRequirement, "uptime-requirement", flags.UptimeRequirement, "Fraction of time a validator must be online to receive rewards. Defaults to `0.6`")

	StartnodeCmd.Flags().IntVar(&flags.RetryBootstrapWarnFrequency, "bootstrap-retry-warn-frequency", flags.RetryBootstrapWarnFrequency, "Specifies how many times bootstrap should be retried")
	StartnodeCmd.Flags().BoolVar(&flags.RetryBootstrap, "bootstrap-retry-enabled", flags.RetryBootstrap, "Specifies whether bootstrap should be retried")

	StartnodeCmd.Flags().StringVar(&flags.HealthCheckFreqKey, "health-check-frequency", flags.HealthCheckFreqKey, "Time between health checks")
	StartnodeCmd.Flags().StringVar(&flags.HealthCheckAveragerHalflifeKey, "health-check-averager-halflife", flags.HealthCheckAveragerHalflifeKey, "Halflife of averager when calculating a running average in a health check")

	StartnodeCmd.Flags().IntVar(&flags.RouterHealthMaxOutstandingRequestsKey, "router-health-max-outstanding-requests", flags.RouterHealthMaxOutstandingRequestsKey, "Node reports unhealthy if there are more than this many outstanding consensus requests (Get, PullQuery, etc.) over all chains")
	StartnodeCmd.Flags().Float64Var(&flags.RouterHealthMaxDropRateKey, "router-health-max-drop-rate", flags.RouterHealthMaxDropRateKey, "Node reports unhealthy if the router drops more than this portion of messages.")

	StartnodeCmd.Flags().BoolVar(&flags.IndexEnabled, "index-enabled", flags.IndexEnabled, "If true, index all accepted containers and transactions and expose them via an API")
	StartnodeCmd.Flags().BoolVar(&flags.PluginModeEnabled, "plugin-mode-enabled", flags.PluginModeEnabled, "Whether the app should run as a plugin. Defaults to false.")
	StartnodeCmd.Flags().BoolVar(&flags.MeterVMsEnabled, "meter-vms-enabled", flags.MeterVMsEnabled, "Whether the MeterVMs should be enabled on each VM.")
}
