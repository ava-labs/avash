// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

// Package processmgr manages processes launched by avash
package processmgr

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avash/cfg"
	"github.com/olekukonko/tablewriter"
)

// ProcessManager is a system for managing processes in the system
type ProcessManager struct {
	// Key: Process name
	// Value: The corresponding process
	processes map[string]*Process
}

// AddProcess places a process into the process manager with an associated name
func (pm *ProcessManager) AddProcess(cmdstr string, proctype string, args []string, name string, metadata string, ih InputHandler, oh OutputHandler, eh OutputHandler) error {
	pname := strings.TrimSpace(name)
	if pname == "" {
		return fmt.Errorf("Process name cannot be empty")
	}
	_, exists := pm.processes[pname]
	if exists {
		return fmt.Errorf("Process with name %s already exists", pname)
	}
	cout := make(chan []byte)
	cerr := make(chan []byte)
	cin := make(chan []byte)
	stop := make(chan bool)
	kill := make(chan bool)
	fail := make(chan error)
	p := &Process{
		cmdstr:    cmdstr,
		args:      args,
		name:      pname,
		proctype:  proctype,
		metadata:  metadata,
		cout:      cout,
		cerr:      cerr,
		cin:       cin,
		stop:      stop,
		kill:      kill,
		fail:      fail,
		inhandle:  ih,
		outhandle: oh,
		errhandle: eh,
	}
	pm.processes[name] = p

	return nil
}

// StartProcess starts the process at the name
func (pm *ProcessManager) StartProcess(name string) error {
	p, ok := pm.processes[name]
	if !ok {
		return fmt.Errorf("Process does not exist, cannot start: %s", name)
	}
	done := make(chan bool)
	go p.Start(done)
	<-done
	return nil
}

// StopProcess stops the process at the name
func (pm *ProcessManager) StopProcess(name string) error {
	p, ok := pm.processes[name]
	if !ok {
		return fmt.Errorf("Process does not exist, cannot stop: %s", name)
	}
	return p.Stop()
}

// StopAllProcesses calls Stop() on every running process, logging errors
func (pm *ProcessManager) StopAllProcesses() {
	existsRunning := false
	for name := range pm.processes {
		if pm.processes[name].running {
			existsRunning = true
			err := pm.StopProcess(name)
			if err != nil {
				cfg.Config.Log.Error(err.Error())
			}
		}
	}
	if !existsRunning {
		cfg.Config.Log.Info("No processes currently running.")
	}
}

// KillProcess kills the process at the name
func (pm *ProcessManager) KillProcess(name string) error {
	p, ok := pm.processes[name]
	if !ok {
		return fmt.Errorf("Process does not exist, cannot kill: %s", name)
	}
	return p.Kill()
}

// KillAllProcesses calls Kill() on every running process, logging errors
func (pm *ProcessManager) KillAllProcesses() {
	existsRunning := false
	for name := range pm.processes {
		if pm.processes[name].running {
			existsRunning = true
			err := pm.KillProcess(name)
			if err != nil {
				cfg.Config.Log.Error(err.Error())
			}
		}
	}
	if !existsRunning {
		cfg.Config.Log.Info("No processes currently running.")
	}
}

// StartAllProcesses calls Start() on every stopped process, logging errors
func (pm *ProcessManager) StartAllProcesses() {
	existsStopped := false
	for name := range pm.processes {
		if !pm.processes[name].running {
			existsStopped = true
			err := pm.StartProcess(name)
			if err != nil {
				cfg.Config.Log.Error(err.Error())
			}
		}
	}
	if !existsStopped {
		cfg.Config.Log.Info("All processes currently running.")
	}
}

// RemoveProcess removes a process from the list of available named processes
func (pm *ProcessManager) RemoveProcess(name string) error {
	if _, ok := pm.processes[name]; !ok {
		return fmt.Errorf("Process does not exist, cannot remove: %s", name)
	}
	if pm.processes[name].running {
		if err := pm.StopProcess(name); err != nil {
			return err
		}
	}
	delete(pm.processes, name)
	cfg.Config.Log.Info("Process removed: %s", name)
	return nil
}

// RemoveAllProcesses removes a process from the list of available named processes
func (pm *ProcessManager) RemoveAllProcesses() {
	pm.StopAllProcesses()
	totalProcesses := len(pm.processes)
	processesRemoved := 0
	for name := range pm.processes {
		if err := pm.RemoveProcess(name); err != nil {
			cfg.Config.Log.Error(err.Error())
			continue
		}
		processesRemoved++
	}
	cfg.Config.Log.Info("%d/%d processes removed", processesRemoved, totalProcesses)
}

// ProcessTable returns a formatted metadata table for the data provided
func (pm *ProcessManager) ProcessTable(table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{"Name", "Status", "Metadata", "Command"})
	table.SetBorder(false)

	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgMagentaColor, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgMagentaColor, tablewriter.FgWhiteColor})

	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Normal},
		tablewriter.Colors{tablewriter.Normal},
		tablewriter.Colors{tablewriter.Normal})
	psd := pm.ProcessSummary()
	table.AppendBulk(*psd)
	table.SetReflowDuringAutoWrap(true)
	return table
}

// ProcessSummary returns data table of all processes and their statuses
func (pm *ProcessManager) ProcessSummary() *[][]string {
	var data [][]string
	for _, val := range pm.processes {
		var running string
		if val.running {
			running = "running"
		} else if val.failed {
			running = "defunct"
		} else {
			running = "stopped"
		}
		cmd := val.cmdstr + " " + strings.Join(val.args, " ")
		line := []string{val.name, running, val.metadata, cmd}
		data = append(data, line)
	}
	return &data

}

// Metadata returns the metadata given the process name
func (pm *ProcessManager) Metadata(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("Process name required")
	}
	if p, ok := pm.processes[name]; ok {
		return p.metadata, nil
	}
	return "", fmt.Errorf("Process does not exist, cannot get metadata: %s", name)
}

// HasRunning returns true if there exists a running process, otherwise false
func (pm *ProcessManager) HasRunning() bool {
	for _, val := range pm.processes {
		if val.running {
			return true
		}
	}
	return false
}

// ProcManager is a global
var ProcManager = ProcessManager{
	processes: make(map[string]*Process),
}
