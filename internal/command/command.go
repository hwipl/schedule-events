package command

import (
	"context"
	"log"
	"os/exec"
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
func (c *Command) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.Executable, c.Arguments...)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}

// Add add command to the command list
func Add(command *Command) {
	commands.Add(command)
}

// Get returns the command identified by name
func Get(name string) *Command {
	return commands.Get(name)
}
