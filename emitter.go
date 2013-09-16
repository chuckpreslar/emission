// Package emission provides an event emitter.
package emission

import (
  "reflect"
  "sync"
)

const (
  DEFAULT_MAX_LISTENERS = 10
)

type Listener func(...interface{})

// Equals compares the reflect Value of the listener
// to another.
func (l Listener) Equals(o Listener) bool {
  return reflect.ValueOf(l) == reflect.ValueOf(o)
}

// Execute calls `l(args...)`.
func (l Listener) Execute(args ...interface{}) {
  l(args...)
}

type Event struct {
  listeners []Listener // The registered Listener functions of the event.
}

type Emitter struct {
  events       map[string]*Event // A map of string to a pointer for an Event.
  maxListeners int               // Maximum number of Listener functions per event for an emitter.
}

// AddListener adds the Listener function (fn) to the Emitter's event (e)
// listener array.
func (emitter *Emitter) AddListener(e string, fn Listener) *Emitter {
  if nil == fn {
    return emitter
  }

  var (
    event *Event
    ok    bool
  )

  if event, ok = emitter.events[e]; !ok {
    event = &Event{[]Listener{}}
    emitter.events[e] = event
  }

  if emitter.maxListeners == -1 || emitter.maxListeners >= len(event.listeners)+1 {
    event.listeners = append(event.listeners, fn)
  }

  return emitter
}

// On is an alias method for AddListener.
func (emitter *Emitter) On(e string, fn Listener) *Emitter {
  return emitter.AddListener(e, fn)
}

// RemoveListener loops through an Emitter's events and listeners, comparing
// the string value of the given Listener function (fn) since go
// does not allow you to compare functions.  If a match is found,
// it is removed from the event's listeners array.
func (emitter *Emitter) RemoveListener(e string, fn Listener) *Emitter {
  if ev, ok := emitter.events[e]; ok {
    for i, l := range ev.listeners {
      if fn.Equals(l) {
        ev.listeners = append(ev.listeners[:i], ev.listeners[i+1:]...)
      }
    }
  }

  return emitter
}

// Off is an alias method for RemoveListener.
func (emitter *Emitter) Off(e string, fn Listener) *Emitter {
  return emitter.RemoveListener(e, fn)
}

// Once adds a Listener function (fn) to an event (e) that will run a maximum of one time
// before being removed from it's listeners array.
func (emitter *Emitter) Once(e string, fn Listener) *Emitter {
  if nil == fn {
    return emitter
  }

  var run Listener

  run = func(args ...interface{}) {
    fn.Execute(args...)
    emitter.RemoveListener(e, run)
  }

  emitter.AddListener(e, run)
  return emitter
}

// Emit triggers an event (e), passing along arguments (args) to each of the event's
// listeners.  Each Listener function is ran as a go routine.
func (emitter *Emitter) Emit(e string, args ...interface{}) *Emitter {
  if _, ok := emitter.events[e]; !ok {
    return emitter
  }

  var wg sync.WaitGroup
  var listeners = emitter.events[e].listeners

  wg.Add(len(listeners))

  for _, fn := range listeners {
    go func(fn Listener) {
      defer wg.Done()
      fn.Execute(args...)
    }(fn)
  }

  wg.Wait()
  return emitter
}

// SetMaxListeners sets an emitters maximum listeners per event.
// If `max` is passed
func (emitter *Emitter) SetMaxListeners(max int) *Emitter {
  emitter.maxListeners = max
  return emitter
}

// Returns a pointer to an Emitter struct.
func NewEmitter() (emitter *Emitter) {
  emitter = new(Emitter)
  emitter.events = make(map[string]*Event)
  emitter.maxListeners = DEFAULT_MAX_LISTENERS
  return
}
