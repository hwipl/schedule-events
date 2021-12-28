package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hwipl/schedule-events/internal/command"
)

// handleCommands handles a client "commands" GET request
func handleCommandsGet(w http.ResponseWriter, r *http.Request) {
	cmds := command.List()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(cmds)
	if err != nil {
		log.Println(err)
	}
}

// handleCommands handles a client "commands" request
func handleCommands(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleCommandsGet(w, r)
	}
}

// Runs starts the server
func Run() {
	http.HandleFunc("/commands/", handleCommands)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
