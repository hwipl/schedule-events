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
	serverAddr   = ":8080"
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// set command line arguments
	flag.StringVar(&commandsFile, "commands", commandsFile,
		"read commands from `file`")
	flag.StringVar(&serverAddr, "address", serverAddr,
		"listen on or connect to `addr`")
	flag.Parse()

	// parse address
	if serverAddr == "" {
		log.Fatal("no address specified")
	}

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
	server.Run(serverAddr)
}
