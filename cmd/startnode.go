/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"encoding/json"
	"errors"
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
	Long: `Starts an ava client node using pmgo and gives it a name. Example:
	startnode MyNode1 --public-ip=127.0.0.1 --staking-port=5001 --http-port=9650 ... `,
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
		// Set flags to default for next `startnode` call
		flags = node.DefaultFlags()
		mdbytes, _ := json.MarshalIndent(md, " ", "    ")
		metadata := string(mdbytes)
		meta := flags.Meta
		if meta != "" {
			metadata = meta
		}
		avalocation := flags.ClientLocation
		if avalocation == "" {
			avalocation = cfg.Config.AvaLocation
		}
		err = pmgr.ProcManager.AddProcess(avalocation, "ava node", args, name, metadata, nil, nil, nil)
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
	StartnodeCmd.Flags().UintVar(&flags.AvaTxFee, "ava-tx-fee", flags.AvaTxFee, "Ava transaction fee, in $nAva.")

	StartnodeCmd.Flags().StringVar(&flags.PluginDir, "plugin-dir", flags.PluginDir, "Directory to search for plugins")

	StartnodeCmd.Flags().BoolVar(&flags.APIAdminEnabled, "api-admin-enabled", flags.APIAdminEnabled, "If true, this node exposes the Admin API")
	StartnodeCmd.Flags().BoolVar(&flags.APIKeystoreEnabled, "api-keystore-enabled", flags.APIKeystoreEnabled, "If true, this node exposes the Keystore API")
	StartnodeCmd.Flags().BoolVar(&flags.APIMetricsEnabled, "api-metrics-enabled", flags.APIMetricsEnabled, "If true, this node exposes the Metrics API")
	StartnodeCmd.Flags().BoolVar(&flags.APIIPCsEnabled, "api-ipcs-enabled", flags.APIIPCsEnabled, "If true, IPCs can be opened")

	StartnodeCmd.Flags().StringVar(&flags.PublicIP, "public-ip", flags.PublicIP, "Public IP of this node.")
	StartnodeCmd.Flags().StringVar(&flags.NetworkID, "network-id", flags.NetworkID, "Network ID this node will connect to.")
	StartnodeCmd.Flags().UintVar(&flags.XputServerPort, "xput-server-port", flags.XputServerPort, "Port of the deprecated throughput test server.")
	StartnodeCmd.Flags().BoolVar(&flags.XputServerEnabled, "xput-server-enabled", flags.XputServerEnabled, "If true, throughput test server is created.")
	StartnodeCmd.Flags().BoolVar(&flags.SignatureVerificationEnabled, "signature-verification-enabled", flags.SignatureVerificationEnabled, "Turn on signature verification.")

	StartnodeCmd.Flags().UintVar(&flags.HTTPPort, "http-port", flags.HTTPPort, "Port of the HTTP server.")
	StartnodeCmd.Flags().BoolVar(&flags.HTTPTLSEnabled, "http-tls-enabled", flags.HTTPTLSEnabled, "Upgrade the HTTP server to HTTPS.")
	StartnodeCmd.Flags().StringVar(&flags.HTTPTLSCertFile, "http-tls-cert-file", flags.HTTPTLSCertFile, "TLS certificate file for the HTTPS server.")
	StartnodeCmd.Flags().StringVar(&flags.HTTPTLSKeyFile, "http-tls-key-file", flags.HTTPTLSKeyFile, "TLS private key file for the HTTPS server.")

	StartnodeCmd.Flags().StringVar(&flags.BootstrapIPs, "bootstrap-ips", flags.BootstrapIPs, "Comma separated list of bootstrap nodes to connect to. Example: 127.0.0.1:9630,127.0.0.1:9620")
	StartnodeCmd.Flags().StringVar(&flags.BootstrapIDs, "bootstrap-ids", flags.BootstrapIDs, "Comma separated list of bootstrap peer ids to connect to. Example: JR4dVmy6ffUGAKCBDkyCbeZbyHQBeDsET,8CrVPQZ4VSqgL8zTdvL14G8HqAfrBr4z")

	StartnodeCmd.Flags().BoolVar(&flags.DBEnabled, "db-enabled", flags.DBEnabled, "Turn on persistent storage.")
	StartnodeCmd.Flags().StringVar(&flags.DBDir, "db-dir", flags.DBDir, "Database directory for Ava state.")

	StartnodeCmd.Flags().StringVar(&flags.LogLevel, "log-level", flags.LogLevel, "Specify the log level. Should be one of {verbo, debug, info, warn, error, fatal, off}")
	StartnodeCmd.Flags().StringVar(&flags.LogDir, "log-dir", flags.LogDir, "Name of directory for the node's logging.")

	StartnodeCmd.Flags().IntVar(&flags.SnowAvalancheBatchSize, "snow-avalanche-batch-size", flags.SnowAvalancheBatchSize, "Number of operations to batch in each new vertex.")
	StartnodeCmd.Flags().IntVar(&flags.SnowAvalancheNumParents, "snow-avalanche-num-parents", flags.SnowAvalancheNumParents, "Number of vertexes for reference from each new vertex.")
	StartnodeCmd.Flags().IntVar(&flags.SnowSampleSize, "snow-sample-size", flags.SnowSampleSize, "Number of nodes to query for each network poll.")
	StartnodeCmd.Flags().IntVar(&flags.SnowQuorumSize, "snow-quorum-size", flags.SnowQuorumSize, "Alpha value to use for required number positive results.")
	StartnodeCmd.Flags().IntVar(&flags.SnowVirtuousCommitThreshold, "snow-virtuous-commit-threshold", flags.SnowVirtuousCommitThreshold, "Beta value to use for virtuous transactions.")
	StartnodeCmd.Flags().IntVar(&flags.SnowRogueCommitThreshold, "snow-rogue-commit-threshold", flags.SnowRogueCommitThreshold, "Beta value to use for rogue transactions.")

	StartnodeCmd.Flags().BoolVar(&flags.P2PTLSEnabled, "p2p-tls-enabled", flags.P2PTLSEnabled, "Require TLS to authenticate network communications")
	StartnodeCmd.Flags().BoolVar(&flags.StakingTLSEnabled, "staking-tls-enabled", flags.StakingTLSEnabled, "Enable staking. If enabled, Network TLS is required.")
	StartnodeCmd.Flags().UintVar(&flags.StakingPort, "staking-port", flags.StakingPort, "Port of the consensus server.")
	StartnodeCmd.Flags().StringVar(&flags.StakingTLSCertFile, "staking-tls-cert-file", flags.StakingTLSCertFile, "TLS certificate file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.crt")
	StartnodeCmd.Flags().StringVar(&flags.StakingTLSKeyFile, "staking-tls-key-file", flags.StakingTLSKeyFile, "TLS private key file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.key")
}
