/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"

	"github.com/ava-labs/avash/cfg"
	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/spf13/cobra"
)

// Metadata struct for storing metadata, available to commands
type Metadata struct {
	Serverhost     string `json:"public-ip"`
	Stakingport    string `json:"staking-port"`
	HTTPport       string `json:"http-port"`
	Dbdir          string `json:"db-dir"`
	Datadir        string `json:"data-dir"`
	Logsdir        string `json:"log-dir"`
	Loglevel       string `json:"log-level"`
	StakerCertPath string `json:"staking-tls-cert-file"`
	StakerKeyPath  string `json:"staking-tls-key-file"`
}

type nodeFlags struct {
    clientLocation                  string
    meta                            string
    dataDir                         string
	assertionsEnabled				bool
	avaTxFee						uint
	publicIP						string
	networkID						string
	xputServerPort					uint
	signatureVerificationEnabled	bool
	apiIpcsEnabled					bool
	httpPort						uint
	httpTLSEnabled					bool
	httpTLSCertFile					string
	httpTLSKeyFile					string
	bootstrapIps					string
	bootstrapIds					string
	dbEnabled						bool
	dbDir							string
	logLevel						string
	logDir							string
	snowAvalancheBatchSize			int
	snowAvalancheNumParents			int
	snowSampleSize					int
	snowQuorumSize					int
	snowVirtuousCommitThreshold		int
	snowRogueCommitThreshold		int
	stakingTLSEnabled				bool
	stakingPort						uint
	stakingTLSKeyFile				string
	stakingTLSCertFile				string
}

func defaultNodeFlags() nodeFlags {
	return nodeFlags{
        clientLocation:                 "",
        meta:                           "",
        dataDir:                        "",
		assertionsEnabled:				true,
		avaTxFee:						0,
		publicIP:						"127.0.0.1",
		networkID:						"local",
		xputServerPort:					9652,
		signatureVerificationEnabled:	true,
		apiIpcsEnabled:			        true,
		httpPort:                       9650,
		httpTLSEnabled:                 false,
		httpTLSCertFile:                "",
		httpTLSKeyFile:                 "",
		bootstrapIps:                   "",
		bootstrapIds:                   "",
		dbEnabled:                      true,
		dbDir:                          "db1",
		logLevel:                       "info",
		logDir:                         "logs",
		snowAvalancheBatchSize:         30,
		snowAvalancheNumParents:        5,
		snowSampleSize:                 2,
		snowQuorumSize:                 2,
		snowVirtuousCommitThreshold:    5,
		snowRogueCommitThreshold:       10,
		stakingTLSEnabled:              false,
		stakingPort:                    9651,
        stakingTLSKeyFile:              "",
		stakingTLSCertFile:             "",
	}
}

var flags nodeFlags

