package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hwipl/schedule-events/internal/event"
)

// getEvents retrieves the event list from the server and prints it
func getEvents(addr string) {
	log.Println("Getting events from server")

	// get events from server
	url := fmt.Sprintf("http://%s/events", addr)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode > 299 {
		log.Fatal(resp.StatusCode)
	}

	// make sure it's a valid json Event array
	events := []*event.Event{}
	if err := json.Unmarshal(body, &events); err != nil {
		log.Fatal(err)
	}

	// print as indented json
	var out bytes.Buffer
	json.Indent(&out, body, "", "    ")
	fmt.Println(&out)
}

// Run starts the client connecting to addr and executing op
func Run(addr, op string) {
	log.Println("Starting client connecting to:", addr)
	switch op {
	default:
		getEvents(addr)
	}
}
