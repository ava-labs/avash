// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cmd

import (
	"strconv"
	"time"

	"github.com/ava-labs/avash/cfg"
	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// ProcmanagerCmd represents the procmanager command
var ProcmanagerCmd = &cobra.Command{
	Use:   "procmanager",
	Short: "Access the process manager for the avash client.",
	Long: `Access the process manager for the avash client. Using this 
	command you can list, stop, and start processes registered with the 
	process manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// PMListCmd represents the list operation on the procmanager command
var PMListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the processes currently running.",
	Long:  `Lists the processes currently running in tabular format.`,
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(AvalancheShell.rl.Stdout())
		table = pmgr.ProcManager.ProcessTable(table)
		table.Render()
	},
}

// PMMetadataCmd represents the list operation on the procmanager command
var PMMetadataCmd = &cobra.Command{
	Use:   "metadata [node name]",
	Short: "Prints the metadata associated with the node name.",
	Long:  `Prints the metadata associated with the node name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 && args[0] != "" {
			log := cfg.Config.Log
			name := args[0]
			metadata, err := pmgr.ProcManager.Metadata(name)
			if err != nil {
				log.Error(err.Error())
			}
			log.Info(metadata)
		} else {
			cmd.Help()
		}
	},
}

// PMStartCmd represents the start operation on the procmanager command
var PMStartCmd = &cobra.Command{
	Use:   "start [node name] [optional: delay in secs]",
	Short: "Starts the process named if not currently running.",
	Long:  `Starts the process named if not currently running.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 && args[0] != "" {
			log := cfg.Config.Log
			name := args[0]
			delay := time.Duration(0)
			if len(args) >= 2 {
				if v, e := strconv.ParseInt(args[1], 10, 64); e == nil && v > 0 {
					delay = time.Duration(v)
					log.Info("process will start in %ds: %s", int(delay), name)
				}
			}
			start := func() {
				err := pmgr.ProcManager.StartProcess(name)
				if err != nil {
					log.Error(err.Error())
				}
			}
			delayRun(start, delay)
		} else {
			cmd.Help()
		}
	},
}

// PMStopCmd represents the stop operation on the procmanager command
var PMStopCmd = &cobra.Command{
	Use:   "stop [node name] [optional: delay in secs]",
	Short: "Stops the process named if currently running.",
	Long:  `Stops the process named if currently running.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 && args[0] != "" {
			log := cfg.Config.Log
			name := args[0]
			delay := time.Duration(0)
			if len(args) >= 2 {
				if v, e := strconv.ParseInt(args[1], 10, 64); e == nil && v > 0 {
					delay = time.Duration(v)
					log.Info("process will stop in %ds: %s", int(delay), name)
				}
			}
			stop := func() {
				err := pmgr.ProcManager.StopProcess(name)
				if err != nil {
					log.Error(err.Error())
				}
			}
			delayRun(stop, delay)
		} else {
			cmd.Help()
		}
	},
}

// PMKillCmd represents the stop operation on the procmanager command
var PMKillCmd = &cobra.Command{
	Use:   "kill [node name] [optional: delay in secs]",
	Short: "Kills the process named if currently running.",
	Long:  `Kills the process named if currently running.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 && args[0] != "" {
			log := cfg.Config.Log
			name := args[0]
			delay := time.Duration(0)
			if len(args) >= 2 {
				if v, e := strconv.ParseInt(args[1], 10, 64); e == nil && v > 0 {
					delay = time.Duration(v)
					log.Info("process will stop in %ds: %s", int(delay), name)
				}
			}
			kill := func() {
				err := pmgr.ProcManager.KillProcess(name)
				if err != nil {
					log.Error(err.Error())
				}
			}
			delayRun(kill, delay)
		} else {
			cmd.Help()
		}
	},
}

// PMKillAllCmd stops all processes in the procmanager
var PMKillAllCmd = &cobra.Command{
	Use:   "killall [optional: delay in secs]",
	Short: "Kills all processes if currently running.",
	Long:  `Kills all processes if currently running.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		delay := time.Duration(0)
		if len(args) >= 1 {
			if v, e := strconv.ParseInt(args[0], 10, 64); e == nil && v > 0 {
				delay = time.Duration(v)
				log.Info("all processes will be killed in %ds", int(delay))
			}
		}
		delayRun(pmgr.ProcManager.KillAllProcesses, delay)
	},
}

