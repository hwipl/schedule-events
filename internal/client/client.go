package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hwipl/schedule-events/internal/command"
	"github.com/hwipl/schedule-events/internal/event"
)

// get retrieves content from url
func get(url string) []byte {
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
	return body
}

// getCommands retrieves the command list from the server and prints it
func getCommands(addr string) {
	log.Println("Getting commands from server")

	// get commands from server
	url := fmt.Sprintf("http://%s/commands", addr)
	body := get(url)

	// make sure it's a valid json Command array
	commands := []*command.Command{}
	if err := json.Unmarshal(body, &commands); err != nil {
		log.Fatal(err)
	}

	// print as indented json
	var out bytes.Buffer
	json.Indent(&out, body, "", "    ")
	fmt.Println(&out)
}

// getEvents retrieves the event list from the server and prints it
func getEvents(addr string) {
	log.Println("Getting events from server")

	// get events from server
	url := fmt.Sprintf("http://%s/events", addr)
	body := get(url)

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

// setEvents sends the client's event list to the server for scheduling
func setEvents(addr string) {
	log.Println("Sending events to server")

	// send events to server
	url := fmt.Sprintf("http://%s/events/", addr)
	for _, e := range event.List() {
		log.Println("Sending event:", e.Name)

		b, err := e.JSON()
		if err != nil {
			log.Fatal(err)
		}
		r := bytes.NewReader(b)
		resp, err := http.Post(url, "application/json", r)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
		io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			log.Println(err)
		}
		if resp.StatusCode > 299 {
			log.Println(resp.StatusCode)
		}
	}
}

// Run starts the client connecting to addr and executing op
func Run(addr, op string) {
	log.Println("Starting client connecting to:", addr)
	switch op {
	case "get-commands":
		getCommands(addr)
	case "get-events":
		getEvents(addr)
	case "set-events":
		setEvents(addr)
	case "":
		getEvents(addr)
	default:
		log.Fatal("invalid operation: ", op)
	}
}
