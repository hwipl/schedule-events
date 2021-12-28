package cmd

import (
	"flag"
	"log"

	"github.com/hwipl/schedule-events/internal/command"
	"github.com/hwipl/schedule-events/internal/server"
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	cmdFile := flag.String("commands", "", "read commands from `file`")
	flag.Parse()

	// parse commands file
	if *cmdFile == "" {
		log.Fatal("no commands file specified")
	}
	if err := command.CommandsFromJSON(*cmdFile); err != nil {
		log.Fatal(err)
	}
}

// Run is the main entry point
func Run() {
	parseCommandLine()
	server.Run(":8080")
}
