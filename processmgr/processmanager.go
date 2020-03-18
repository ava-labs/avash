/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

// Package processmgr manages processes launched by avash
package processmgr

import (
	"errors"
	"fmt"
	"os/exec"
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
func (p *ProcessManager) AddProcess(cmdstr string, proctype string, args []string, name string, metadata string, ih InputHandler, oh OutputHandler, eh OutputHandler) error {
	pname := strings.TrimSpace(name)
	if pname == "" {
		return errors.New("Process name cannot be empty")
	}
	_, exists := p.processes[pname]
	if exists {
		return fmt.Errorf("Process with name %s already exists", pname)
	}
	cout := make(chan []byte)
	cerr := make(chan []byte)
	cin := make(chan []byte)
	stop := make(chan bool)
	kill := make(chan bool)
	fail := make(chan error)
	proc := &Process{
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
	p.processes[name] = proc

	return nil
}

// StartProcess starts the process at the name
func (p *ProcessManager) StartProcess(name string) error {
	if _, ok := p.processes[name]; !ok {
		return fmt.Errorf("Can't start process %s because it does not exist", name)

	}
	if p.processes[name].running {
		return fmt.Errorf("Proccess %s is already running", name)
	}
	p.processes[name].cmd = exec.Command(p.processes[name].cmdstr, p.processes[name].args...)
	done := make(chan bool)
	go p.processes[name].Start(done)
	<-done
	return nil
}

// StopProcess stops the process at the name
func (p *ProcessManager) StopProcess(name string) error {
	if _, ok := p.processes[name]; !ok {
		return fmt.Errorf("Cannot stop process '%s' because it doesn't exist", name)

	}
	if !p.processes[name].running {
		return fmt.Errorf("Cannot stop process '%s' because it isn't running", name)
	}
	return p.processes[name].Stop()
}

// StopAllProcesses goes through each process and calls Stop(), returning name of process that failed and error
func (p *ProcessManager) StopAllProcesses() (string, error) {
	for name := range p.processes {
		if p.processes[name].running {
			err := p.StopProcess(name)
			if err != nil {
				cfg.Config.Log.Error("Error while stopping process '%s': %s", name, err)
				return name, err
			}
		}
	}
	return "", nil
}

// KillProcess kills the process at the name
func (p *ProcessManager) KillProcess(name string) error {
	if _, ok := p.processes[name]; ok {
		if !p.processes[name].running {
			return errors.New("process cannot kill: '" + name + "' isn't running.")
		}
	} else {
		return errors.New("process cannot kill: '" + name + "' doesn't exist.")
	}
	return p.processes[name].Kill()
}

// KillAllProcesses goes through each process and calls Kill(), returning name of process that failed and error
func (p *ProcessManager) KillAllProcesses() (string, error) {
	for name := range p.processes {
		if p.processes[name].running {
			err := p.KillProcess(name)
			if err != nil {
				return name, err
			}
		}
	}
	return "", nil
}

// StartAllProcesses goes through each process and calls Start(), returning name of process that failed and error
func (p *ProcessManager) StartAllProcesses() (string, error) {
	for name := range p.processes {
		if !p.processes[name].running {
			err := p.StartProcess(name)
			if err != nil {
				return name, err
			}
		}
	}
	return "", nil
}

// RemoveProcess removes a process from the list of available named processes
func (p *ProcessManager) RemoveProcess(name string) error {
	if _, ok := p.processes[name]; !ok {
		return fmt.Errorf("Cannot remove process '%s' because it doesn't exist", name)
	}
	p.StopProcess(name)
	delete(p.processes, name)
	cfg.Config.Log.Info("Process %s removed", name)
	return nil
}

// ProcessTable returns a formatted tablmetadatae for the data provided
func (p *ProcessManager) ProcessTable(table *tablewriter.Table) *tablewriter.Table {
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
	psd := p.ProcessSummary()
	table.AppendBulk(*psd)
	table.SetReflowDuringAutoWrap(true)
	return table
}

// ProcessSummary returns data table of all processes and their statuses
func (p *ProcessManager) ProcessSummary() *[][]string {
	var data [][]string
	for _, val := range p.processes {
		var running string
		if val.running {
			running = "running"
		} else if val.failed {
			running = "defunct"
		} else {
			running = "stopped"
		}
		line := []string{val.name, running, val.metadata, val.cmdstr}
		data = append(data, line)
	}
	return &data

}

// Metadata returns the metadata given the node name
func (p *ProcessManager) Metadata(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("Node name required")
	}
	if pr, ok := p.processes[name]; ok {
		return pr.metadata, nil
	}
	return "", fmt.Errorf("No such node: %s", name)
}

// ProcManager is a global
var ProcManager = ProcessManager{
	processes: make(map[string]*Process),
}
