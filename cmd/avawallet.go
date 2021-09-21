// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avash/cfg"
	"github.com/ava-labs/avash/node"
	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/spf13/cobra"

	"github.com/ybbus/jsonrpc"
)

// AVAXWalletCmd represents the avaxwallet command
var AVAXWalletCmd = &cobra.Command{
	Use:   "avaxwallet",
	Short: "Tools for interacting with AVAX Payments over the network.",
	Long: `Tools for interacting with AVAX Payments over the network. Using this 
	command you can send, and get the status of a transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

const (
	defaultEncoding = formatting.CB58
)

// AVAXWalletNewKeyCmd creates a new private key
var AVAXWalletNewKeyCmd = &cobra.Command{
	Use:   "newkey",
	Short: "Creates a random private key.",
	Long:  `Creates a random private key.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		factory := crypto.FactorySECP256K1R{}
		if skGen, err := factory.NewPrivateKey(); err == nil {
			sk := skGen.(*crypto.PrivateKeySECP256K1R)
			str, err := formatting.EncodeWithChecksum(defaultEncoding, sk.Bytes())
			if err != nil {
				log.Error("could not encode private key")
			}
			log.Info("PrivateKey: PrivateKey-%s", str)
		} else {
			log.Error("could not create private key")
		}
	},
}

// AVAXWalletSendCmd will send a transaction through a node
var AVAXWalletSendCmd = &cobra.Command{
	Use:   "send [node name] [tx string]",
	Short: "Sends a transaction to a node.",
	Long:  `Sends a transaction to a node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			log := cfg.Config.Log
			if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
				var md node.Metadata
				metaBytes := []byte(meta)
				if err := json.Unmarshal(metaBytes, &md); err == nil {
					jrpcloc := fmt.Sprintf("http://%s:%s/ext/bc/avm", md.Serverhost, md.HTTPport)
					rpcClient := jsonrpc.NewClient(jrpcloc)
					response, err := rpcClient.Call("avm.issueTx", struct {
						Tx string
					}{
						Tx: args[1],
					})
					if err != nil {
						log.Error("error sent tx: %s", args[1])
						log.Error("rpcClient returned error: %s", err.Error())
					} else if response.Error != nil {
						log.Error("error sent tx: %s", args[1])
						log.Error("rpcClient returned error: %d, %s", response.Error.Code, response.Error.Message)
					} else {
						var s struct {
							TxID string
						}
						err = response.GetObject(&s)
						if err != nil {
							log.Error("error on parsing response: %s", err.Error())
						} else {
							log.Info("TxID:%s", s.TxID)
						}
					}
				} else {
					log.Error("unable to unmarshal metadata for node %s: %s", args[0], err.Error())
				}
			} else {
				log.Error("node not found: %s", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAXWalletStatusCmd will get the status of a transaction for a particular node
var AVAXWalletStatusCmd = &cobra.Command{
	Use:   "status [node name] [tx id]",
	Short: "Checks the status of a transaction on a node.",
	Long:  `Checks the status of a transaction on a node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			log := cfg.Config.Log
			if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
				var md node.Metadata
				metaBytes := []byte(meta)
				if err := json.Unmarshal(metaBytes, &md); err == nil {
					jrpcloc := fmt.Sprintf("http://%s:%s/ext/bc/avm", md.Serverhost, md.HTTPport)
					rpcClient := jsonrpc.NewClient(jrpcloc)
					response, err := rpcClient.Call("avm.getTxStatus", struct {
						TxID string
					}{
						TxID: args[1],
					})
					if err != nil {
						log.Error("error sent txid: %s", args[1])
						log.Error("rpcClient returned error: %s", err.Error())
					} else if response.Error != nil {
						log.Error("error sent txid: %s", args[1])
						log.Error("rpcClient returned error: %d, %s", response.Error.Code, response.Error.Message)
					} else {
						var s struct {
							Status string
						}
						err = response.GetObject(&s)
						if err != nil {
							log.Error("error on parsing response: %s", err.Error())
						} else {
							log.Info("Status:%s", s.Status)
						}
					}
				} else {
					log.Error("unable to unmarshal metadata for node %s: %s", args[0], err.Error())
				}
			} else {
				log.Error("node not found: %s", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAXWalletGetBalanceCmd will get the balance of an address from a node
var AVAXWalletGetBalanceCmd = &cobra.Command{
	Use:   "balance [node name] [address]",
	Short: "Checks the balance of an address from a node.",
	Long:  `Checks the balance of an address from a node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			log := cfg.Config.Log
			if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
				var md node.Metadata
				metaBytes := []byte(meta)
				if err := json.Unmarshal(metaBytes, &md); err == nil {
					jrpcloc := fmt.Sprintf("http://%s:%s/ext/bc/avm", md.Serverhost, md.HTTPport)
					rpcClient := jsonrpc.NewClient(jrpcloc)
					response, err := rpcClient.Call("avm.getBalance", struct {
						Address string
						AssetID string
					}{
						Address: args[1],
						AssetID: "AVAX",
					})
					if err != nil {
						log.Error("error sent address: %s", args[1])
						log.Error("rpcClient returned error: %s", err.Error())
					} else if response.Error != nil {
						log.Error("error sent address: %s", args[1])
						log.Error("rpcClient returned error: %d, %s", response.Error.Code, response.Error.Message)
					} else {
						var s struct {
							Balance string
						}
						err = response.GetObject(&s)
						if err != nil {
							log.Error("error on parsing response: %s", err.Error())
						} else {
							log.Info("Balance: %s", s.Balance)
						}
					}
				} else {
					log.Error("unable to unmarshal metadata for node %s: %s", args[0], err.Error())
				}
			} else {
				log.Error("node not found: %s", args[0])
			}
		} else {
			cmd.Help()
		}
	},
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
	AVAXWalletCmd.AddCommand(AVAXWalletNewKeyCmd)
	AVAXWalletCmd.AddCommand(AVAXWalletGetBalanceCmd)
	AVAXWalletCmd.AddCommand(AVAXWalletSendCmd)
	AVAXWalletCmd.AddCommand(AVAXWalletStatusCmd)
}
