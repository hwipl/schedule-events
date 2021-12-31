package event

import (
	"reflect"
	"testing"
	"time"
)

// TestEventListAdd tests adding events to an eventList
func TestEventListAdd(t *testing.T) {
	// prepare event list, some test events, test function
	evtList := newEventList()
	evt1 := &Event{Name: "evt1"}
	evt2 := &Event{Name: "evt1"} // duplicate name for overwrite test
	evt3 := &Event{Name: "evt3"}
	evt4 := &Event{Name: "evt4"}
	test := func(want, got *Event) {
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// test adding new entry to empty list
	evtList.Add(evt1)
	test(evt1, evtList.Get(evt1.Name))

	// test overwriting existing entry
	evtList.Add(evt2)
	test(evt2, evtList.Get(evt1.Name))

	// test adding more entries
	evtList.Add(evt3)
	evtList.Add(evt4)
	test(evt2, evtList.Get(evt2.Name))
	test(evt3, evtList.Get(evt3.Name))
	test(evt4, evtList.Get(evt4.Name))
}

// TestEventListGet tests getting events from an eventList
func TestEventListGet(t *testing.T) {
	// prepare event list, some test events, test function
	evtList := newEventList()
	evt1 := &Event{Name: "evt1"}
	evt2 := &Event{Name: "evt2"}
	evt3 := &Event{Name: "evt3"}
	test := func(want, got *Event) {
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// test empty list
	test(nil, evtList.Get("does not exist"))

	// test with 1 entry
	evtList.Add(evt1)

	test(evt1, evtList.Get(evt1.Name))
	test(nil, evtList.Get("does not exist"))

	// test with more entries
	evtList.Add(evt2)
	evtList.Add(evt3)

	test(evt1, evtList.Get(evt1.Name))
	test(evt2, evtList.Get(evt2.Name))
	test(evt3, evtList.Get(evt3.Name))
	test(nil, evtList.Get("does not exist"))
}

// TestEventListList tests listing events in an eventList
func TestEventListList(t *testing.T) {
	// prepare event list, some test events, test function
	evtList := newEventList()
	evt1 := &Event{Name: "evt1"}
	evt2 := &Event{Name: "evt2"}
	evt3 := &Event{Name: "evt3"}
	test := func(want, got []*Event) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	// test empty list
	test([]*Event{}, evtList.List())

	// test one element
	evtList.Add(evt1)
	test([]*Event{evt1}, evtList.List())

	// test more elements
	evtList.Add(evt2)
	evtList.Add(evt3)
	test([]*Event{evt1, evt2, evt3}, evtList.List())
}

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
	// event without expiration date
	e1 := &Event{}
	e1.scheduleWait(-1)
	e1.scheduleWait(0)
	e1.scheduleWait(1 * time.Second)

	// expired event
	e2 := &Event{}
	e2.StopDate = e2.StopDate.Add(1) // set valid date in the past
	e2.scheduleWait(-1)
	e2.scheduleWait(0)
	e2.scheduleWait(1 * time.Second)

	// event that will expire in 5 seconds from now
	e3 := &Event{
		StartDate: time.Now(),
		StopDate:  time.Now().Add(5 * time.Second),
	}
	e3.scheduleWait(-1)
	e3.scheduleWait(0)
	e3.scheduleWait(1 * time.Second)
	e3.scheduleWait(4 * time.Second) // expired
	e3.scheduleWait(5 * time.Second) // expired
}

// TestSchedule tests scheduling events
func TestSchedule(t *testing.T) {
	// start time in the past
	e1 := &Event{}
	e1.Schedule()

	// start time now (probably slightly in the past ;))
	e2 := &Event{StartDate: time.Now()}
	e2.Schedule()

	// start time in the future
	e3 := &Event{StartDate: time.Now().Add(time.Second)}
	e3.Schedule()

	// periodic event
	e4 := &Event{
		StopDate: time.Now().Add(time.Second),
		Periodic: true,
		WaitMin:  100 * time.Millisecond,
	}
	e4.Schedule()
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
