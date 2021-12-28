package cmd

import (
	"flag"
	"log"

	"github.com/hwipl/schedule-events/internal/command"
	"github.com/hwipl/schedule-events/internal/server"
)

var (
	// parsed command line arguments
	commandsFile = "config.json"
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	flag.StringVar(&commandsFile, "commands", commandsFile,
		"read commands from `file`")
	flag.Parse()

	// parse commands file
	if commandsFile == "" {
		log.Fatal("no commands file specified")
	}
	if err := command.CommandsFromJSON(commandsFile); err != nil {
		log.Fatal(err)
	}
}

// Run is the main entry point
func Run() {
	parseCommandLine()
	server.Run(":8080")
}
