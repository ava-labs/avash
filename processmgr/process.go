// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package processmgr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/ava-labs/avash/cfg"
)

// InputHandler is a generic function for handling input from cin
type InputHandler func(p []byte) (bool, error)

// OutputHandler recieves the information
type OutputHandler func(b bytes.Buffer) error

// Process declares the necessary data for tracking a process
type Process struct {
	cmdstr    string
	args      []string
	cmd       *exec.Cmd
	name      string
	proctype  string
	metadata  string
	running   bool
	failed    bool
	output    io.ReadCloser
	errput    io.ReadCloser
	input     io.WriteCloser
	cout      chan []byte
	cerr      chan []byte
	cin       chan []byte
	stop      chan bool
	kill      chan bool
	fail      chan error
	inhandle  InputHandler
	outhandle OutputHandler
	errhandle OutputHandler
}

// Start begins a new process
func (p *Process) Start(done chan bool) {
	log := cfg.Config.Log
	if p.running {
		log.Error("Process is already running, cannot start: %s", p.name)
		done <- true
		return
	}
	log.Info("Starting process %s.", p.name)
	p.cmd = exec.Command(p.cmdstr, p.args...)
	log.Info("Command: %s\n", p.cmd.Args)

	selfStopped := false
	go func() {
		err := p.cmd.Start()
		if err != nil {
			p.fail <- err
			return
		}
		p.running = true
		p.failed = false
		done <- true
		err = p.cmd.Wait()
		if !selfStopped {
			p.fail <- err
		}
	}()

	closegen := func() {
		p.cmd.Stdin = nil
		p.cmd.Stderr = nil
		p.cmd.Stdout = nil
		p.cmd.Wait()
		p.cmd.Process = nil
	}

	defer closegen()

	for {
		select {
		case sp := <-p.stop:
			log.Info("Calling stop() on %s", p.name)
			if sp {
				selfStopped = true
				if err := p.endProcess(false); err != nil {
					log.Error("SIGINT failed on process: %s: %s", p.name, err.Error())
					p.stop <- false
					return
				}
				log.Info("SIGINT called on process: %s", p.name)
				p.stop <- true
				return
			}
		case kl := <-p.kill:
			log.Info("Calling kill() on %s.", p.name)
			if kl {
				selfStopped = true
				if err := p.endProcess(true); err != nil {
					log.Error("SIGTERM failed on process: %s: %s", p.name, err.Error())
					p.kill <- false
					return
				}
				log.Info("SIGTERM called on process: %s", p.name)
				p.kill <- true
				return
			}
		case fl := <-p.fail:
			p.failed = true
			errMsg := "inspect for process validity (command, args, flags) or FATAL output in related logs"
			if fl != nil {
				errMsg = fl.Error()
			}
			log.Error("Process failure: %s: %s", p.name, errMsg)
			// Specific case for a bad `p.cmd.Start()` call
			if !p.running {
				done <- false
			}
			p.running = false
			return
		}
	}
}

// Stop ends a process with SIGINT
func (p *Process) Stop() error {
	if !p.running {
		return fmt.Errorf("Process is not running, cannot stop: %s", p.name)
	}
	p.stop <- true
	result := <-p.stop
	p.running = false
	if !result {
		p.failed = true
		return fmt.Errorf("Unable to properly stop process: %s", p.name)
	}
	return nil
}

// Kill ends a process with SIGTERM
func (p *Process) Kill() error {
	if !p.running {
		return fmt.Errorf("Process is not running, cannot kill: %s", p.name)
	}
	p.kill <- true
	result := <-p.kill
	p.running = false
	if !result {
		p.failed = true
		return fmt.Errorf("Unable to properly kill process: %s", p.name)
	}
	return nil
}

func (p *Process) endProcess(killer bool) error {
	if killer {
		if err := p.cmd.Process.Kill(); err != nil {
			return err
		}
	} else {
		if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
			if err := p.cmd.Process.Kill(); err != nil {
				return err
			}
		}
	}
	return nil
}
