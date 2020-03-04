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
	flag "github.com/spf13/pflag"

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

		name := args[0]

		datadir := cfg.Config.DataDir
		basename := sanitize.BaseName(name)
		datapath := datadir + "/" + basename
		if basename == "" {
			fmt.Println("Process name can't be empty")
			return
		}
		f := cmd.Flags()

		k, _ := f.GetInt("snow-sample-size")
		alpha, _ := f.GetInt("snow-quorum-size")
		beta1, _ := f.GetInt("snow-virtuous-commit-threshold")
		beta2, _ := f.GetInt("snow-rogue-commit-threshold")

		err := validateConsensusArgs(k, alpha, beta1, beta2)
		if err != nil {
			fmt.Println(err)
			return
		}

		args, md := flagsToArgs(f, sanitize.Path(datapath))
		mdbytes, _ := json.MarshalIndent(md, " ", "    ")
		metadata := string(mdbytes)
		meta, _ := cmd.Flags().GetString("meta")
		if meta != "" {
			metadata = meta
		}
		avalocation, _ := f.GetString("client-location")
		if avalocation == "" {
			avalocation = cfg.Config.AvaLocation
		}
		err = pmgr.ProcManager.AddProcess(avalocation, "ava node", args, name, metadata, nil, nil, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		pmgr.ProcManager.StartProcess(name)
		fmt.Printf("Created process %s\n", name)
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

func flagsToArgs(f *flag.FlagSet, basedir string) ([]string, Metadata) {

	// Assertions
	assertions, _ := f.GetBool("assertions-enabled")
	useassertions := "false"
	if assertions {
		useassertions = "true"
	}

	// Transaction fees
	txfee, _ := f.GetUint("ava-tx-fee")

	// Network ID
	networkid, _ := f.GetString("network-id")

	// Host/port settings
	sh, _ := f.GetString("public-ip")
	sp, _ := f.GetUint("staking-port")
	rp, _ := f.GetUint("xput-server-port")
	hp, _ := f.GetUint("http-port")
	bootstrapips, _ := f.GetString("bootstrap-ips")
	bootstrapids, _ := f.GetString("bootstrap-ids")

	// Paths/directories
	dbdir, _ := f.GetString("db-dir")
	datadir, _ := f.GetString("data-dir")
	logdir, _ := f.GetString("log-dir")

	// Staking settings
	wd, _ := os.Getwd()

	// Assertions
	henabled, _ := f.GetBool("http-tls-enabled")
	httptlsenabled := "false"
	if henabled {
		httptlsenabled = "true"
	}

	hcert, _ := f.GetString("http-tls-cert-file")
	hkey, _ := f.GetString("http-tls-key-file")

	httptlscertparam := ""
	httptlskeyparam := ""
	if henabled {
		httptlscertparam = "--http-tls-cert-file=" + hcert
		httptlskeyparam = "--http-tls-key-file=" + hkey
	}

	// Signature verification
	sigver, _ := f.GetBool("signature-verification-enabled")
	sigverenabled := "false"
	if sigver {
		sigverenabled = "true"
	}

	// IPC
	ipc, _ := f.GetBool("ipc-enabled")
	ipcEnabled := "true"
	if ipc == false {
		ipcEnabled = "false"
	}

	stakingenabled, _ := f.GetBool("staking-tls-enabled")

	// If the path given in the flag doesn't begin with "/", treat it as relative
	// to the directory of the avash binary
	stakerCertFile, _ := f.GetString("staking-tls-cert-file")
	if stakerCertFile != "" && string(stakerCertFile[0]) != "/" {
		stakerCertFile = fmt.Sprintf("%s/%s", wd, stakerCertFile)
	}

	stakerKeyFile, _ := f.GetString("staking-tls-key-file")
	if stakerKeyFile != "" && string(stakerKeyFile[0]) != "/" {
		stakerKeyFile = fmt.Sprintf("%s/%s", wd, stakerKeyFile)
	}

	stakercertparam := ""
	stakerkeyparam := ""
	if stakingenabled {
		stakercertparam = stakerCertFile
		stakerkeyparam = stakerKeyFile
	}

	requirestaking := "false"
	if stakingenabled {
		requirestaking = "true"
	}

	// Log settings
	logLevel, _ := f.GetString("log-level")

	// Db settings
	hasdb, _ := f.GetBool("db-enabled")
	usedb := "false"
	if hasdb {
		usedb = "true"
	}

	// Consensus parameters
	k, _ := f.GetInt("snow-sample-size")
	alpha, _ := f.GetInt("snow-quorum-size")
	beta1, _ := f.GetInt("snow-virtuous-commit-threshold")
	beta2, _ := f.GetInt("snow-rogue-commit-threshold")
	batch, _ := f.GetInt("snow-avalanche-batch-size")
	numparents, _ := f.GetInt("snow-avalanche-num-parents")

	args := []string{
		"--assertions-enabled=" + useassertions,
		"--ava-tx-fee=" + strconv.FormatUint(uint64(txfee), 10),
		"--public-ip=" + sh,
		"--network-id=" + networkid,
		"--xput-server-port=" + strconv.FormatUint(uint64(rp), 10),
		"--signature-verification-enabled=" + sigverenabled,
		"--ipc-enabled=" + ipcEnabled,
		"--http-port=" + strconv.FormatUint(uint64(hp), 10),
		"--http-tls-enabled=" + httptlsenabled,
		"--http-tls-cert-file=" + httptlscertparam,
		"--http-tls-key-file=" + httptlskeyparam,
		"--bootstrap-ips=" + bootstrapips,
		"--bootstrap-ids=" + bootstrapids,
		"--db-enabled=" + usedb,
		"--db-dir=" + basedir + "/" + dbdir,
		"--log-level=" + logLevel,
		"--log-dir=" + basedir + "/" + logdir,
		"--snow-avalanche-batch-size=" + strconv.Itoa(batch),
		"--snow-avalanche-num-parents=" + strconv.Itoa(numparents),
		"--snow-sample-size=" + strconv.Itoa(k),
		"--snow-quorum-size=" + strconv.Itoa(alpha),
		"--snow-virtuous-commit-threshold=" + strconv.Itoa(beta1),
		"--snow-rogue-commit-threshold=" + strconv.Itoa(beta2),
		"--staking-tls-enabled=" + requirestaking,
		"--staking-port=" + strconv.FormatUint(uint64(sp), 10),
		"--staking-tls-key-file=" + stakerkeyparam,
		"--staking-tls-cert-file=" + stakercertparam,
	}

	metadata := Metadata{
		Serverhost:     sh,
		Stakingport:    strconv.FormatUint(uint64(sp), 10),
		HTTPport:       strconv.FormatUint(uint64(hp), 10),
		Dbdir:          basedir + "/" + dbdir,
		Datadir:        basedir + "/" + datadir,
		Logsdir:        basedir + "/" + logdir,
		Loglevel:       logLevel,
		StakerCertPath: stakerCertFile,
		StakerKeyPath:  stakerKeyFile,
	}

	return args, metadata
}

func init() {
	StartnodeCmd.Flags().String("client-location", "", "Path to AVA node client, defaulting to the config file's value.")
	StartnodeCmd.Flags().String("meta", "", "Override default metadata for the node process.")
	StartnodeCmd.Flags().String("data-dir", "", "Name of directory for the data stash.")

	StartnodeCmd.Flags().Bool("assertions-enabled", true, "Turn on assertion execution.")
	StartnodeCmd.Flags().Uint("ava-tx-fee", 0, "Ava transaction fee, in $nAva.")

	StartnodeCmd.Flags().String("public-ip", "127.0.0.1", "Public IP of this node.")
	StartnodeCmd.Flags().String("network-id", "12345", "Network ID this node will connect to.")
	StartnodeCmd.Flags().Uint("xput-server-port", 9652, "Port of the deprecated throughput test server.")
	StartnodeCmd.Flags().Bool("signature-verification-enabled", true, "Turn on signature verification.")
	StartnodeCmd.Flags().Bool("ipc-enabled", true, "Turn on IPC.")

	StartnodeCmd.Flags().Uint("http-port", 9650, "Port of the HTTP server.")
	StartnodeCmd.Flags().Bool("http-tls-enabled", false, "Upgrade the HTTP server to HTTPS.")
	StartnodeCmd.Flags().String("http-tls-cert-file", "", "TLS certificate file for the HTTPS server.")
	StartnodeCmd.Flags().String("http-tls-key-file", "", "TLS private key file for the HTTPS server.")

	StartnodeCmd.Flags().String("bootstrap-ips", "", "Comma separated list of bootstrap nodes to connect to. Example: 127.0.0.1:9630,127.0.0.1:9620")
	StartnodeCmd.Flags().String("bootstrap-ids", "", "Comma separated list of bootstrap peer ids to connect to. Example: JR4dVmy6ffUGAKCBDkyCbeZbyHQBeDsET,8CrVPQZ4VSqgL8zTdvL14G8HqAfrBr4z")

	StartnodeCmd.Flags().Bool("db-enabled", true, "Turn on persistent storage.")
	StartnodeCmd.Flags().String("db-dir", "db1", "Database directory for Ava state.")

	StartnodeCmd.Flags().String("log-level", "all", "Specify the log level. Should be one of {all, debug, info, warn, error, fatal, off}")
	StartnodeCmd.Flags().String("log-dir", "logs", "Name of directory for the node's logging.")

	StartnodeCmd.Flags().Int("snow-avalanche-batch-size", 30, "Number of operations to batch in each new vertex.")
	StartnodeCmd.Flags().Int("snow-avalanche-num-parents", 5, "Number of vertexes for reference from each new vertex.")
	StartnodeCmd.Flags().Int("snow-sample-size", 2, "Number of nodes to query for each network poll.")
	StartnodeCmd.Flags().Int("snow-quorum-size", 2, "Alpha value to use for required number positive results.")
	StartnodeCmd.Flags().Int("snow-virtuous-commit-threshold", 5, "Beta value to use for virtuous transactions.")
	StartnodeCmd.Flags().Int("snow-rogue-commit-threshold", 10, "Beta value to use for rogue transactions.")

	StartnodeCmd.Flags().Bool("staking-tls-enabled", false, "Require TLS to authenticate staking connections.")
	StartnodeCmd.Flags().Uint("staking-port", 9651, "Port of the consensus server.")
	StartnodeCmd.Flags().String("staking-tls-cert-file", "", "TLS certificate file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.crt")
	StartnodeCmd.Flags().String("staking-tls-key-file", "", "TLS private key file for staking connections. Relative to the avash binary if doesn't start with '/'. Ex: certs/keys1/staker.key")
}