// StartnodeCmd represents the startnode command
var StartnodeCmd = &cobra.Command{
	Use:   "startnode [node name] args...",
	Short: "Starts a node process and gives it a name.",
	Long: `Starts an ava client node using pmgo and gives it a name. Example:
	
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
            flags.snowSampleSize,
            flags.snowQuorumSize,
            flags.snowVirtuousCommitThreshold,
            flags.snowRogueCommitThreshold,
        )
		if err != nil {
			log.Error(err.Error())
			return
		}

		args, md := flagsToArgs(sanitize.Path(datapath))
		mdbytes, _ := json.MarshalIndent(md, " ", "    ")
		metadata := string(mdbytes)
		meta := flags.meta
		if meta != "" {
			metadata = meta
		}
		avalocation := flags.clientLocation
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

func flagsToArgs(basedir string) ([]string, Metadata) {

    // Port targets
    httpPortString := strconv.FormatUint(uint64(flags.httpPort), 10)
    stakingPortString := strconv.FormatUint(uint64(flags.stakingPort), 10)

	// Paths/directories
	dbPath := basedir + "/" + flags.dbDir
	dataPath := basedir + "/" + flags.dataDir
	logPath := basedir + "/" + flags.logDir

	// Staking settings
	wd, _ := os.Getwd()

	// If the path given in the flag doesn't begin with "/", treat it as relative
    // to the directory of the avash binary
    httpCertFile := flags.httpTLSCertFile
	if httpCertFile != "" && string(httpCertFile[0]) != "/" {
		httpCertFile = fmt.Sprintf("%s/%s", wd, httpCertFile)
	}

    httpKeyFile := flags.httpTLSKeyFile
	if httpKeyFile != "" && string(httpKeyFile[0]) != "/" {
		httpKeyFile = fmt.Sprintf("%s/%s", wd, httpKeyFile)
	}

	stakerCertFile := flags.stakingTLSCertFile
	if stakerCertFile != "" && string(stakerCertFile[0]) != "/" {
		stakerCertFile = fmt.Sprintf("%s/%s", wd, stakerCertFile)
	}

	stakerKeyFile := flags.stakingTLSKeyFile
	if stakerKeyFile != "" && string(stakerKeyFile[0]) != "/" {
		stakerKeyFile = fmt.Sprintf("%s/%s", wd, stakerKeyFile)
	}

	args := []string{
		"--assertions-enabled=" + strconv.FormatBool(flags.assertionsEnabled),
		"--ava-tx-fee=" + strconv.FormatUint(uint64(flags.avaTxFee), 10),
		"--public-ip=" + flags.publicIP,
		"--network-id=" + flags.networkID,
		"--xput-server-port=" + strconv.FormatUint(uint64(flags.xputServerPort), 10),
		"--signature-verification-enabled=" + strconv.FormatBool(flags.signatureVerificationEnabled),
		"--api-ipcs-enabled=" + strconv.FormatBool(flags.apiIpcsEnabled),
		"--http-port=" + httpPortString,
		"--http-tls-enabled=" + strconv.FormatBool(flags.httpTLSEnabled),
		"--http-tls-cert-file=" + httpCertFile,
		"--http-tls-key-file=" + httpKeyFile,
		"--bootstrap-ips=" + flags.bootstrapIps,
		"--bootstrap-ids=" + flags.bootstrapIds,
		"--db-enabled=" + strconv.FormatBool(flags.dbEnabled),
		"--db-dir=" + dbPath,
		"--log-level=" + flags.logLevel,
		"--log-dir=" + logPath,
		"--snow-avalanche-batch-size=" + strconv.Itoa(flags.snowAvalancheBatchSize),
		"--snow-avalanche-num-parents=" + strconv.Itoa(flags.snowAvalancheNumParents),
		"--snow-sample-size=" + strconv.Itoa(flags.snowSampleSize),
		"--snow-quorum-size=" + strconv.Itoa(flags.snowQuorumSize),
		"--snow-virtuous-commit-threshold=" + strconv.Itoa(flags.snowVirtuousCommitThreshold),
		"--snow-rogue-commit-threshold=" + strconv.Itoa(flags.snowRogueCommitThreshold),
		"--staking-tls-enabled=" + strconv.FormatBool(flags.stakingTLSEnabled),
		"--staking-port=" + stakingPortString,
		"--staking-tls-key-file=" + stakerKeyFile,
		"--staking-tls-cert-file=" + stakerCertFile,
	}

	metadata := Metadata{
		Serverhost:     flags.publicIP,
		Stakingport:    stakingPortString,
		HTTPport:       httpPortString,
		Dbdir:          dbPath,
		Datadir:        dataPath,
		Logsdir:        logPath,
		Loglevel:       flags.logLevel,
		StakerCertPath: stakerCertFile,
		StakerKeyPath:  stakerKeyFile,
	}

    // Reset flags for next `startnode` call
    flags = defaultNodeFlags()
	return args, metadata
}

func init() {
    flags = defaultNodeFlags()
	StartnodeCmd.Flags().StringVar(&flags.clientLocation, "client-location", flags.clientLocation, "Path to AVA node client, defaulting to the config file's value.")
	StartnodeCmd.Flags().StringVar(&flags.meta, "meta", flags.meta, "Override default metadata for the node process.")
    StartnodeCmd.Flags().StringVar(&flags.dataDir, "data-dir", flags.dataDir, "Name of directory for the data stash.")

	StartnodeCmd.Flags().BoolVar(&flags.assertionsEnabled, "assertions-enabled", flags.assertionsEnabled, "Turn on assertion execution.")
	StartnodeCmd.Flags().UintVar(&flags.avaTxFee, "ava-tx-fee", flags.avaTxFee, "Ava transaction fee, in $nAva.")

	StartnodeCmd.Flags().BoolVar(&flags.apiIpcsEnabled, "api-ipcs-enabled", flags.apiIpcsEnabled, "Turn on IPC.")
	StartnodeCmd.Flags().StringVar(&flags.publicIP, "public-ip", flags.publicIP, "Public IP of this node.")
	StartnodeCmd.Flags().StringVar(&flags.networkID, "network-id", flags.networkID, "Network ID this node will connect to.")
	StartnodeCmd.Flags().UintVar(&flags.xputServerPort, "xput-server-port", flags.xputServerPort, "Port of the deprecated throughput test server.")
	StartnodeCmd.Flags().BoolVar(&flags.signatureVerificationEnabled, "signature-verification-enabled", flags.signatureVerificationEnabled, "Turn on signature verification.")

	StartnodeCmd.Flags().UintVar(&flags.httpPort, "http-port", flags.httpPort, "Port of the HTTP server.")
	StartnodeCmd.Flags().BoolVar(&flags.httpTLSEnabled, "http-tls-enabled", flags.httpTLSEnabled, "Upgrade the HTTP server to HTTPS.")
	StartnodeCmd.Flags().StringVar(&flags.httpTLSCertFile, "http-tls-cert-file", flags.httpTLSCertFile, "TLS certificate file for the HTTPS server.")
	StartnodeCmd.Flags().StringVar(&flags.httpTLSKeyFile, "http-tls-key-file", flags.httpTLSKeyFile, "TLS private key file for the HTTPS server.")

	StartnodeCmd.Flags().StringVar(&flags.bootstrapIps, "bootstrap-ips", flags.bootstrapIps, "Comma separated list of bootstrap nodes to connect to. Example: 127.0.0.1:9630,127.0.0.1:9620")
	StartnodeCmd.Flags().StringVar(&flags.bootstrapIds, "bootstrap-ids", flags.bootstrapIds, "Comma separated list of bootstrap peer ids to connect to. Example: JR4dVmy6ffUGAKCBDkyCbeZbyHQBeDsET,8CrVPQZ4VSqgL8zTdvL14G8HqAfrBr4z")

	StartnodeCmd.Flags().BoolVar(&flags.dbEnabled, "db-enabled", flags.dbEnabled, "Turn on persistent storage.")
	StartnodeCmd.Flags().StringVar(&flags.dbDir, "db-dir", flags.dbDir, "Database directory for Ava state.")

	StartnodeCmd.Flags().StringVar(&flags.logLevel, "log-level", flags.logLevel, "Specify the log level. Should be one of {verbo, debug, info, warn, error, fatal, off}")
	StartnodeCmd.Flags().StringVar(&flags.logDir, "log-dir", flags.logDir, "Name of directory for the node's logging.")

	StartnodeCmd.Flags().IntVar(&flags.snowAvalancheBatchSize, "snow-avalanche-batch-size", flags.snowAvalancheBatchSize, "Number of operations to batch in each new vertex.")
	StartnodeCmd.Flags().IntVar(&flags.snowAvalancheNumParents, "snow-avalanche-num-parents", flags.snowAvalancheNumParents, "Number of vertexes for reference from each new vertex.")
	StartnodeCmd.Flags().IntVar(&flags.snowSampleSize, "snow-sample-size", flags.snowSampleSize, "Number of nodes to query for each network poll.")
	StartnodeCmd.Flags().IntVar(&flags.snowQuorumSize, "snow-quorum-size", flags.snowQuorumSize, "Alpha value to use for required number positive results.")
	StartnodeCmd.Flags().IntVar(&flags.snowVirtuousCommitThreshold, "snow-virtuous-commit-threshold", flags.snowVirtuousCommitThreshold, "Beta value to use for virtuous transactions.")
	StartnodeCmd.Flags().IntVar(&flags.snowRogueCommitThreshold, "snow-rogue-commit-threshold", flags.snowRogueCommitThreshold, "Beta value to use for rogue transactions.")

	StartnodeCmd.Flags().BoolVar(&flags.stakingTLSEnabled, "staking-tls-enabled", flags.stakingTLSEnabled, "Require TLS to authenticate staking connections.")
	StartnodeCmd.Flags().UintVar(&flags.stakingPort, "staking-port", flags.stakingPort, "Port of the consensus server.")
	StartnodeCmd.Flags().StringVar(&flags.stakingTLSCertFile, "staking-tls-cert-file", flags.stakingTLSCertFile, "TLS certificate file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.crt")
	StartnodeCmd.Flags().StringVar(&flags.stakingTLSKeyFile, "staking-tls-key-file", flags.stakingTLSKeyFile, "TLS private key file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.key")
}
