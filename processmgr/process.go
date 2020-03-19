/*
Copyright © 2019 AVA Labs <collin@avalabs.org>
*/

package processmgr

import (
	"bytes"
	"errors"
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
	log.Info("\rStarting process %s.", p.name)
	log.Info("Command: %s\n", p.cmd.Args)
	if p.running {
		log.Error("Process %s is already running", p.name)
		return
	}
	p.running = true
	done <- true

	go func() {
		err := p.cmd.Run()
		// Procmanager not expecting termination
		if p.running {
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
			log.Info("\rCalling stop() on %s", p.name)
			if sp {
				if err := p.endProcess(false); err != nil {
					log.Error("SIGINT failed on process: %s", p.name)
					p.stop <- false
					return
				}
				log.Info("SIGINT called on process: %s", p.name)
				p.stop <- true
				return
			}
		case kl := <-p.kill:
			log.Info("\rCalling kill() on %s.", p.name)
			if kl {
				if err := p.endProcess(true); err != nil {
					log.Error("SIGTERM failed on process: %s", p.name)
					p.kill <- false
					return
				}
				log.Info("SIGTERM called on process: %s", p.name)
				p.kill <- true
				return
			}
		case fl := <-p.fail:
			errMsg := "inspect related logs for FATAL output"
			if fl != nil {
				errMsg = fl.Error()
			}
			log.Error("\rProcess %s failure: %s", p.name, errMsg)
			p.running = false
			p.failed = true
			return
		}
	}
}

// Stop ends a process with SIGINT
func (p *Process) Stop() error {
	if !p.running {
		return fmt.Errorf("Cannot stop process %s because it is not running", p.name)
	}
	p.stop <- true
	result := <-p.stop
	if result {
		p.running = false
	} else {
		return fmt.Errorf("Unable to stop process %s: ", p.name)
	}
	return nil
}

// Kill ends a process with SIGTERM
func (p *Process) Kill() error {
	if p.running {
		p.kill <- true
		result := <-p.kill
		if result {
			p.running = false
		} else {
			return errors.New("Unable to kill process: " + p.name)
		}
	} else {
		return errors.New("Process is not running, cannot kill: " + p.name)
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
