package command

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"sort"
	"sync"
	"time"
)

var (
	// commands stores a list of all commands
	commands = newCommandList()
)

// commandList is a list of commands identified by their name
type commandList struct {
	sync.Mutex
	m map[string]*Command
}

// Add adds command to the command list
func (c *commandList) Add(command *Command) {
	c.Lock()
	defer c.Unlock()

	c.m[command.Name] = command
}

// Get returns the command identified by its name
func (c *commandList) Get(name string) *Command {
	c.Lock()
	defer c.Unlock()

	return c.m[name]
}

// List returns all commands sorted by their name
func (c *commandList) List() []*Command {
	c.Lock()
	defer c.Unlock()

	cmds := []*Command{}
	for _, c := range c.m {
		cmds = append(cmds, c)
	}
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Name <
			cmds[j].Name
	})
	return cmds
}

// newCommandList returns a new commandList
func newCommandList() *commandList {
	return &commandList{
		m: make(map[string]*Command),
	}
}

// Command is an executable command
type Command struct {
	Name       string
	Executable string
	Arguments  []string
	Timeout    time.Duration
}

// Run executes the command
func (c *Command) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.Executable, c.Arguments...)
	return cmd.Run()
}

// Add adds command to the command list
func Add(command *Command) {
	commands.Add(command)
}

// Get returns the command identified by name
func Get(name string) *Command {
	return commands.Get(name)
}

// List returns all commands in the command list
func List() []*Command {
	return commands.List()
}

// CommandsFromJSON loads commands from the json file in path and adds them to
// the command list
func CommandsFromJSON(path string) error {
	// read file
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// parse commands
	cmds := []*Command{}
	if err := json.Unmarshal(file, &cmds); err != nil {
		return err
	}

	// add commands to command list
	for _, c := range cmds {
		Add(c)
	}
	return nil
}
