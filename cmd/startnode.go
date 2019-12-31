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
	Serverhost     string `json:"serverhost"`
	Serverport     string `json:"serverport"`
	Jrpchost       string `json:"jrpchost"`
	Jrpcport       string `json:"jrpcport"`
	Dbdir          string `json:"dbdir"`
	Genesisdir     string `json:"genesisdir"`
	Logsdir        string `json:"logsdir"`
	Loglevel       string `json:"loglevel"`
	StakerCertPath string `json:"stakerCertPath"`
	StakerKeyPath  string `json:"stakerKeyPath"`
}

// StartnodeCmd represents the startnode command
var StartnodeCmd = &cobra.Command{
	Use:   "startnode [node name] args...",
	Short: "Starts a node process and gives it a name.",
	Long: `Starts an ava client node using pmgo and gives it a name. Example:
	
startnode 127.0.0.1 9001 localhost 9002 9008`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		name := args[0]

		datadir := cfg.Viper.GetString("datadir")
		if datadir == "" {
			wd, _ := os.Getwd()
			datadir = wd + "/stash"
		}
		basename := sanitize.BaseName(name)
		datapath := datadir + "/" + basename
		if basename == "" {
			fmt.Println("Process name can't be empty")
			return
		}
		f := cmd.Flags()

		k, _ := f.GetInt("k")
		alpha, _ := f.GetInt("alpha")
		beta1, _ := f.GetInt("beta1")
		beta2, _ := f.GetInt("beta2")

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
		avalocation, _ := f.GetString("avaloc")
		if avalocation == "" {
			avalocation = cfg.Viper.GetString("avalocation")
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

	// Host/port settings
	sh, _ := f.GetString("serverhost")
	sp, _ := f.GetString("serverport")
	rh, _ := f.GetString("rpchost")
	rp, _ := f.GetString("rpcport")
	jh, _ := f.GetString("jrpchost")
	jp, _ := f.GetString("jrpcport")
	bootstrapips, _ := f.GetString("bootstrapips")

	// Paths/directories
	dbdir, _ := f.GetString("dbdir")
	gendir, _ := f.GetString("genesisdir")
	logdir, _ := f.GetString("logsdir")

	// Staking settings
	wd, _ := os.Getwd()
	stakingenabled, _ := f.GetBool("staking_enabled")
	stakerCertFile, _ := f.GetString("stake_cert_file")
	// If the path given in the flag doesn't begin with "/", treat it as relative
	// to the directory of the avash binary
	if stakerCertFile != "" && string(stakerCertFile[0]) != "/" {
		stakerCertFile = fmt.Sprintf("%s/%s", wd, stakerCertFile)
	}
	stakerKeyFile, _ := f.GetString("stake_key_file")
	if stakerKeyFile != "" && string(stakerKeyFile[0]) != "/" {
		stakerKeyFile = fmt.Sprintf("%s/%s", wd, stakerKeyFile)
	}
	requirestaking := "--require_staking=true"
	if !stakingenabled {
		requirestaking = "--require_staking=false"
	}

	// Log settings
	logLevel, _ := f.GetString("loglevel")

	// Db settings
	hasdb, _ := f.GetBool("db")
	usedb := "--db=false"
	if hasdb {
		usedb = "--db=true"
	}

	// Consensus parameters
	k, _ := f.GetInt("k")
	alpha, _ := f.GetInt("alpha")
	beta1, _ := f.GetInt("beta1")
	beta2, _ := f.GetInt("beta2")

	args := []string{
		"--server_ip=" + sh + ":" + sp,
		"--rpc_ip=" + rh + ":" + rp,
		"--jrpc_ip=" + jh,
		"--jrpc_port=" + jp,
		"--bootstrap_ips=" + bootstrapips,
		"--db_dir=" + basedir + "/" + dbdir,
		"--genesis_dir=" + basedir + "/" + gendir,
		"--log_level=" + logLevel,
		"--logs_dir=" + basedir + "/" + logdir,
		requirestaking,
		usedb,
		"--k=" + strconv.Itoa(k),
		"--alpha=" + strconv.Itoa(alpha),
		"--beta1=" + strconv.Itoa(beta1),
		"--beta2=" + strconv.Itoa(beta2),
		"--stake_cert_file=" + stakerCertFile,
		"--stake_key_file=" + stakerKeyFile,
	}

	metadata := Metadata{
		Serverhost:     sh,
		Serverport:     sp,
		Jrpchost:       jh,
		Jrpcport:       jp,
		Dbdir:          basedir + "/" + dbdir,
		Genesisdir:     basedir + "/" + gendir,
		Logsdir:        basedir + "/" + logdir,
		Loglevel:       logLevel,
		StakerCertPath: wd + stakerCertFile,
		StakerKeyPath:  wd + stakerKeyFile,
	}

	return args, metadata
}

func init() {
	StartnodeCmd.Flags().String("avaloc", "", "Path to AVA node binary.")
	StartnodeCmd.Flags().String("meta", "", "Override default metadata for the node process.")

	StartnodeCmd.Flags().String("serverhost", "127.0.0.1", "Server host for the node.")
	StartnodeCmd.Flags().String("serverport", "9651", "Server port for the node.")
	StartnodeCmd.Flags().String("rpchost", "127.0.0.1", "RPC host for the node.")
	StartnodeCmd.Flags().String("rpcport", "9652", "RPC port for the node.")
	StartnodeCmd.Flags().String("jrpchost", "127.0.0.1", "JSON RPC host for the node.")
	StartnodeCmd.Flags().String("jrpcport", "9650", "JSON RPC port for the node.")

	StartnodeCmd.Flags().String("bootstrapips", "", "Comma separated list of bootstrap nodes to connect to. Example: 127.0.0.1:9630,127.0.0.1:9620")

	StartnodeCmd.Flags().String("dbdir", "db1", "Name of database folder for the node.")
	StartnodeCmd.Flags().String("genesisdir", "data", "Name of directory for the genesis key.")
	StartnodeCmd.Flags().String("loglevel", "info", "Specify the log level. Should be one of {all, debug, info, warn, error, fatal, off}")
	StartnodeCmd.Flags().String("logsdir", "logs", "Name of directory for the node's logging.")

	StartnodeCmd.Flags().Bool("staking_enabled", false, "Turns on staking.")
	StartnodeCmd.Flags().Bool("db", true, "Sets if data should be persistently stored.")

	StartnodeCmd.Flags().Int("k", 2, "Number of nodes to query for each network poll.")
	StartnodeCmd.Flags().Int("alpha", 2, "Alpha value to use for required number positive results.")
	StartnodeCmd.Flags().Int("beta1", 5, "Beta value to use for virtuous transactions.")
	StartnodeCmd.Flags().Int("beta2", 10, "Beta value to use for rogue transactions.")

	StartnodeCmd.Flags().String("stake_cert_file", "", "The path of the staker certificate file. Relative to the avash binary iff doesn't start with /. Ex: certs/keys1/staker.crt")
	StartnodeCmd.Flags().String("stake_key_file", "", "The path of the staker certificate key. Relative to the avash binary iff doesn't start with /. Ex: certs/keys1/staker.key")
}
