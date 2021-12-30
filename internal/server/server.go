package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hwipl/schedule-events/internal/command"
	"github.com/hwipl/schedule-events/internal/event"
)

// handleCommandsGet handles a client "commands" GET request
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

// handleEventsGet handles a client "events" GET request
func handleEventsGet(w http.ResponseWriter, r *http.Request) {
	events := event.List()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}

// handleEvents handles a client "events" request
func handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleEventsGet(w, r)
	}
}

// Runs starts the server listening on addr
func Run(addr string) {
	log.Println("Starting server listening on:", addr)

	// schedule all events
	for _, e := range event.List() {
		log.Println("Scheduling event:", e.Name)
		e.Schedule()
	}

	// start http server
	http.HandleFunc("/commands/", handleCommands)
	http.HandleFunc("/events/", handleEvents)

	log.Fatal(http.ListenAndServe(addr, nil))
}
