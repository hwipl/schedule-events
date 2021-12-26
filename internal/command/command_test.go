package command

import (
	"reflect"
	"testing"
	"time"
)

// TestCommandListAdd tests adding commands to a commandList
func TestCommandListAdd(t *testing.T) {
	// prepare command list, some test commands, test function
	cmdList := newCommandList()
	cmd1 := &Command{Name: "cmd1"}
	cmd2 := &Command{Name: "cmd1"} // duplicate name for overwrite test
	cmd3 := &Command{Name: "cmd3"}
	cmd4 := &Command{Name: "cmd4"}
	test := func(want, got *Command) {
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// test adding new entry to empty list
	cmdList.Add(cmd1)
	test(cmd1, cmdList.Get(cmd1.Name))

	// test overwriting existing entry
	cmdList.Add(cmd2)
	test(cmd2, cmdList.Get(cmd1.Name))

	// test adding more entries
	cmdList.Add(cmd3)
	cmdList.Add(cmd4)
	test(cmd2, cmdList.Get(cmd2.Name))
	test(cmd3, cmdList.Get(cmd3.Name))
	test(cmd4, cmdList.Get(cmd4.Name))
}

// TestCommandListGet tests getting commands from a commandList
func TestCommandListGet(t *testing.T) {
	// prepare command list, some test commands, test function
	cmdList := newCommandList()
	cmd1 := &Command{Name: "cmd1"}
	cmd2 := &Command{Name: "cmd2"}
	cmd3 := &Command{Name: "cmd3"}
	test := func(want, got *Command) {
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// test empty list
	test(nil, cmdList.Get("does not exist"))

	// test with 1 entry
	cmdList.Add(cmd1)

	test(cmd1, cmdList.Get(cmd1.Name))
	test(nil, cmdList.Get("does not exist"))

	// test with more entries
	cmdList.Add(cmd2)
	cmdList.Add(cmd3)

	test(cmd1, cmdList.Get(cmd1.Name))
	test(cmd2, cmdList.Get(cmd2.Name))
	test(cmd3, cmdList.Get(cmd3.Name))
	test(nil, cmdList.Get("does not exist"))
}

// TestCommandListList tests listing commands in a commandList
func TestCommandListList(t *testing.T) {
	// prepare command list, some test commands, test function
	cmdList := newCommandList()
	cmd1 := &Command{Name: "cmd1"}
	cmd2 := &Command{Name: "cmd2"}
	cmd3 := &Command{Name: "cmd3"}
	test := func(want, got []*Command) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// test empty list
	test([]*Command{}, cmdList.List())

	// test one element
	cmdList.Add(cmd1)
	test([]*Command{cmd1}, cmdList.List())

	// test more elements
	cmdList.Add(cmd2)
	cmdList.Add(cmd3)
	test([]*Command{cmd1, cmd2, cmd3}, cmdList.List())
}

// TestCommandRun tests running commands
func TestCommandRun(t *testing.T) {
	cmd1 := &Command{
		Name:       "list",
		Executable: "ls",
	}
	cmd2 := &Command{
		Name:       "sleep",
		Executable: "sleep",
		Arguments:  []string{"1"},
	}

	// test timeout, no args
	if err := cmd1.Run(); err == nil {
		t.Errorf("got %v, want !nil", err)
	}

	// test timeout, with args
	if err := cmd2.Run(); err == nil {
		t.Errorf("got %v, want !nil", err)
	}

	// test successful run, no args
	cmd1.Timeout = 10 * time.Second
	if err := cmd1.Run(); err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	// test successful run, with args
	cmd2.Timeout = 10 * time.Second
	if err := cmd2.Run(); err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
}
