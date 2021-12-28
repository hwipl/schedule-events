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
	serverMode   = false
)

// parseCommandLine parses the command line arguments
func parseCommandLine() {
	// set command line arguments
	flag.StringVar(&commandsFile, "commands", commandsFile,
		"read commands from `file`")
	flag.StringVar(&serverAddr, "address", serverAddr,
		"listen on or connect to `addr`")
	flag.BoolVar(&serverMode, "server", serverMode, "run as server")
	flag.Parse()

	// parse address
	if serverAddr == "" {
		log.Fatal("no address specified")
	}

	// parse commands file
	if serverMode {
		if commandsFile == "" {
			log.Fatal("no commands file specified")
		}
		err := command.CommandsFromJSON(commandsFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Run is the main entry point
func Run() {
	parseCommandLine()
	if serverMode {
		server.Run(serverAddr)
	}
}
