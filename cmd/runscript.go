// Copyright © 2021 AVA Labs, Inc.
// All rights reserved.

package cmd

import (
	//"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	//"context"

	"go.uber.org/multierr"

	"github.com/ava-labs/avash/cfg"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
)

// RunScriptCmd represents the exit command
var RunScriptCmd = &cobra.Command{
	Use:   "runscript [script file]",
	Short: "Runs the provided script.",
	Long:  `Runs the script provided in the argument, relative to the present working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			log := cfg.Config.Log
			L := lua.NewState( /*lua.Options{
				RegistrySize:        1024 * 20,   // this is the initial size of the registry
				RegistryMaxSize:     1024 * 8000, // this is the maximum size that the registry can grow to. If set to `0` (the default) then the registry will not auto grow
				RegistryGrowStep:    32,          // this is how much to step up the registry by each time it runs out of space. The default is `32`.
				CallStackSize:       240,         // this is the maximum callstack size of this LState
				MinimizeStackMemory: true,        // Defaults to `false` if not specified. If set, the callstack will auto grow and shrink as needed up to a max of `CallStackSize`. If not set, the callstack will be fixed at `CallStackSize`.

			}*/)
			L.OpenLibs()
			defer L.Close()
			//ctx, cancel := context.WithCancel(context.Background())
			//L.SetContext(ctx)
			//defer cancel()

			/* set new Lua functions here */
			L.SetGlobal("avash_call", L.NewFunction(AvashCall))
			L.SetGlobal("avash_sleepmicro", L.NewFunction(AvashSleepMicro))
			L.SetGlobal("avash_setvar", L.NewFunction(AvashSetVar))
			//L.SetGlobal("avash_coroutine", L.NewFunction(AvashCoroutine))

			filename := args[0]
			log.Info("RunScript: Running " + filename)

			if err := L.DoFile(filename); err != nil {
				log.Error("RunScript: Failed to run " + filename + "\n" + err.Error())
			} else {
				log.Info("RunScript: Successfully ran " + filename)
			}
		} else {
			cmd.Help()
		}
	},
}

func capture() func() (string, error) {
	re, we, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	ro, wo, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	saveO := os.Stdout
	os.Stdout = wo

	saveE := os.Stderr
	os.Stderr = we

	var buf strings.Builder

	go func() {
		_, StderrIoCopyError := io.Copy(&buf, re)
		StderrReadCloseError := re.Close()
		_, StdoutIoCopyError := io.Copy(&buf, ro)
		StdoutReadCloseError := ro.Close()
		AllErrors := multierr.Combine(
			StderrIoCopyError,
			StderrReadCloseError,
			StdoutIoCopyError,
			StdoutReadCloseError,
		)
		done <- AllErrors
	}()

	return func() (string, error) {
		os.Stderr = saveE
		we.Close()
		os.Stdout = saveO
		wo.Close()
		err := <-done
		return buf.String(), err
	}
}

// AvashSleepMicro function to sleep for N microseconds
func AvashSleepMicro(L *lua.LState) int { /* returns number of results */
	lv := time.Duration(L.ToInt(1))
	time.Sleep(lv * time.Microsecond)
	return 0
}

// AvashSetVar sets a variable to a string, necessary because `varstore set` can't deal with spaces yet
func AvashSetVar(L *lua.LState) int {
	log := cfg.Config.Log
	varscope := L.ToString(1)
	varname := L.ToString(2)
	varvalue := L.ToString(3)
	if varscope == "" || varname == "" || varvalue == "" {
		log.Error("Error: AvashSetVar provided insufficient number of arguments, expected 3")
		return 0
	}
	if store, err := AvashVars.Get(varscope); err == nil {
		store.Set(varname, varvalue)
	} else {
		log.Error("Error: AvashSetVar scope not found: " + varscope)
	}
	return 0
}

// AvashCall hooks avash calls into scripts
func AvashCall(L *lua.LState) int { /* returns number of results */
	lv := L.ToString(1) /* get argument */
	cmd, flags, err := AvalancheShell.root.Find(strings.Fields(lv))
	if err != nil {
		AvalancheShell.rl.Terminal.Write([]byte(err.Error()))
	}
	AvalancheShell.addHistory(cmd, flags)
	cmd.ParseFlags(flags)
	captureDone := capture()
	cmd.Run(cmd, flags)
	capturedOutout, err := captureDone()
	log := cfg.Config.Log
	if err != nil {
		L.Push(lua.LString("Error: Unable to execute in capture: " + err.Error()))
		log.Error("Error: Unable to execute in capture: " + err.Error())
		log.Error("Captured Output: " + capturedOutout)
		return 1
	}
	L.Push(lua.LString(strings.TrimSpace(capturedOutout))) /* push result */
	return 1                                               /* number of results */
}

/*
// AvashCoroutine launches a simple Lua coroutine in Golang
func AvashCoroutine(L *lua.LState) int {
	coname := L.ToString(1)  // get coroutine name
	coscope := L.ToString(2) // get coroutine scope
	if coname == "" {
		fmt.Println("Error: AvashCoroutine provided insufficient number of arguments, expected at least 2")
		return 0
	}
	store, err := AvashVars.Get(coscope)
	if err != nil {
		fmt.Printf("Error: AvashCoroutine can't find output scope.\n")
		return 0
	}
	fn := L.GetGlobal(coname).(*lua.LFunction) // get coroutine function
	fnargs := []lua.LValue{}                   // get coroutine arguments
	paramLen := L.GetTop()                     // get top of the parameter stack
	// skip argument #1&2 (coname, coscope) and gather other arguments
	for i := 3; i <= paramLen; i++ {
		fnargs = append(fnargs, L.Get(i))
	}
	// coroutine needs to start here
	go gocoro(L, store, fn, fnargs)
	return 0
}

// separated for readability, coroutine used in AvashCoroutine
func gocoro(rootL *lua.LState, store varScope, fn *lua.LFunction, fnargs []lua.LValue) {
	co, cocancel := rootL.NewThread() // create a thread for the coroutine
	defer cocancel()
	store.Set("state", "yield")
	state, err, values := rootL.Resume(co, fn, fnargs...)
	for {
		if state == lua.ResumeError {
			stateErr(store, "Error: Coroutine state has error: %s\n", err.Error())
			break
		}

		if state == lua.ResumeOK {
			v, err := json.MarshalIndent(values, "", "  ")
			if err != nil {
				stateErr(store, "Error: Coroutine cannot marshal values: %s\n", err.Error())
				break
			}
			store.Set("state", "ok")
			store.Set("values", string(v))
			break
		}
		state, err, values = rootL.Resume(co, fn)
	}
}

func stateErr(store *varScope, estr string, args ...interface{}) {
	e := fmt.Sprintf(estr, args...)
	fmt.Print(e)
	store.Set("state", "error")
	store.Set("value", e)
}
*/
