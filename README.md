# Avash &mdash; The Avalanche Shell Client

This is a temporary stateful shell execution environment used to deploy networks locally, manage their processes, and run network tests.

Avash opens a shell environment of its own. This environment is completely wiped when Avash exits. Any Avalanche nodes deployed by Avash should be exited as well, leaving only their stash (containing only their log files) behind.

Avash provides the ability to run Lua scripts which can execute a sequence of shell commands in Avash. This allows for automation of regular tasks. For instance, different network configurations can be programmed into a lua script and deployed as-needed, allowing for rapid tests against various network types.

## Installation

### Requirements

* Golang 1.15.5+
* An Avalanche Client Implementing Avalanche Standard CLI Flags

### Quick Setup

Install and build an Avalanche client

```zsh
go get github.com/ava-labs/avash
cd $GOPATH/src/github.com/ava-labs/avash
go build
```

 Now you can fire up a 5 node staking network:

```zsh
./avash
Config file set: /Users/username/.avash.yaml
Avash successfully configured.
avash> runscript scripts/five_node_staking.lua
RunScript: Running scripts/five_node_staking.lua
RunScript: Successfully ran scripts/five_node_staking.lua
```

For full documentation of Avash configuration and commands, please see the official [Avalanche Documentation](https://docs.avax.network/build/tools/avash).

## Using Avash

### Opening a shell

Super easy, just type `./avash` and it will open a shell environment.

#### Configuration

By default Avash will look for a configuration file named either `.avash.yaml` or `.avash.yml` located in the following paths.

* `$HOME/`
* `.`
* `/etc/avash/`

If no config file is found then Avash will create one at `$HOME/.avash.yaml`.

```zsh
./avash
Config file not found: .avash.yaml
Created empty config file: /Users/username/.avash.yaml
```

Alternatively you can pass in a `--config` flag with a path to your config file. **NOTE** you must put the full path. `~/` **will not** resolve to `$HOME/`.

```zsh
 ./avash --config=/Users/username/path/to/config/my-config-file.yaml
Config file set: /Users/username/path/to/config/my-config-file.yaml
Avash successfully configured.
```

If no config file is found at the path which was passed to `--config` then Avash will create one at `$HOME/`. Avash will use the filename which was passed to `--config`.

```zsh
./avash --config=/Users/username/path/to/config/my-config-file.yaml
Config file not found: /Users/username/path/to/config/my-config-file.yaml
Created empty config file: /Users/username/my-config-file.yaml
```

If you have multiple config files Avash will load the values from a single file in decreasing preference:

* `--config`
* `$HOME/`
* `.`
* `/etc/avash/`

#### Help

For your first command, type `help` in Avash to see the commands available.

You can also type `help [command]` to see the list of options available for that command.

Ex:

```zsh
help procmanager
help procmanager start
```

### Commands

* `avaxwallet` - Tools for interacting with Avalanche Payments over the network.
* `callrpc` - Issues an RPC call to a node.
* `exit` - Exit the shell.
* `help` - Help about any command.
* `network` - Tools for interacting with remote hosts.
* `procmanager` - Access the process manager for the avash client.
* `runscript` - Runs the provided script.
* `setoutput` - Sets shell log output.
* `startnode` - Starts a node process and gives it a name.
* `varstore` - Tools for creating variable stores and printing variables within them.

### Writing Scripts

Avash imports the [gopher-lua library](https://github.com/yuin/gopher-lua) to run lua scripts.

Scripts have certain hooks available to them which allows the user to write code which invokes the current Avash environment.

The functions available to Lua are:

* `avash_call` - Takes a string and runs it as an Avash command, returning output
* `avash_sleepmicro` - Takes an unsigned integer representing microseconds and sleeps for that long
* `avash_setvar` - Takes a variable scope (string), a variable name (string), and a variable (string) and places it in the variable store. The scope must already have been created.

 When writing Lua, the standard Lua functionality is available to automate the execution of series of Avash commands. This allows a developer to automate:

* Local network deployments
* Sending transations, both virtuous and conflicting
* Order transaction test cases
* Save the value of UTXO sets and test results to disk
* Compare the values of two nodes UTXO sets
* Track expected results and compare them with real nodes

Example Lua scripts are in [the `./scripts` directory](./scripts/).

### Funding a Wallet

On a local network, the 3 blockchains on the default subnet&mdash;the X-Chain, C-Chain and P-Chain&mdash;each have a pre-funded private key, `PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN`. This private key has 300m AVAX on the X-Chain, 50m AVAX on the C-Chain and 30m AVAX on the P-Chain&mdash;20m of which is unlocked and 10m which is locked and stakeable. For more details, see [Fund a local test network tutorial](https://docs.avax.network/build/tutorials/platform/fund-a-local-test-network).
