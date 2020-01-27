/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	pmgr "github.com/ava-labs/avash/processmgr"
	dagwallet "github.com/ava-labs/avash/wallets/dags"
	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/formatting"
	"github.com/ava-labs/gecko/vms/spdagvm"
	"github.com/spf13/cobra"

	"github.com/ava-labs/gecko/utils/crypto"

	"github.com/ybbus/jsonrpc"
)

// AVAWalletCmd represents the avawallet command
var AVAWalletCmd = &cobra.Command{
	Use:   "avawallet [operation]",
	Short: "Tools for interacting with AVA Payments over the network.",
	Long: `Tools for interacting with AVA Payments over the network. Using this 
	command you can create, send, and get the status of a transaction.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("avawallet requires an operation. Available: create, addkey, maketx, refresh, remove, send, createkey")
	},
}

// AVAWalletCreateCmd creates a new named wallet
var AVAWalletCreateCmd = &cobra.Command{
	Use:   "create [wallet name]",
	Short: "Creates a wallet.",
	Long:  `Creates a wallet persistent for this session.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			dagwallet.Wallets[args[0]] = dagwallet.NewWallet()
			fmt.Printf("wallet created: %s\n", args[0])
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletNewKeyCmd creates a new private key
var AVAWalletNewKeyCmd = &cobra.Command{
	Use:   "newkey",
	Short: "Creates a random private key.",
	Long:  `Creates a random private key.`,
	Run: func(cmd *cobra.Command, args []string) {
		factory := crypto.FactorySECP256K1R{}
		if skGen, err := factory.NewPrivateKey(); err == nil {
			sk := skGen.(*crypto.PrivateKeySECP256K1R)
			fb := formatting.FormatBytes{}
			fb.Bytes = sk.Bytes()
			fmt.Printf("Pk:%s\n", fb.String())
		} else {
			fmt.Printf("could not create private key\n")
		}
	},
}

// AVAWalletAddKeyCmd adds a private key to a wallet
var AVAWalletAddKeyCmd = &cobra.Command{
	Use:   "addkey [wallet name] [private key]",
	Short: "Adds a private key to a wallet.",
	Long:  `Adds a private key to a wallet from a b58 string and returns its address. Reminder: refresh the UTXOs after keys are imported.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if w, ok := dagwallet.Wallets[args[0]]; ok {
				factory := crypto.FactorySECP256K1R{}
				fb := formatting.FormatBytes{}
				fb.FromString(args[1])
				if skGen, err := factory.ToPrivateKey(fb.Bytes); err == nil {
					sk := skGen.(*crypto.PrivateKeySECP256K1R)
					w.ImportKey(sk)
					fmt.Printf("Addr:%s\n", skGen.PublicKey().Address().String())
				} else {
					fmt.Printf("unable to add key %s: %s\n", args[1], err.Error())
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletMakeTxCmd will create a transaction and return its signed string
var AVAWalletMakeTxCmd = &cobra.Command{
	Use:   "maketx [wallet name] [destination address] [amount]",
	Short: "Creates a signed transaction.",
	Long:  `Creates a signed transaction for an amount to an address. Returns the a string of the transaction.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 3 {
			if w, ok := dagwallet.Wallets[args[0]]; ok {
				if amount, err := strconv.ParseUint(args[2], 10, 64); err == nil {
					fb := formatting.FormatBytes{}
					fb.FromString(args[1])
					toAddr, err := ids.ToShortID(fb.Bytes)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					signedTx := w.CreateTx(amount, 0, 1, []ids.ShortID{toAddr})
					if signedTx != nil {
						if err := signedTx.Verify(0); err == nil {
							fb.Bytes = signedTx.Bytes()
							fmt.Printf("Tx:%s\n", fb.String())
						} else {
							fmt.Println("signedTx cannot verify")
						}
					} else {
						fmt.Println("unable to create tx, check UTXO set")
					}
				} else {
					fmt.Printf("amount %s cannot convert to uint64\n", args[2])
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletRemoveCmd will remove a transaction from the UTXO set
var AVAWalletRemoveCmd = &cobra.Command{
	Use:   "remove [wallet name] [tx string]",
	Short: "Removes a transaction from a wallet's UTXO set.",
	Long:  `Removes a transaction from a wallet's UTXO set.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if w, ok := dagwallet.Wallets[args[0]]; ok {
				fb := formatting.FormatBytes{}
				fb.FromString(args[1])
				txBytes := fb.Bytes
				codec := spdagvm.Codec{}
				tx, err := codec.UnmarshalTx(txBytes)
				if err == nil {
					for _, in := range tx.Ins() {
						utxoID := in.InputID()
						w.RemoveUtxo(utxoID)
					}
					fmt.Printf("transaction removed: %s\n", args[1])
				} else {
					fmt.Printf("cannot unmarshal tx: %s\n", args[1])
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletSpendCmd will spend (update inputs and outputs) a transaction from the UTXO set
var AVAWalletSpendCmd = &cobra.Command{
	Use:   "spend [wallet name] [tx string]",
	Short: "Spends a transaction from a wallet's UTXO set.",
	Long:  `Spends a transaction from a wallet's UTXO set.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if w, ok := dagwallet.Wallets[args[0]]; ok {
				fb := formatting.FormatBytes{}
				fb.FromString(args[1])
				txBytes := fb.Bytes
				codec := spdagvm.Codec{}
				tx, err := codec.UnmarshalTx(txBytes)
				if err == nil {
					w.SpendTx(tx)
					fmt.Printf("transaction spent: %s\n", args[1])
				} else {
					fmt.Printf("cannot unmarshal tx: %s\n", args[1])
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletSendCmd will send a transaction through a node
var AVAWalletSendCmd = &cobra.Command{
	Use:   "send [node name] [tx string]",
	Short: "Sends a transaction to a node.",
	Long:  `Sends a transaction to a node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
				var md Metadata
				metaBytes := []byte(meta)
				if err := json.Unmarshal(metaBytes, &md); err == nil {
					jrpcloc := fmt.Sprintf("http://%s:%s/ext/ava", md.Jrpchost, md.Jrpcport)
					rpcClient := jsonrpc.NewClient(jrpcloc)
					response, err := rpcClient.Call("Ava.IssueTx", struct {
						Tx string
					}{
						Tx: args[1],
					})
					if err != nil {
						fmt.Printf("error sent tx: %s\n", args[1])
						fmt.Printf("rpcClient returned error: %s\n", err.Error())
					} else if response.Error != nil {
						fmt.Printf("error sent tx: %s\n", args[1])
						fmt.Printf("rpcClient returned error: %d, %s\n", response.Error.Code, response.Error.Message)
					} else {
						var s struct {
							TxID string
						}
						err = response.GetObject(&s)
						if err != nil {
							fmt.Printf("error on parsing response: %s\n", err.Error())
						} else {
							fmt.Printf("TxID:%s\n", s.TxID)
						}
					}
				} else {
					fmt.Printf("unable to unmarshal metadata for node %s: %s\n", args[0], err.Error())
				}
			} else {
				fmt.Printf("node not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletStatusCmd will get the status of a transaction for a particular node
var AVAWalletStatusCmd = &cobra.Command{
	Use:   "status [node name] [tx id]",
	Short: "Checks the status of a transaction on a node.",
	Long:  `Checks the status of a transaction on a node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
				var md Metadata
				metaBytes := []byte(meta)
				if err := json.Unmarshal(metaBytes, &md); err == nil {
					jrpcloc := fmt.Sprintf("http://%s:%s/ext/ava", md.Jrpchost, md.Jrpcport)
					rpcClient := jsonrpc.NewClient(jrpcloc)
					response, err := rpcClient.Call("Ava.GetTxStatus", struct {
						TxID string
					}{
						TxID: args[1],
					})
					if err != nil {
						fmt.Printf("error sent txid: %s\n", args[1])
						fmt.Printf("rpcClient returned error: %s\n", err.Error())
					} else if response.Error != nil {
						fmt.Printf("error sent txid: %s\n", args[1])
						fmt.Printf("rpcClient returned error: %d, %s\n", response.Error.Code, response.Error.Message)
					} else {
						var s struct {
							Status string
						}
						err = response.GetObject(&s)
						if err != nil {
							fmt.Printf("error on parsing response: %s\n", err.Error())
						} else {
							fmt.Printf("Status:%s\n", s.Status)
						}
					}
				} else {
					fmt.Printf("unable to unmarshal metadata for node %s: %s\n", args[0], err.Error())
				}
			} else {
				fmt.Printf("node not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletGetBalanceCmd will get the balance of an address from a node
var AVAWalletGetBalanceCmd = &cobra.Command{
	Use:   "balance [node name] [address]",
	Short: "Checks the balance of an address from a node.",
	Long:  `Checks the balance of an address from a node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
				var md Metadata
				metaBytes := []byte(meta)
				if err := json.Unmarshal(metaBytes, &md); err == nil {
					jrpcloc := fmt.Sprintf("http://%s:%s/ext/wallet", md.Jrpchost, md.Jrpcport)
					rpcClient := jsonrpc.NewClient(jrpcloc)
					response, err := rpcClient.Call("Wallet.GetBalance", struct {
						SubnetAlias string
						Address     string
						AssetID     string
					}{
						SubnetAlias: "ava",
						Address:     args[1],
						AssetID:     "ava",
					})
					if err != nil {
						fmt.Printf("error sent address: %s\n", args[1])
						fmt.Printf("rpcClient returned error: %s\n", err.Error())
					} else if response.Error != nil {
						fmt.Printf("error sent address: %s\n", args[1])
						fmt.Printf("rpcClient returned error: %d, %s\n", response.Error.Code, response.Error.Message)
					} else {
						var s struct {
							Balance uint64
						}
						err = response.GetObject(&s)
						if err != nil {
							fmt.Printf("error on parsing response: %s\n", err.Error())
						} else {
							fmt.Printf("Balance:%d\n", s.Balance)
						}
					}
				} else {
					fmt.Printf("unable to unmarshal metadata for node %s: %s\n", args[0], err.Error())
				}
			} else {
				fmt.Printf("node not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletRefreshCmd will send a transaction through a node
var AVAWalletRefreshCmd = &cobra.Command{
	Use:   "refresh [node name] [wallet name]",
	Short: "Refreshes UTXO set from node.",
	Long:  `Refreshes UTXO set from node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if w, ok := dagwallet.Wallets[args[1]]; ok {
				if meta, err := pmgr.ProcManager.Metadata(args[0]); err == nil {
					var md Metadata
					metaBytes := []byte(meta)
					if err := json.Unmarshal(metaBytes, &md); err == nil {
						jrpcloc := fmt.Sprintf("http://%s:%s/ext/ava", md.Jrpchost, md.Jrpcport)
						rpcClient := jsonrpc.NewClient(jrpcloc)

						response, err := rpcClient.Call("Ava.GetUTXOs", struct {
							Addresses []string
						}{
							Addresses: w.Addresses(),
						})
						if err != nil {
							fmt.Printf("rpcClient returned error: %s\n", err.Error())
						} else if response.Error != nil {
							fmt.Printf("rpcClient returned error: %d, %s\n", response.Error.Code, response.Error.Message)
						} else {
							var s struct {
								UTXOs []string
							}
							err = response.GetObject(&s)
							if err != nil {
								fmt.Printf("error on parsing response: %s\n", err.Error())
							} else {
								fb := formatting.FormatBytes{}
								acodec := spdagvm.Codec{}
								for _, aUTXO := range s.UTXOs {
									fb.FromString(aUTXO)
									if utxo, err := acodec.UnmarshalUTXO(fb.Bytes); err == nil {
										w.AddUtxo(utxo)
									} else {
										fmt.Printf("unable to add UTXO: %s\n", aUTXO)
									}
								}
								//fmt.Printf("[%s]\n", strings.Join(s.UTXOs, ","))
								fmt.Printf("utxo set refreshed on wallet %s from node %s\n", args[1], args[0])
							}
						}
					} else {
						fmt.Printf("unable to unmarshal metadata for node %s: %s\n", args[0], err.Error())
					}
				} else {
					fmt.Printf("node not found: %s\n", args[0])
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[1])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletWriteUTXOCmd writes the UTXOs of a wallet to the filename specified in the stash
var AVAWalletWriteUTXOCmd = &cobra.Command{
	Use:   "writeutxo [wallet name A] [filename]",
	Short: "Writes the UTXO set to a file.",
	Long:  `Writes the UTXO set to a file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			if wallet, ok := dagwallet.Wallets[args[0]]; ok {
				pwd, _ := os.Getwd()
				stashdir := pwd + "/stash"
				basename := filepath.Base(args[1])
				basedir := filepath.Dir(stashdir + "/" + args[1])

				os.MkdirAll(basedir, os.ModePerm)
				outputfile := basedir + "/" + basename
				utxoset := wallet.GetUtxos()

				if marshalled, err := utxoset.JSON(); err == nil {
					if err := ioutil.WriteFile(outputfile, marshalled, 0755); err != nil {
						fmt.Printf("unable to write file: %s - %s\n", string(outputfile), err.Error())
					} else {
						fmt.Printf("UTXO Set written to: %s\n", outputfile)
					}
				} else {
					fmt.Printf("unable to marshal: %s\n", err.Error())
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AVAWalletCompareCmd compares the UTXO set between two wallets, stores difference in a variable
var AVAWalletCompareCmd = &cobra.Command{
	Use:   "compare [wallet name A] [wallet name B] [variable scope] [variable name]",
	Short: "Compares the UTXO set between two wallets.",
	Long:  `Compares the UTXO set between two wallets.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 4 {
			if w1, ok := dagwallet.Wallets[args[0]]; ok {
				if w2, ok := dagwallet.Wallets[args[1]]; ok {
					if store, err := AvashVars.Get(args[2]); err == nil {
						us1 := w1.GetUtxos()
						us2 := w2.GetUtxos()
						diff := us1.SetDiff(us2)
						diffByte, err := json.MarshalIndent(diff, "", "    ")
						if err != nil {
							fmt.Println("unable to marshal: ", err.Error())
						} else {
							store.Set(args[3], string(diffByte))
						}
					} else {
						fmt.Println("store not found: " + args[2])
					}
				} else {
					fmt.Printf("wallet not found: %s\n", args[1])
				}
			} else {
				fmt.Printf("wallet not found: %s\n", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

/*
avawallet
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
	AVAWalletCmd.AddCommand(AVAWalletCreateCmd)
	AVAWalletCmd.AddCommand(AVAWalletNewKeyCmd)
	AVAWalletCmd.AddCommand(AVAWalletAddKeyCmd)
	AVAWalletCmd.AddCommand(AVAWalletGetBalanceCmd)
	AVAWalletCmd.AddCommand(AVAWalletMakeTxCmd)
	AVAWalletCmd.AddCommand(AVAWalletRemoveCmd)
	AVAWalletCmd.AddCommand(AVAWalletSpendCmd)
	AVAWalletCmd.AddCommand(AVAWalletSendCmd)
	AVAWalletCmd.AddCommand(AVAWalletRefreshCmd)
	AVAWalletCmd.AddCommand(AVAWalletCompareCmd)
	AVAWalletCmd.AddCommand(AVAWalletStatusCmd)
	AVAWalletCmd.AddCommand(AVAWalletWriteUTXOCmd)
}
