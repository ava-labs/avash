# Avash &mdash; The AVA Shell Client

This is a temporary stateful shell execution environment used to deploy networks locally, manage their processes, and run network tests.

Avash opens a shell environment of its own. This environment is completely wiped when Avash exits. Any AVA nodes deployed by Avash should be exited as well, leaving only their stash behind.

Avash provides the ability to run Lua scripts which can execute a sequence of shell commands in Avash. This allows for automation of regular tasks. For instance, different network configurations can be programmed into a lua script and deployed as-needed, allowing for rapid tests against various network types.

## Installation

### Requirements

 * Golang 1.13+
 * An AVA Client Implementing AVA Standard CLI Flags

### Quick Setup

 1. Install and build an AVA client
 2. `go get github.com/ava-labs/avash`
 3. `cd $GOPATH/src/github.com/ava-labs/avash`
 4. `go build`

For full documentation of Avash configuration and commands, please see the official [AVA Documentation](https://docs.ava.network/v1.0/en/tools/avash/).

## Using Avash

### Opening a shell

Super easy, just type `./avash` and it will open a shell environement.

For your first command, type `help` in Avash to see the commands available. 

You can also type `help [command]` to see the list of options available for that command.

Ex:

```sh
help procmanager
help procmanager start
```

### Commands

 * avawallet - Tools for interacting with AVA Payments over the network.
 * exit - Exit the shell.
 * help - Help about any command.
 * procmanager - Access the process manager for the avash client.
 * runscript - Runs the provided script.
 * setoutput - Sets shell log output.
 * startnode - Starts a node process and gives it a name.
 * varstore - Tools for creating variable stores and printing variables within them.

### Writing Scripts

Avash imports the gopher-lua library (https://github.com/yuin/gopher-lua) to run lua scripts.

Scripts have certain hooks available to them which allows the user to write code which invokes the current Avash environment.

The functions available to Lua are:

 * avash_call - Takes a string and runs it as an Avash command, returning output
 * avash_sleepmicro - Takes an unsigned integer representing microseconds and sleeps for that long
 * avash_setvar - Takes a variable scope (string), a variable name (string), and a variable (string) and places it in the variable store. The scope must already have been created.

 When writing Lua, the standard Lua functionality is available to automate the execution of series of Avash commands. This allows a developer to automate:

 * Local network deployments 
 * Sending transations, both virtuous and conflicting
 * Order transaction test cases
 * Save the value of UTXO sets and test results to disk
 * Compare the values of two nodes UTXO sets
 * Track expected results and compare them with real nodes
 
 Example Lua scripts are in the `./scripts` folder. 


