// Package emission provides an event emitter.
package emission

import (
  "fmt"
  "strconv"
  "sync"
  "time"
)

// Maximum number of listeners an event can have.
const DEFAULT_MAX_LISTENERS = 10

type listener func(...interface{})

type event struct {
  listeners []listener
}

type Emitter struct {
  events       map[string]*event
  maxListeners int
  mutex        sync.Mutex
}

// AddListener adds the listener function (fn) to the Emitter's event (e)
// listener array.
func (emitter *Emitter) AddListener(e string, fn listener) *Emitter {
  emitter.mutex.Lock()
  defer emitter.mutex.Unlock()
  if nil == fn {
    return emitter
  }
  _, ok := emitter.events[e]
  if !ok {
    emitter.events[e] = &event{[]listener{}}
  }
  event := emitter.events[e]
  if emitter.maxListeners >= len(event.listeners)+1 {
    event.listeners = append(event.listeners, fn)
  }
  return emitter
}

// On is an alias method for AddListener.
func (emitter *Emitter) On(e string, fn listener) *Emitter {
  return emitter.AddListener(e, fn)
}

// RemoveListener loops through an Emitter's events and listeners, comparing
// the string value of the given listener function (fn) since go
// does not allow you to compare functions.  If a match is found,
// it is removed from the event's listeners array.
func (emitter *Emitter) RemoveListener(e string, fn listener) *Emitter {
  emitter.mutex.Lock()
  defer emitter.mutex.Unlock()
  ev, ok := emitter.events[e]
  if ok {
    for i, l := range ev.listeners {
      if fmt.Sprintf("%v", l) == fmt.Sprintf("%v", fn) {
        ev.listeners = append(ev.listeners[:i], ev.listeners[i+1:]...)
      }
    }
  }
  return emitter
}

// Off is an alias method for RemoveListener.
func (emitter *Emitter) Off(e string, fn listener) *Emitter {
  return emitter.RemoveListener(e, fn)
}

// Once adds a listener function (fn) to an event (e) that will run a maximum of one time
// before being removed from it's listener array.
func (emitter *Emitter) Once(e string, fn listener) *Emitter {
  if nil == fn {
    return emitter
  }
  var run listener
  run = func(args ...interface{}) {
    fn(args...)
    emitter.RemoveListener(e, run)
  }
  emitter.AddListener(e, run)
  return emitter
}

// Throttled is an event handler `fn` for event `e` that can only be called once within a given
// duration `interval` (in milliseconds).
//
// TODO: A throttled function currently has no way of being removed.
func (emitter *Emitter) Throttled(e string, interval time.Duration,
  fn listener) *Emitter {
  var init int64
  emitter.AddListener(e, func(args ...interface{}) {
    if init == 0 {
      init = emitter.timestamp()
    } else if (emitter.timestamp() - (int64(interval) / 1e6)) <= init {
      return
    }
    init = emitter.timestamp()
    fn(args...)
  })
  return emitter
}

// Emit triggers an event (e), passing along arguments (args) to each of the event's
// listeners.  Each listener function is ran as a go routine.
func (emitter *Emitter) Emit(e string, args ...interface{}) *Emitter {
  if _, ok := emitter.events[e]; !ok {
    return emitter
  }
  var wg sync.WaitGroup
  var listeners = emitter.events[e].listeners
  wg.Add(len(listeners))
  for _, fn := range listeners {
    go func(fn listener) {
      defer wg.Done()
      fn(args...)
    }(fn)
  }
  wg.Wait()
  return emitter
}

// SetMaxListeners sets an emitters maximum listeners per event.
func (emitter *Emitter) SetMaxListeners(max int) *Emitter {
  emitter.maxListeners = max
  return emitter
}

// timestamp is a helper function to generate an accurate UNIX timestamp.
func (emitter *Emitter) timestamp() int64 {
  t := time.Now()
  ts, _ := strconv.ParseInt(fmt.Sprintf("%d%03d", t.Unix(), t.Nanosecond()/1e6), 10, 64)
  return ts
}

// Returns a pointer to an Emitter struct.
func NewEmitter() *Emitter {
  return &Emitter{
    make(map[string]*event),
    DEFAULT_MAX_LISTENERS,
    sync.Mutex{},
  }
}
