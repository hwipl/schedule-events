package command

import (
	"context"
	"log"
	"os/exec"
	"time"
)

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
