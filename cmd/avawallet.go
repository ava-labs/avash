package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
	"github.com/spf13/cobra"
	"github.com/ybbus/jsonrpc"

	pmgr "github.com/ava-labs/avash/processmgr"
)

// Register adds the a command to parent.
// Used when root and commands are in different packages.
func Register(parent *cobra.Command) {
	parent.AddCommand(avaxWalletCmd)
}

var metadata = func(name string) (string, error) {
	return pmgr.ProcManager.Metadata(name)
}

const (
	defaultEncoding = formatting.CB58
)

var avmRPCClient = func(host string, port string) jsonrpc.RPCClient {
	url := fmt.Sprintf("http://%s:%s/ext/bc/avm", host, port)
	return jsonrpc.NewClient(url)
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
	meta, err := metadata(args[0])
	if err != nil {
		return fmt.Errorf("node not found: %s", args[0])
	}

	var md node.Metadata
	if err := json.Unmarshal([]byte(meta), &md); err != nil {
		return fmt.Errorf("failed to unmarshal metadata for node %s: %v", args[0], err)
	}

	rpcClient := avmRPCClient(md.Serverhost, md.HTTPport)
	response, err := rpcClient.Call("avm.issueTx", struct {
		Tx string
	}{
		Tx: args[1],
	})
	switch {
	case err != nil:
		return fmt.Errorf("failed to issue tx %s : %v", args[1], err)
	case response.Error != nil:
		return fmt.Errorf("failed to issue tx %s : %d, %s", args[1], response.Error.Code, response.Error.Message)
	}

	var s struct {
		TxID string
	}
	if err = response.GetObject(&s); err != nil {
		return fmt.Errorf("failed to unmarshal tx response: %v", err)
	}

	log.Info("TxID:%s", s.TxID)
	return nil
}

func statusCmdRunE(cmd *cobra.Command, args []string) error {
	log := cfg.Config.Log
	meta, err := metadata(args[0])
	if err != nil {
		return fmt.Errorf("node not found: %s", args[0])
	}

	var md node.Metadata
	if err := json.Unmarshal([]byte(meta), &md); err != nil {
		return fmt.Errorf("failed to unmarshal metadata for node %s: %v", args[0], err)
	}

	rpcClient := avmRPCClient(md.Serverhost, md.HTTPport)
	response, err := rpcClient.Call("avm.getTxStatus", struct {
		TxID string
	}{
		TxID: args[1],
	})

	switch {
	case err != nil:
		return fmt.Errorf("failed to issue tx %s : %v", args[1], err)
	case response.Error != nil:
		return fmt.Errorf("failed to issue tx %s : %d, %s", args[1], response.Error.Code, response.Error.Message)
	}

	var s struct {
		Status string
	}
	if err = response.GetObject(&s); err != nil {
		return fmt.Errorf("failed to unmarshal status response: %v", err)
	}

	log.Info("Status:%s", s.Status)
	return nil
}

func getBalanceCmdRunE(cmd *cobra.Command, args []string) error {
	log := cfg.Config.Log
	meta, err := metadata(args[0])
	if err != nil {
		return fmt.Errorf("node not found: %s", args[0])
	}

	var md node.Metadata
	if err := json.Unmarshal([]byte(meta), &md); err == nil {
		return fmt.Errorf("failed to unmarshal metadata for node %s: %v", args[0], err)
	}

	rpcClient := avmRPCClient(md.Serverhost, md.HTTPport)
	response, err := rpcClient.Call("avm.getBalance", struct {
		Address string
		AssetID string
	}{
		Address: args[1],
		AssetID: "AVAX",
	})

	switch {
	case err != nil:
		return fmt.Errorf("failed to issue tx %s : %v", args[1], err)
	case response.Error != nil:
		return fmt.Errorf("failed to issue tx %s : %d, %s", args[1], response.Error.Code, response.Error.Message)
	}

	var s struct {
		Balance string
	}
	if err = response.GetObject(&s); err != nil {
		return fmt.Errorf("failed to unmarshal status response: %v", err)
	}

	log.Info("Balance: %s", s.Balance)
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
