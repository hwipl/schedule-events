package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/hwipl/schedule-events/internal/command"
	"github.com/hwipl/schedule-events/internal/event"
)

const (
	// maxEventPostLength is the maximum content length of an
	// event post request
	maxEventPostLength = 512
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

// handleEventsPost handles a client "events" POST request
func handleEventsPost(w http.ResponseWriter, r *http.Request) {
	// TODO: add error replies?
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("invalid content type")
		return
	}
	if r.ContentLength <= 0 || r.ContentLength > maxEventPostLength {
		log.Println("invalid content length")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	event, err := event.NewFromJSON(body)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: do something with the event
	log.Println(event)
}

// handleEvents handles a client "events" request
func handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleEventsGet(w, r)
	case http.MethodPost:
		handleEventsPost(w, r)
	}
}

// Run starts the server listening on addr
func Run(addr string) {
	log.Println("Starting server listening on:", addr)

	// schedule all events
	for _, e := range event.List() {
		go e.Schedule()
	}

	// start http server
	http.HandleFunc("/commands/", handleCommands)
	http.HandleFunc("/events/", handleEvents)

	log.Fatal(http.ListenAndServe(addr, nil))
}