// PMStopAllCmd stops all processes in the procmanager
var PMStopAllCmd = &cobra.Command{
	Use:   "stopall [optional: delay in secs]",
	Short: "Stops all processes if currently running.",
	Long:  `Stops all processes if currently running.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		delay := time.Duration(0)
		if len(args) >= 1 {
			if v, e := strconv.ParseInt(args[0], 10, 64); e == nil && v > 0 {
				delay = time.Duration(v)
				log.Info("all processes will stop in %ds", int(delay))
			}
		}
		delayRun(pmgr.ProcManager.StopAllProcesses, delay)
	},
}

// PMStartAllCmd starts all processes in the procmanager
var PMStartAllCmd = &cobra.Command{
	Use:   "startall [optional: delay in secs]",
	Short: "Starts all processes if currently stopped.",
	Long:  `Starts all processes if currently stopped.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		delay := time.Duration(0)
		if len(args) >= 1 {
			if v, e := strconv.ParseInt(args[0], 10, 64); e == nil && v > 0 {
				delay = time.Duration(v)
				log.Info("all processes will start in %ds", int(delay))
			}
		}
		delayRun(pmgr.ProcManager.StartAllProcesses, delay)
	},
}

// PMRemoveCmd represents the list operation on the procmanager command
var PMRemoveCmd = &cobra.Command{
	Use:   "remove [node name] [optional: delay in secs]",
	Short: "Removes the process named.",
	Long:  `Removes the process named. It will stop the process if it is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !(len(args) >= 1 && args[0] != "") {
			cmd.Help()
		}
		log := cfg.Config.Log
		name := args[0]
		delay := time.Duration(0)
		if len(args) >= 2 {
			if v, e := strconv.ParseInt(args[1], 10, 64); e == nil && v > 0 {
				delay = time.Duration(v)
				log.Info("process will be removed in %ds: %s", int(delay), name)
			}
		}
		remove := func() {
			err := pmgr.ProcManager.RemoveProcess(name)
			if err != nil {
				log.Error(err.Error())
			}
		}
		delayRun(remove, delay)
	},
}

// PMRemoveAllCmd represents the list operation on the procmanager command
var PMRemoveAllCmd = &cobra.Command{
	Use:   "removeall [optional: delay in secs]",
	Short: "Removes all processes.",
	Long:  `Removes all processes. It will stop the process if it is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		delay := time.Duration(0)
		if len(args) >= 1 {
			if v, e := strconv.ParseInt(args[0], 10, 64); e == nil && v > 0 {
				delay = time.Duration(v)
				log.Info("all processes will be removed in %ds", int(delay))
			}
		}
		delayRun(pmgr.ProcManager.RemoveAllProcesses, delay)
	},
}

func delayRun(f func(), delay time.Duration) {
	if delay == 0 {
		f()
		return
	}
	timer := time.NewTimer(delay * time.Second)
	go func() {
		<-timer.C
		f()
	}()
}

func init() {
	ProcmanagerCmd.AddCommand(PMKillCmd)
	ProcmanagerCmd.AddCommand(PMKillAllCmd)
	ProcmanagerCmd.AddCommand(PMListCmd)
	ProcmanagerCmd.AddCommand(PMMetadataCmd)
	ProcmanagerCmd.AddCommand(PMRemoveCmd)
	ProcmanagerCmd.AddCommand(PMRemoveAllCmd)
	ProcmanagerCmd.AddCommand(PMStopCmd)
	ProcmanagerCmd.AddCommand(PMStopAllCmd)
	ProcmanagerCmd.AddCommand(PMStartAllCmd)
	ProcmanagerCmd.AddCommand(PMStartCmd)
}
