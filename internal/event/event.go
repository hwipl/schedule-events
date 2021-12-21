package event

import (
	"log"
	"time"
)

// Event is an event that can be scheduled
type Event struct {
	Command   string
	StartDate time.Time
	StopDate  time.Time
	Periodic  bool
	Wait      string
	done      bool
}

// Run executes the event's command once
func (e *Event) Run() {
	// TODO: implement it
}

// getWait returns the next wait duration for the event
func (e *Event) getWait() time.Duration {
	// TODO: parse Wait and return proper wait duration
	return 5 * time.Second
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
		wait = e.getWait()
		e.scheduleWait(wait)
	}
}
