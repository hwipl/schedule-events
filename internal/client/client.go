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

// getCommandsAll retrieves all commands from the server and prints them
func getCommandsAll(addr string) {
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

// getCommandsOne retrieves the command with name from the server and prints it
func getCommandsOne(addr, name string) {
	// get command from server
	url := fmt.Sprintf("http://%s/commands/%s", addr, name)
	body := get(url)

	// make sure it's a valid json Command
	cmd := command.Command{}
	if err := json.Unmarshal(body, &cmd); err != nil {
		log.Fatal(err)
	}

	// print as indented json
	var out bytes.Buffer
	json.Indent(&out, body, "", "    ")
	fmt.Println(&out)
}

// getCommands retrieves the command list from the server and prints it
func getCommands(addr string) {
	cmds := command.List()
	if len(cmds) == 0 {
		log.Println("Getting all commands from server")
		getCommandsAll(addr)
		return
	}
	for _, cmd := range cmds {
		log.Println("Getting command from server:", cmd.Name)
		getCommandsOne(addr, cmd.Name)
	}
}

// getEventsAll retrieves all events from the server and prints them
func getEventsAll(addr string) {
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

// getEventsOne retrieves the event with name from the server and prints it
func getEventsOne(addr, name string) {
	// get event from server
	url := fmt.Sprintf("http://%s/events/%s", addr, name)
	body := get(url)

	// make sure it's a valid json Event
	evt := event.Event{}
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Fatal(err)
	}

	// print as indented json
	var out bytes.Buffer
	json.Indent(&out, body, "", "    ")
	fmt.Println(&out)
}

// getEvents retrieves the event list from the server and prints it
func getEvents(addr string) {
	evts := event.List()
	if len(evts) == 0 {
		log.Println("Getting all events from server")
		getEventsAll(addr)
		return
	}
	for _, evt := range evts {
		log.Println("Getting event from server:", evt.Name)
		getEventsOne(addr, evt.Name)
	}
}

// getStatus retrieves the status from the server and prints it
func getStatus(addr string) {
	log.Println("Getting status from server")

	url := fmt.Sprintf("http://%s/status", addr)
	body := get(url)
	fmt.Println(string(body))
}

// handleResponse a response discarding the body and checking the status code
func handleResponse(resp *http.Response) {
	defer resp.Body.Close()
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode > 299 {
		log.Fatal(resp.StatusCode)
	}

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
			log.Fatal(err)
		}
		handleResponse(resp)
	}
}

// setStatus sends status request to server
func setStatus(addr, status string) {
	url := fmt.Sprintf("http://%s/status/%s", addr, status)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		log.Fatal(err)
	}
	handleResponse(resp)
}

// shutdown sends a shutdown request to the server
func shutdown(addr string) {
	log.Println("Sending shutdown to server")
	setStatus(addr, "shutdown")
}

// stop sends a stop request to the server
func stop(addr string) {
	log.Println("Sending stop to server")
	setStatus(addr, "stop")
}

// delEvents deletes events on the server
func delEvents(addr string) {
	log.Println("Deleting events on server")

	for _, e := range event.List() {
		log.Println("Deleting event:", e.Name)

		url := fmt.Sprintf("http://%s/events/%s", addr, e.Name)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		handleResponse(resp)
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
	case "delete-events":
		delEvents(addr)
	case "get-status":
		getStatus(addr)
	case "shutdown":
		shutdown(addr)
	case "stop":
		stop(addr)
	case "":
		getEvents(addr)
	default:
		log.Fatal("invalid operation: ", op)
	}
}
