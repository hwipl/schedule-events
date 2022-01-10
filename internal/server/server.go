package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hwipl/schedule-events/internal/command"
	"github.com/hwipl/schedule-events/internal/event"
)

const (
	// maxEventPostLength is the maximum content length of an
	// event post request
	maxEventPostLength = 512
)

var (
	// server is the http server
	server *http.Server
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
	evt, err := event.NewFromJSON(body)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: add event validation?

	// add and schedule event
	log.Println("Adding new event:", evt.Name)
	if event.Add(evt) {
		go evt.Schedule()
	}
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

// handleStatusGet handles a client "status" GET request
func handleStatusGet(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Status: OK\n")
	if err != nil {
		log.Println(err)
	}
}

// handleStatusPost handles a client "status" GET request
func handleStatusPost(w http.ResponseWriter, r *http.Request) {
	switch html.EscapeString(r.URL.Path) {
	case "/status/shutdown":
		Shutdown()
	case "/status/stop":
		Stop()
	}
}

// handleStatus handles a client "status" request
func handleStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatusGet(w, r)
	case http.MethodPost:
		handleStatusPost(w, r)
	}
}

// Shutdown shuts the server down
func Shutdown() {
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

// Stop stops all events on the server
func Stop() {
	for _, e := range event.List() {
		e.Stop()
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
	http.HandleFunc("/status/", handleStatus)

	server = &http.Server{Addr: addr}
	log.Println(server.ListenAndServe())

	// server stopped, stop all events
	Stop()
	time.Sleep(1 * time.Second) // TODO: improve
}
