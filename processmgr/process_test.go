package processmgr

import (
	"sync"
	"testing"
)

func newTestProcess(code uint) (*Process, chan bool) {
	var cmdstr string
	var args []string
	switch code {
	case 0:
		cmdstr = "sleep"
		args = []string{"10"}
	case 1:
		cmdstr = "fake-command"
		args = nil
	}
	return &Process{
		cmdstr: cmdstr,
		args:   args,
		stop:   make(chan bool),
		kill:   make(chan bool),
		fail:   make(chan error),
	}, make(chan bool)
}

// Calls `p.Start` with `d` and returns a `WaitGroup` to block on `p` stopping
func syncStart(p *Process, d chan bool) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Start(d)
	}()
	return &wg
}

func TestProcessStart(t *testing.T) {
	p1, d1 := newTestProcess(0)
	p2, d2 := newTestProcess(1)
	p3, d3 := newTestProcess(0)
	p4, d4 := newTestProcess(0)
	p5, d5 := newTestProcess(0)
	p6, d6 := newTestProcess(1)

	t.Run("GoodExec", func(t *testing.T) {
		go p1.Start(d1)
		<-d1

		t.Logf("%+v", p1)
		if proc := p1.cmd.Process; proc == nil {
			t.Fatalf("P.Cmd.Process returned %v expected not %v", proc, nil)
		} else if running := p1.running; running != true {
			t.Fatalf("P.Running returned %t expected %t", running, true)
		} else if failed := p1.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("BadExec", func(t *testing.T) {
		go p2.Start(d2)
		<-d2

		if proc := p2.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p2.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p2.failed; failed != true {
			t.Fatalf("P.Failed returned %t expected %t", failed, true)
		}
	})
	t.Run("KilledExec", func(t *testing.T) {
		done := make(chan bool)
		go func() {
			p3.Start(d3)
			done <- true
		}()
		<-d3
		p3.cmd.Process.Kill()
		<-done

		if proc := p3.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p3.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p3.failed; failed != true {
			t.Fatalf("P.Failed returned %t expected %t", failed, true)
		}
	})
	t.Run("Running", func(t *testing.T) {
		go p4.Start(d4)
		<-d4
		go p4.Start(d4)
		<-d4

		if proc := p4.cmd.Process; proc == nil {
			t.Fatalf("P.Cmd.Process returned %v expected not %v", proc, nil)
		} else if running := p4.running; running != true {
			t.Fatalf("P.Running returned %t expected %t", running, true)
		} else if failed := p4.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("GoodFailed", func(t *testing.T) {
		p5.failed = true
		go p5.Start(d5)
		<-d5

		if proc := p5.cmd.Process; proc == nil {
			t.Fatalf("P.Cmd.Process returned %v expected not %v", proc, nil)
		} else if running := p5.running; running != true {
			t.Fatalf("P.Running returned %t expected %t", running, true)
		} else if failed := p5.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("BadFailed", func(t *testing.T) {
		p6.failed = true
		go p6.Start(d6)
		<-d6

		if proc := p6.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p6.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p6.failed; failed != true {
			t.Fatalf("P.Failed returned %t expected %t", failed, true)
		}
	})

	if p1.cmd.Process != nil {
		p1.cmd.Process.Kill()
	}
	if p2.cmd.Process != nil {
		p2.cmd.Process.Kill()
	}
	if p3.cmd.Process != nil {
		p3.cmd.Process.Kill()
	}
	if p4.cmd.Process != nil {
		p4.cmd.Process.Kill()
	}
	if p5.cmd.Process != nil {
		p5.cmd.Process.Kill()
	}
	if p6.cmd.Process != nil {
		p6.cmd.Process.Kill()
	}
}

func TestProcessStop(t *testing.T) {
	p1, d1 := newTestProcess(0)
	p2, d2 := newTestProcess(0)
	p3, d3 := newTestProcess(1)

	t.Run("Running", func(t *testing.T) {
		wg := syncStart(p1, d1)
		<-d1
		p1.Stop()
		wg.Wait()

		if proc := p1.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p1.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p1.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("Stopped", func(t *testing.T) {
		wg := syncStart(p2, d2)
		<-d2
		p2.Stop()
		p2.Stop()
		wg.Wait()

		if proc := p2.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p2.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p2.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("Failed", func(t *testing.T) {
		wg := syncStart(p3, d3)
		<-d3
		p3.Stop()
		wg.Wait()

		if proc := p3.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p3.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p3.failed; failed != true {
			t.Fatalf("P.Failed returned %t expected %t", failed, true)
		}
	})

	if p1.cmd.Process != nil {
		p1.cmd.Process.Kill()
	}
	if p2.cmd.Process != nil {
		p2.cmd.Process.Kill()
	}
	if p3.cmd.Process != nil {
		p3.cmd.Process.Kill()
	}
}

func TestProcessKill(t *testing.T) {
	p1, d1 := newTestProcess(0)
	p2, d2 := newTestProcess(0)
	p3, d3 := newTestProcess(1)

	t.Run("Running", func(t *testing.T) {
		wg := syncStart(p1, d1)
		<-d1
		p1.Kill()
		wg.Wait()

		if proc := p1.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p1.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p1.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("Stopped", func(t *testing.T) {
		wg := syncStart(p2, d2)
		<-d2
		p2.Stop()
		p2.Kill()
		wg.Wait()

		if proc := p2.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p2.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p2.failed; failed != false {
			t.Fatalf("P.Failed returned %t expected %t", failed, false)
		}
	})
	t.Run("Failed", func(t *testing.T) {
		wg := syncStart(p3, d3)
		<-d3
		p3.Kill()
		wg.Wait()

		if proc := p3.cmd.Process; proc != nil {
			t.Fatalf("P.Cmd.Process returned %v expected %v", proc, nil)
		} else if running := p3.running; running != false {
			t.Fatalf("P.Running returned %t expected %t", running, false)
		} else if failed := p3.failed; failed != true {
			t.Fatalf("P.Failed returned %t expected %t", failed, true)
		}
	})

	if p1.cmd.Process != nil {
		p1.cmd.Process.Kill()
	}
	if p2.cmd.Process != nil {
		p2.cmd.Process.Kill()
	}
	if p3.cmd.Process != nil {
		p3.cmd.Process.Kill()
	}
}