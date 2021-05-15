// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ava-labs/avash/cfg"
	"github.com/spf13/cobra"
	"github.com/yourbasic/radix"
)

// VarScope is a scope of the variable
type VarScope struct {
	Name      string
	Variables map[string]string
}

// List lists the variables in the scope
func (v *VarScope) List() []string {
	results := []string{}
	for k := range v.Variables {
		results = append(results, k)
	}
	return results
}

// Get gets the variable by name in the scope
func (v *VarScope) Get(varname string) (string, error) {
	if variable, ok := v.Variables[varname]; ok {
		return variable, nil
	}
	return "", fmt.Errorf("variable not found: %s", varname)
}

// Set sets the variable at a name to a value
func (v *VarScope) Set(varname string, value string) {
	v.Variables[varname] = value
}

// JSON returns the json representation of the variable scope
func (v *VarScope) JSON() ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

// VarStore stores scopes of variables to store
type VarStore struct {
	Stores map[string]VarScope
}

// Create will make a new variable scope
func (v *VarStore) Create(store string) error {
	if _, ok := v.Stores[store]; ok {
		return fmt.Errorf("store exists: %s", store)
	}
	v.Stores[store] = VarScope{
		Name:      store,
		Variables: map[string]string{},
	}
	return nil
}

// List lists the scopes available
func (v *VarStore) List() []string {
	results := []string{}
	for k := range v.Stores {
		results = append(results, k)
	}
	return results
}

// Get will retrieve the scope defined at the name passed in
func (v *VarStore) Get(store string) (VarScope, error) {
	if variable, ok := v.Stores[store]; ok {
		return variable, nil
	}
	return VarScope{}, fmt.Errorf("store not found: %s", store)
}

// VarStoreCmd represents the vars command
var VarStoreCmd = &cobra.Command{
	Use:   "varstore",
	Short: "Tools for creating variable stores and printing variables within them.",
	Long: `Tools for creating variable stores and printing variables within them. Using this 
	command you can create variable stores, list all variables they store, and print data 
	placed into these stores. Variable assigment and update is often managed by avash commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// VarStoreCreateCmd will attempt to get a genesis key and send a transaction
var VarStoreCreateCmd = &cobra.Command{
	Use:   "create [store name]",
	Short: "Creates a variable store.",
	Long: `Creates a variable store. If it exists, it prints "name conflict" otherwise 
	it prints "store created".`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			log := cfg.Config.Log
			store := args[0]
			if err := AvashVars.Create(store); err == nil {
				log.Info("store created: " + store)
			} else {
				log.Error("name conflict: " + store)
			}
		} else {
			cmd.Help()
		}
	},
}

// VarStoreListCmd will attempt to get a genesis key and send a transaction
var VarStoreListCmd = &cobra.Command{
	Use:   "list [store name]",
	Short: "Lists all stores. If store provided, lists all variables in the store.",
	Long: `Lists all stores. If store provided, lists all variables in the store. 
	If the store exists, it will print a new-line separated string of variables in 
	this store. If the store does not exist, it will print "store not found".`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		results := []string{}
		if len(args) >= 1 {
			if store, err := AvashVars.Get(args[0]); err == nil {
				results = store.List()
			} else {
				log.Error("store not found:" + args[0])
			}
		} else {
			results = AvashVars.List()
		}
		radix.Sort(results)
		for _, v := range results {
			log.Info(v)
		}
	},
}

// VarStorePrintCmd will attempt to get a genesis key and send a transaction
var VarStorePrintCmd = &cobra.Command{
	Use:   "print [store] [variable]",
	Short: "Prints a variable that is within the store.",
	Long:  `Prints a variable that is within the store. If it doesn't exist, it prints the default JSON string "{}".`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			log := cfg.Config.Log
			if store, err := AvashVars.Get(args[0]); err == nil {
				if v, e := store.Get(args[1]); e == nil {
					log.Info(v)
				} else {
					log.Info("{}")
				}
			} else {
				log.Info("{}")
			}
		} else {
			cmd.Help()
		}
	},
}

// VarStoreSetCmd will attempt to get a genesis key and send a transaction
var VarStoreSetCmd = &cobra.Command{
	Use:   "set [store] [variable] [value]",
	Short: "Sets a simple variable that within the store.",
	Long:  `Sets a simple variable that within the store. Store must exist. May not have spaces, even quoted. Existing values are overwritten.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 3 {
			log := cfg.Config.Log
			if store, err := AvashVars.Get(args[0]); err == nil {
				store.Set(args[1], args[2])
				log.Info("variable set: %q.%q=%q", args[0], args[1], args[2])
			} else {
				log.Error("store not found: " + args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// VarStoreStoreDumpCmd writes the store to the filename specified in the stash
var VarStoreStoreDumpCmd = &cobra.Command{
	Use:   "storedump [store] [filename]",
	Short: "Writes the store to a file.",
	Long:  `Writes the store to a file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 {
			log := cfg.Config.Log
			if store, err := AvashVars.Get(args[0]); err == nil {
				stashdir := cfg.Config.DataDir
				basename := filepath.Base(args[1])
				basedir := filepath.Dir(stashdir + "/" + args[1])

				os.MkdirAll(basedir, os.ModePerm)
				outputfile := basedir + "/" + basename

				if marshalled, err := store.JSON(); err == nil {
					if err := ioutil.WriteFile(outputfile, marshalled, 0755); err != nil {
						log.Error("unable to write file: %s - %s", string(outputfile), err.Error())
					} else {
						log.Info("VarStore written to: %s", outputfile)
					}
				} else {
					log.Error("unable to marshal: %s", err.Error())
				}
			} else {
				log.Error("store not found: %s", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// VarStoreVarDumpCmd writes the variable to the filename specified in the stash
var VarStoreVarDumpCmd = &cobra.Command{
	Use:   "vardump [store] [variable] [filename]",
	Short: "Writes the variable to a file.",
	Long:  `Writes the variable set to a file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 3 {
			log := cfg.Config.Log
			if store, err := AvashVars.Get(args[0]); err == nil {
				if variable, e := store.Get(args[1]); e == nil {
					stashdir := cfg.Config.DataDir
					basename := filepath.Base(args[2])
					basedir := filepath.Dir(stashdir + "/" + args[2])

					os.MkdirAll(basedir, os.ModePerm)
					outputfile := basedir + "/" + basename
					if err := ioutil.WriteFile(outputfile, []byte(variable), 0755); err != nil {
						log.Error("unable to write file: %s - %s", string(outputfile), err.Error())
					} else {
						log.Info("VarStore written to: %s", outputfile)
					}
				} else {
					log.Error("variable not found: %s -> %s", args[0], args[1])
				}
			} else {
				log.Error("store not found: %s", args[0])
			}
		} else {
			cmd.Help()
		}
	},
}

// AvashVars is the variable store.
var AvashVars VarStore

func init() {
	VarStoreCmd.AddCommand(VarStoreCreateCmd)
	VarStoreCmd.AddCommand(VarStoreStoreDumpCmd)
	VarStoreCmd.AddCommand(VarStoreListCmd)
	VarStoreCmd.AddCommand(VarStorePrintCmd)
	VarStoreCmd.AddCommand(VarStoreSetCmd)
	VarStoreCmd.AddCommand(VarStoreVarDumpCmd)
	AvashVars = VarStore{
		Stores: map[string]VarScope{},
	}
}
