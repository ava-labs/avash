package cmd

import (
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/vms/avm"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
	"github.com/spf13/cobra"

	pmgr "github.com/ava-labs/avash/processmgr"
)

// Client is a wrapper interface for avm.Client
type Client interface {
	IssueTx(txBytes []byte) (ids.ID, error)
	GetTxStatus(txID ids.ID) (choices.Status, error)
	GetBalance(addr string, assetID string, includePartial bool) (*avm.GetBalanceReply, error)
}

// Register adds the a command to parent.
// Used when root and commands are in different packages.
func Register(parent *cobra.Command) {
	parent.AddCommand(avaxWalletCmd)
}

var metadata = func(name string) (*node.Metadata, error) {
	return pmgr.ProcManager.NodeMetadata(name)
}

const (
	defaultEncoding         = formatting.CB58
	defaultAVMClientTimeout = time.Minute
)

var avmClient = func(host, port string, requestTimeout time.Duration) Client {
	return avm.NewClient(fmt.Sprintf("http://%s:%s", host, port), "avm", requestTimeout)
}

var avaxWalletCmd = &cobra.Command{
	Use:   "avaxwallet",
	Short: "Tools for interacting with AVAX Payments over the network.",
	Long: `Tools for interacting with AVAX Payments over the network. Using this 
	command you can send, and get the status of a transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// AVAXWalletNewKeyCmd creates a new private key
var avaxWalletNewKeyCmd = &cobra.Command{
	Use:   "newkey",
	Short: "Creates a random private key.",
	Long:  `Creates a random private key.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		if err := newKeyCmdRunE(cmd, args); err != nil {
			log.Error("%v", err)
		}
	},
}

// avaxWalletSendCmd will send a transaction through a node
var avaxWalletSendCmd = &cobra.Command{
	Use:   "send [node name] [tx string]",
	Short: "Sends a transaction to a node.",
	Long:  `Sends a transaction to a node.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		if err := sendCmdRunE(cmd, args); err != nil {
			log.Error("%v", err)
		}
	},
}

// avaxWalletStatusCmd will get the status of a transaction for a particular node
var avaxWalletStatusCmd = &cobra.Command{
	Use:   "status [node name] [tx id]",
	Short: "Checks the status of a transaction on a node.",
	Long:  `Checks the status of a transaction on a node.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		if err := statusCmdRunE(cmd, args); err != nil {
			log.Error("%v", err)
		}
	},
}

// AVAXWalletGetBalanceCmd will get the balance of an address from a node
var avaxWalletGetBalanceCmd = &cobra.Command{
	Use:   "balance [node name] [address]",
	Short: "Checks the balance of an address from a node.",
	Long:  `Checks the balance of an address from a node.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		if err := getBalanceCmdRunE(cmd, args); err != nil {
			log.Error("%v", err)
		}
	},
}

func newKeyCmdRunE(cmd *cobra.Command, args []string) error {
	log := cfg.Config.Log
	factory := crypto.FactorySECP256K1R{}
	skGen, err := factory.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("failed to create private key: %v", err)
	}

	sk, ok := skGen.(*crypto.PrivateKeySECP256K1R)
	if !ok {
		return fmt.Errorf("failed to create private key")
	}

	str, err := formatting.Encode(defaultEncoding, sk.Bytes())
	if err != nil {
		return fmt.Errorf("failed to encode private key")
	}

	log.Info("PrivateKey: PrivateKey-%s", str)
	return nil
}

func sendCmdRunE(cmd *cobra.Command, args []string) error {
	log := cfg.Config.Log
	md, err := metadata(args[0])
	if err != nil {
		return err
	}

	id, err := avmClient(md.Serverhost, md.HTTPport, defaultAVMClientTimeout).IssueTx([]byte(args[1]))
	if err != nil {
		return fmt.Errorf("failed to issue tx %s : %v", args[1], err)
	}

	log.Info("TxID: %s", id)
	return nil
}

func statusCmdRunE(cmd *cobra.Command, args []string) error {
	log := cfg.Config.Log
	md, err := metadata(args[0])
	if err != nil {
		return err
	}

	id, err := ids.FromString(args[1])
	if err != nil {
		return err
	}

	s, err := avmClient(md.Serverhost, md.HTTPport, defaultAVMClientTimeout).GetTxStatus(id)
	if err != nil {
		return fmt.Errorf("failed to get tx status %s : %v", args[1], err)
	}

	log.Info("Status: %s", s.String())
	return nil
}

func getBalanceCmdRunE(cmd *cobra.Command, args []string) error {
	log := cfg.Config.Log
	md, err := metadata(args[0])
	if err != nil {
		return err
	}

	b, err := avmClient(md.Serverhost, md.HTTPport, defaultAVMClientTimeout).GetBalance(args[1], "AVAX", false)
	switch {
	case err != nil:
		return fmt.Errorf("failed to get balance of %s : %v", args[1], err)
	case b == nil:
		return fmt.Errorf("get balance of %s is nil", args[1])
	}

	log.Info("Balance: %s", b.Balance)
	return nil
}

/*
avaxwallet
	create [wallet name] -> "wallet created: " + [wallet name]
	addkey [wallet name] [private key] -> address
	balance [node name] [address] -> uint
	status [node name] [tx string] -> [status]
	maketx [wallet name] [destination address] [amount] -> txString
	refresh [node name] [wallet name] -> "wallet refreshed: " + [wallet name]
	remove [wallet name] [tx string] -> "transaction removed: " + [tx string]
	send [node name] [tx string] -> "sent tx: " [tx string]
	newkey -> privateKey
*/

func init() {
	avaxWalletCmd.AddCommand(avaxWalletNewKeyCmd)
	avaxWalletCmd.AddCommand(avaxWalletGetBalanceCmd)
	avaxWalletCmd.AddCommand(avaxWalletSendCmd)
	avaxWalletCmd.AddCommand(avaxWalletStatusCmd)
}
