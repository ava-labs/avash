package node

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FlagsToArgs converts a `Flags` struct into a CLI command flag string
func FlagsToArgs(flags Flags, basedir string, sepBase bool) ([]string, Metadata) {
	// Port targets
	httpPortString := strconv.FormatUint(uint64(flags.HTTPPort), 10)
	stakingPortString := strconv.FormatUint(uint64(flags.StakingPort), 10)

	// Paths/directories
	dbPath := basedir + "/" + flags.DBDir
	dbType := flags.DBType
	dataPath := basedir + "/" + flags.DataDir
	logPath := basedir + "/" + flags.LogDir
	if sepBase {
		dbPath = flags.DBDir
		dataPath = flags.DataDir
		logPath = flags.LogDir
	}

	wd, _ := os.Getwd()
	// If the path given in the flag doesn't begin with "/", treat it as relative
	// to the directory of the avash binary
	httpCertFile := flags.HTTPTLSCertFile
	if httpCertFile != "" && string(httpCertFile[0]) != "/" && !sepBase {
		httpCertFile = fmt.Sprintf("%s/%s", wd, httpCertFile)
	}

	httpKeyFile := flags.HTTPTLSKeyFile
	if httpKeyFile != "" && string(httpKeyFile[0]) != "/" && !sepBase {
		httpKeyFile = fmt.Sprintf("%s/%s", wd, httpKeyFile)
	}

	stakerCertFile := flags.StakingTLSCertFile
	if stakerCertFile != "" && string(stakerCertFile[0]) != "/" && !sepBase {
		stakerCertFile = fmt.Sprintf("%s/%s", wd, stakerCertFile)
	}

	stakerKeyFile := flags.StakingTLSKeyFile
	if stakerKeyFile != "" && string(stakerKeyFile[0]) != "/" && !sepBase {
		stakerKeyFile = fmt.Sprintf("%s/%s", wd, stakerKeyFile)
	}

	args := []string{
		"--assertions-enabled=" + strconv.FormatBool(flags.AssertionsEnabled),
		"--version=" + strconv.FormatBool(flags.Version),
		"--tx-fee=" + strconv.FormatUint(uint64(flags.TxFee), 10),
		"--public-ip=" + flags.PublicIP,
		"--dynamic-update-duration=" + flags.DynamicUpdateDuration,
		"--dynamic-public-ip=" + flags.DynamicPublicIP,
		"--network-id=" + flags.NetworkID,
		"--signature-verification-enabled=" + strconv.FormatBool(flags.SignatureVerificationEnabled),
		"--api-admin-enabled=" + strconv.FormatBool(flags.APIAdminEnabled),
		"--api-ipcs-enabled=" + strconv.FormatBool(flags.APIIPCsEnabled),
		"--api-keystore-enabled=" + strconv.FormatBool(flags.APIKeystoreEnabled),
		"--api-metrics-enabled=" + strconv.FormatBool(flags.APIMetricsEnabled),
		"--http-host=" + flags.HTTPHost,
		"--http-port=" + httpPortString,
		"--http-tls-enabled=" + strconv.FormatBool(flags.HTTPTLSEnabled),
		"--http-tls-cert-file=" + httpCertFile,
		"--http-tls-key-file=" + httpKeyFile,
		"--bootstrap-ips=" + flags.BootstrapIPs,
		"--bootstrap-ids=" + flags.BootstrapIDs,
		"--bootstrap-beacon-connection-timeout=" + flags.BootstrapBeaconConnectionTimeout,
		"--db-type=" + dbType,
		"--db-dir=" + dbPath,
		"--plugin-dir=" + flags.PluginDir,
		"--build-dir=" + flags.BuildDir,
		"--log-level=" + flags.LogLevel,
		"--log-dir=" + logPath,
		"--log-display-level=" + flags.LogDisplayLevel,
		"--log-display-highlight=" + flags.LogDisplayHighlight,
		"--snow-avalanche-batch-size=" + strconv.Itoa(flags.SnowAvalancheBatchSize),
		"--snow-avalanche-num-parents=" + strconv.Itoa(flags.SnowAvalancheNumParents),
		"--snow-sample-size=" + strconv.Itoa(flags.SnowSampleSize),
		"--snow-quorum-size=" + strconv.Itoa(flags.SnowQuorumSize),
		"--snow-virtuous-commit-threshold=" + strconv.Itoa(flags.SnowVirtuousCommitThreshold),
		"--min-delegator-stake=" + strconv.Itoa(flags.MinDelegatorStake),
		"--consensus-shutdown-timeout=" + flags.ConsensusShutdownTimeout,
		"--consensus-gossip-frequency=" + flags.ConsensusGossipFrequency,
		"--min-delegation-fee=" + strconv.Itoa(flags.MinDelegationFee),
		"--min-validator-stake=" + strconv.Itoa(flags.MinValidatorStake),
		"--max-stake-duration=" + flags.MaxStakeDuration,
		"--max-validator-stake=" + strconv.Itoa(flags.MaxValidatorStake),
		"--snow-concurrent-repolls=" + strconv.Itoa(flags.SnowConcurrentRepolls),
		"--stake-minting-period=" + flags.StakeMintingPeriod,
		"--network-initial-timeout=" + flags.NetworkInitialTimeout,
		"--network-minimum-timeout=" + flags.NetworkMinimumTimeout,
		"--network-maximum-timeout=" + flags.NetworkMaximumTimeout,
		fmt.Sprintf("--network-health-max-send-fail-rate=%f", flags.NetworkHealthMaxSendFailRateKey),
		fmt.Sprintf("--network-health-max-portion-send-queue-full=%f", flags.NetworkHealthMaxPortionSendQueueFillKey),
		"--network-health-max-time-since-msg-sent=" + flags.NetworkHealthMaxTimeSinceMsgSentKey,
		"--network-health-max-time-since-msg-received=" + flags.NetworkHealthMaxTimeSinceMsgReceivedKey,
		"--network-health-min-conn-peers=" + strconv.Itoa(flags.NetworkHealthMinConnPeers),
		"--network-timeout-coefficient=" + strconv.Itoa(flags.NetworkTimeoutCoefficient),
		"--network-timeout-halflife=" + flags.NetworkTimeoutHalflife,
		"--network-peer-list-gossip-frequency=" + flags.NetworkPeerListGossipFrequency,
		"--network-peer-list-gossip-size=" + strconv.Itoa(flags.NetworkPeerListGossipSize),
		"--network-peer-list-size=" + strconv.Itoa(flags.NetworkPeerListSize),
		"--staking-enabled=" + strconv.FormatBool(flags.StakingEnabled),
		"--staking-port=" + stakingPortString,
		"--staking-disabled-weight=" + strconv.Itoa(flags.StakingDisabledWeight),
		"--staking-tls-key-file=" + stakerKeyFile,
		"--staking-tls-cert-file=" + stakerCertFile,
		"--api-auth-required=" + strconv.FormatBool(flags.APIAuthRequired),
		"--api-auth-password-file=" + flags.APIAuthPasswordFileKey,
		"--min-stake-duration=" + flags.MinStakeDuration,
		"--whitelisted-subnets=" + flags.WhitelistedSubnets,
		"--api-health-enabled=" + strconv.FormatBool(flags.APIHealthEnabled),
		"--config-file=" + flags.ConfigFile,
		"--api-info-enabled=" + strconv.FormatBool(flags.APIInfoEnabled),
		"--network-compression-enabled=" + strconv.FormatBool(flags.NetworkCompressionEnabled),
		"--ipcs-chain-ids=" + flags.IPCSChainIDs,
		"--ipcs-path=" + flags.IPCSPath,
		"--fd-limit=" + strconv.Itoa(flags.FDLimit),
		"--benchlist-duration=" + flags.BenchlistDuration,
		"--benchlist-fail-threshold=" + strconv.Itoa(flags.BenchlistFailThreshold),
		"--benchlist-min-failing-duration=" + flags.BenchlistMinFailingDuration,
		"--benchlist-peer-summary-enabled=" + strconv.FormatBool(flags.BenchlistPeerSummaryEnabled),
		fmt.Sprintf("--uptime-requirement=%f", flags.UptimeRequirement),
		"--bootstrap-retry-enabled=" + strconv.FormatBool(flags.RetryBootstrap),
		"--health-check-averager-halflife=" + flags.HealthCheckAveragerHalflifeKey,
		"--health-check-frequency=" + flags.HealthCheckFreqKey,
		"--router-health-max-outstanding-requests=" + strconv.Itoa(flags.RouterHealthMaxOutstandingRequestsKey),
		fmt.Sprintf("--router-health-max-drop-rate=%f", flags.RouterHealthMaxDropRateKey),
		"--index-enabled=" + strconv.FormatBool(flags.IndexEnabled),
		"--plugin-mode-enabled=" + strconv.FormatBool(flags.PluginModeEnabled),
		"--meter-vms-enabled=" + strconv.FormatBool(flags.MeterVMsEnabled),
	}
	if sepBase {
		args = append(args, "--data-dir="+basedir)
	}
	args = removeEmptyFlags(args)

	metadata := Metadata{
		Serverhost:     flags.PublicIP,
		Stakingport:    stakingPortString,
		HTTPport:       httpPortString,
		HTTPTLS:        flags.HTTPTLSEnabled,
		Dbdir:          dbPath,
		Datadir:        dataPath,
		Logsdir:        logPath,
		Loglevel:       flags.LogLevel,
		StakingEnabled: flags.StakingEnabled,
		StakerCertPath: stakerCertFile,
		StakerKeyPath:  stakerKeyFile,
	}

	return args, metadata
}

func removeEmptyFlags(args []string) []string {
	var res []string
	for _, f := range args {
		tmp := strings.TrimSpace(f)
		if !strings.HasSuffix(tmp, "=") {
			res = append(res, tmp)
		}
	}
	return res
}
