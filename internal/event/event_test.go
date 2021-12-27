package event

import (
	"testing"
	"time"
)

// TestNextWait tests getting the next wait time
func TestNextWait(t *testing.T) {
	test := func(want, got time.Duration) {
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// fixed zero wait time
	e1 := Event{
		WaitMin: 0,
		WaitMax: 0,
	}
	test(0, e1.nextWait())

	// fixed 1s wait time
	e2 := Event{
		WaitMin: 1 * time.Second,
		WaitMax: 1 * time.Second,
	}
	test(e2.WaitMin, e2.nextWait())

	// wait time < 0, treated as 0
	e3 := Event{
		WaitMin: -1,
		WaitMax: -1,
	}
	test(0, e3.nextWait())

	// min < max, treated as max = min
	e4 := Event{
		WaitMin: 4 * time.Second,
		WaitMax: 0,
	}
	test(e4.WaitMin, e4.nextWait())

	// wait time between 0 and 1s
	e5 := Event{
		WaitMin: 5 * time.Second,
		WaitMax: 6 * time.Second,
	}
	got := e5.nextWait()
	if got < e5.WaitMin || got > e5.WaitMax {
		t.Errorf("got %v, want 0-1s", got)
	}
}

// TestScheduleWait tests scheduling the event after a wait time
func TestScheduleWait(t *testing.T) {
	// expired event (StopDate: 0)
	e1 := &Event{}
	e1.scheduleWait(-1)
	e1.scheduleWait(0)
	e1.scheduleWait(1 * time.Second)

	// event that will expire in 5 seconds from now
	e2 := &Event{
		StartDate: time.Now(),
		StopDate:  time.Now().Add(5 * time.Second),
	}
	e2.scheduleWait(-1)
	e2.scheduleWait(0)
	e2.scheduleWait(1 * time.Second)
	e2.scheduleWait(4 * time.Second) // expired
	e2.scheduleWait(5 * time.Second) // expired
}

// TestJSON tests conversion from and to json
func TestJSON(t *testing.T) {
	e1 := &Event{
		Command:   "test",
		StartDate: time.Now().Local(),
		StopDate:  time.Now().Add(10 * time.Second).Local(),
		Timeout:   30 * time.Second,
		Periodic:  true,
		WaitMin:   0,
		WaitMax:   time.Second,
	}

	// convert to json
	b, err := e1.JSON()
	if err != nil {
		t.Errorf("got %v, want nil", err)
	}

	// convert from json
	e2, err := NewFromJSON(b)
	if err != nil {
		t.Errorf("got %v, want nil", err)
	}
	if *e1 != *e2 {
		t.Errorf("got e1 != e2, want e1 == e2\ne1: %#v\ne2: %#v",
			e1, e2)
	}
}
