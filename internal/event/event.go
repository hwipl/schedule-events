package event

import (
	"context"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

// Event is an event that can be scheduled
type Event struct {
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
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.Command)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}

// getWait returns the next wait duration for the event
func (e *Event) getWait() time.Duration {
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
		wait = e.getWait()
		e.scheduleWait(wait)
	}
}
