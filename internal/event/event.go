package event

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/hwipl/schedule-events/internal/command"
)

var (
	// events stores a list of all events
	events = newEventList()
)

// eventList is a list of events identified by their name
type eventList struct {
	sync.Mutex
	m map[string]*Event
}

// Add adds event to the event list
func (e *eventList) Add(event *Event) {
	e.Lock()
	defer e.Unlock()

	e.m[event.Name] = event
}

// Get returns the event identified by its name
func (e *eventList) Get(name string) *Event {
	e.Lock()
	defer e.Unlock()

	return e.m[name]
}

// List returns all events sorted by their name
func (e *eventList) List() []*Event {
	e.Lock()
	defer e.Unlock()

	evts := []*Event{}
	for _, evt := range e.m {
		evts = append(evts, evt)
	}
	sort.Slice(evts, func(i, j int) bool {
		return evts[i].Name <
			evts[j].Name
	})
	return evts
}

// newEventList returns a new eventList
func newEventList() *eventList {
	return &eventList{
		m: make(map[string]*Event),
	}
}

// Event is an event that can be scheduled
type Event struct {
	Name      string
	Command   string
	StartDate time.Time
	StopDate  time.Time
	Timeout   time.Duration
	Periodic  bool
	WaitMin   time.Duration
	WaitMax   time.Duration
	done      bool
}

// Run executes the event's command once
func (e *Event) Run() {
	c := command.Get(e.Command)
	if c == nil {
		log.Println("command not found")
		return
	}
	err := c.Run()
	if err != nil {
		log.Println(err)
	}
}

// nextWait returns the next wait duration for the event
func (e *Event) nextWait() time.Duration {
	// get minimum and maximum wait times
	min, max := e.WaitMin, e.WaitMax
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if max < min {
		max = min
	}

	// get next wait time, non-random case
	if min == max {
		return min
	}

	// get next wait time, random in milliseconds granularity
	diff := max.Milliseconds() - min.Milliseconds()
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	t := min.Milliseconds() + r.Int63n(diff)

	return time.Duration(t) * time.Millisecond
}

// scheduleWait schedules the event after the wait duration
func (e *Event) scheduleWait(wait time.Duration) {
	if wait < 0 {
		wait = 0
	}
	if time.Now().Add(wait).After(e.StopDate) {
		log.Println("Event done")
		e.done = true
		return
	}
	time.Sleep(wait)
	e.Run()
}

// Schedule schedules the event for execution
func (e *Event) Schedule() {
	// schedule first execution
	wait := e.StartDate.Sub(time.Now())
	e.scheduleWait(wait)

	// schedule periodic executions
	for e.Periodic && !e.done {
		wait = e.nextWait()
		e.scheduleWait(wait)
	}
}

// JSON returns the event as json
func (e *Event) JSON() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NewFromJSON parses an event from json
func NewFromJSON(b []byte) (*Event, error) {
	e := &Event{}
	err := json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Add adds event to the event list
func Add(event *Event) {
	events.Add(event)
}

// Get returns the event identified by name
func Get(name string) *Event {
	return events.Get(name)
}

// List returns all events in the event list
func List() []*Event {
	return events.List()
}

// EventsFromJSON loads events from the json file in path and adds them to
// the event list
func EventsFromJSON(path string) error {
	// read file
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// parse events
	evts := []*Event{}
	if err := json.Unmarshal(file, &evts); err != nil {
		return err
	}

	// add events to event list
	for _, e := range evts {
		Add(e)
	}
	return nil
}
