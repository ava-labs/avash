package processmgr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProcess(t *testing.T) {
	pm := ProcessManager{
		processes: make(map[string]*Process),
	}

	cmd0 := "cmd0"
	proctype0 := "type0"
	args0 := []string{"arg0"}
	name0 := "test0"
	metadata0 := "data0"

	cmd1 := "cmd1"
	proctype1 := "type1"
	args1 := []string{"arg1"}
	name1 := "test1"
	metadata1 := "data1"

	testProcInitWith := func(t *testing.T, p *Process, cmd string, proctype string, args []string, name string, metadata string) {
		if p.cmdstr != cmd {
			t.Fatalf("P.Cmdstr returned %s expected %s", p.cmdstr, cmd)
		} else if p.proctype != proctype {
			t.Fatalf("P.Proctype returned %s expected %s", p.proctype, proctype)
		} else if !reflect.DeepEqual(p.args, args) {
			t.Fatalf("P.Args returned %v expected %v", p.args, args)
		} else if p.name != name {
			t.Fatalf("P.Name returned %s expected %s", p.name, name)
		} else if p.metadata != metadata {
			t.Fatalf("P.Metadata returned %s expected %s", p.metadata, metadata)
		}
	}

	t.Run("InitAdd", func(t *testing.T) {
		pm.AddProcess(cmd0, proctype0, args0, name0, metadata0, nil, nil, nil)

		if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}

		testProcInitWith(t, pm.processes[name0], cmd0, proctype0, args0, name0, metadata0)
	})
	t.Run("DuplicateAdd", func(t *testing.T) {
		pm.AddProcess("fake", "fake", []string{"fake"}, name0, "fake", nil, nil, nil)

		if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}

		testProcInitWith(t, pm.processes[name0], cmd0, proctype0, args0, name0, metadata0)
	})
	t.Run("NewAdd", func(t *testing.T) {
		pm.AddProcess(cmd1, proctype1, args1, name1, metadata1, nil, nil, nil)

		if count := len(pm.processes); count != 2 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 2)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		} else if _, ok := pm.processes[name1]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name1)
		}

		testProcInitWith(t, pm.processes[name1], cmd1, proctype1, args1, name1, metadata1)
	})
	t.Run("EmptyNameAdd", func(t *testing.T) {
		pm.AddProcess("fake", "fake", []string{"fake"}, "", "fake", nil, nil, nil)

		if count := len(pm.processes); count != 2 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 2)
		} else if _, ok := pm.processes[""]; ok {
			t.Fatalf("PM.Processes contains process with empty name")
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		} else if _, ok := pm.processes[name1]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name1)
		}
	})
}

func TestRemoveProcess(t *testing.T) {
	pm := ProcessManager{
		processes: make(map[string]*Process),
	}
	name0 := "test0"
	name1 := "test1"
	name2 := "fake"
	pm.AddProcess("cmd", "fake-cmd", []string{"arg"}, name0, "data", nil, nil, nil)
	pm.AddProcess("cmd", "fake-cmd", []string{"arg"}, name1, "data", nil, nil, nil)

	t.Run("ExistingProc", func(t *testing.T) {
		pm.RemoveProcess(name0)

		if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; ok {
			t.Fatalf("PM.Processes contains %s", name0)
		} else if _, ok := pm.processes[name1]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name1)
		}
	})
	t.Run("MissingProc", func(t *testing.T) {
		err := pm.RemoveProcess(name2)

		if dneErr := fmt.Sprintf("Process does not exist, cannot remove: %s", name2); err.Error() != dneErr {
			t.Fatalf("PM.RemoveProcess returned %v expected %v", err, dneErr)
		} else if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name1]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name1)
		}
	})
}

func TestStartProcess(t *testing.T) {
	pm := ProcessManager{
		processes: make(map[string]*Process),
	}
	name0 := "test"
	name1 := "fake"
	pm.AddProcess("cmd", "fake-cmd", []string{"arg"}, name0, "data", nil, nil, nil)

	t.Run("ExistingProc", func(t *testing.T) {
		err := pm.StartProcess(name0)

		if err != nil {
			t.Fatalf("PM.StartProcess returned %v expected %v", err, nil)
		} else if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}
	})
	t.Run("MissingProc", func(t *testing.T) {
		err := pm.StartProcess(name1)

		if dneErr := fmt.Sprintf("Process does not exist, cannot start: %s", name1); err.Error() != dneErr {
			t.Fatalf("PM.StartProcess returned %v expected %v", err, dneErr)
		} else if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}
	})
}

func TestRemoveAllProcesses(t *testing.T) {
	pm := ProcessManager{
		processes: make(map[string]*Process),
	}

	name := "procname-"
	for i := 0; i < 5; i++ {
		_ = pm.AddProcess("cmd", "fake-cmd", []string{"arg"}, name+string(rune(i)), "data", nil, nil, nil)
	}

	assert.Len(t, pm.processes, 5)
	pm.RemoveAllProcesses()
	assert.Len(t, pm.processes, 0)
}

func TestStopProcess(t *testing.T) {
	pm := ProcessManager{
		processes: make(map[string]*Process),
	}
	name0 := "test"
	name1 := "fake"
	pm.AddProcess("cmd", "fake-cmd", []string{"arg"}, name0, "data", nil, nil, nil)

	t.Run("ExistingProc", func(t *testing.T) {
		pm.StopProcess(name0)

		if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}
	})
	t.Run("MissingProc", func(t *testing.T) {
		err := pm.StopProcess(name1)

		if dneErr := fmt.Sprintf("Process does not exist, cannot stop: %s", name1); err.Error() != dneErr {
			t.Fatalf("PM.StopProcess returned %v expected %v", err, dneErr)
		} else if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}
	})
}

func TestKillProcess(t *testing.T) {
	pm := ProcessManager{
		processes: make(map[string]*Process),
	}
	name0 := "test"
	name1 := "fake"
	pm.AddProcess("cmd", "fake-cmd", []string{"arg"}, name0, "data", nil, nil, nil)

	t.Run("ExistingProc", func(t *testing.T) {
		pm.KillProcess(name0)

		if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}
	})
	t.Run("MissingProc", func(t *testing.T) {
		err := pm.KillProcess(name1)

		if dneErr := fmt.Sprintf("Process does not exist, cannot kill: %s", name1); err.Error() != dneErr {
			t.Fatalf("PM.KillProcess returned %v expected %v", err, dneErr)
		} else if count := len(pm.processes); count != 1 {
			t.Fatalf("PM.Processes has length %d expected %d", count, 1)
		} else if _, ok := pm.processes[name0]; !ok {
			t.Fatalf("PM.Processes does not contain %s", name0)
		}
	})
}
