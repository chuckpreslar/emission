// Package emission provides an event emitter.
package emission

import (
  "reflect"
  "sync"
  "time"
)

// Maximum number of listeners an event can have.
const DEFAULT_MAX_LISTENERS = 10

type listener func(...interface{})

type throttle struct {
  lastCalled time.Time
  interval   time.Duration
  fn         listener
}

type event struct {
  listeners []listener
  throttles []*throttle
}

type Emitter struct {
  events       map[string]*event
  maxListeners int
  mutex        sync.Mutex
}

func createEvent() *event {
  return &event{[]listener{}, []*throttle{}}
}

// AddListener appends the listener function `fn` to the Emitter's listeners
// for event `e`.
func (emitter *Emitter) AddListener(e string, fn listener) *Emitter {
  emitter.mutex.Lock()
  defer emitter.mutex.Unlock()
  if nil == fn {
    return emitter
  }
  _, ok := emitter.events[e]
  if !ok {
    emitter.events[e] = createEvent()
  }
  event := emitter.events[e]
  if emitter.maxListeners >= len(event.listeners)+len(event.throttles)+1 {
    event.listeners = append(event.listeners, fn)
  }
  return emitter
}

// addThrottle appends the throttle `t` to the event's throttles slice.
func (emitter *Emitter) addThrottle(e string, t *throttle) *Emitter {
  emitter.mutex.Lock()
  defer emitter.mutex.Unlock()
  if nil == t.fn {
    return emitter
  }
  _, ok := emitter.events[e]
  if !ok {
    emitter.events[e] = createEvent()
  }
  event := emitter.events[e]
  if emitter.maxListeners >= len(event.listeners)+len(event.throttles)+1 {
    event.throttles = append(event.throttles, t)
  }
  return emitter
}

// On is an alias method for AddListener.
func (emitter *Emitter) On(e string, fn listener) *Emitter {
  return emitter.AddListener(e, fn)
}

func (a listener) Equals(b listener) bool {
  return reflect.ValueOf(a) == reflect.ValueOf(b)
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
      if l.Equals(fn) {
        ev.listeners = append(ev.listeners[:i], ev.listeners[i+1:]...)
      }
    }
    for i, t := range ev.throttles {
      if t.fn.Equals(fn) {
        ev.throttles = append(ev.throttles[:i], ev.throttles[i+1:]...)
      }
    }
  }
  return emitter
}

// Off is an alias method for RemoveListener.
func (emitter *Emitter) Off(e string, fn listener) *Emitter {
  return emitter.RemoveListener(e, fn)
}

// Once adds a listener function `fn` to an event `e` that will run a maximum of
// one time before being removed from it's listener array.
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

func (t *throttle) handler(args ...interface{}) {
  if t.lastCalled.IsZero() {
    t.lastCalled = time.Now()
  } else if time.Now().Sub(t.lastCalled) < t.interval {
    return
  }
  t.lastCalled = time.Now()
  t.fn(args...)
}

// Throttled is an event handler `fn` for event `e` that can only be called once
// within a given duration `interval`.
func (emitter *Emitter) Throttled(e string, interval time.Duration,
  fn listener) *Emitter {
  return emitter.addThrottle(e, &throttle{time.Time{}, interval, fn})
}

// Emit triggers an event `e`, passing along arguments `args` to each of the
// event's listeners.  Each listener function is ran as a goroutine.
func (emitter *Emitter) Emit(e string, args ...interface{}) *Emitter {
  if _, ok := emitter.events[e]; !ok {
    return emitter
  }
  ev := emitter.events[e]
  var wg sync.WaitGroup
  wg.Add(len(ev.listeners) + len(ev.throttles))
  for _, fn := range ev.listeners {
    go func(fn listener) {
      defer wg.Done()
      fn(args...)
    }(fn)
  }
  for _, t := range ev.throttles {
    go func(t *throttle) {
      defer wg.Done()
      t.handler(args...)
    }(t)
  }
  wg.Wait()
  return emitter
}

// EmitSync trigger an event `e`, passing arguments `args` to each of the
// event's listeners.  Each listener function is ran synchronously in the order
// that each listener was added to the event.
func (emitter *Emitter) EmitSync(e string, args ...interface{}) *Emitter {
  if _, ok := emitter.events[e]; !ok {
    return emitter
  }
  var listeners = emitter.events[e].listeners
  var throttles = emitter.events[e].throttles
  for _, fn := range listeners {
    fn(args...)
  }
  for _, t := range throttles {
    t.handler(args...)
  }
  return emitter
}

// SetMaxListeners sets an emitters maximum listeners per event.
// SetMaxListeners will not discard existing listeners beyond the new limit.
func (emitter *Emitter) SetMaxListeners(max int) *Emitter {
  emitter.maxListeners = max
  return emitter
}

// Returns a pointer to an Emitter struct.
func NewEmitter() *Emitter {
  return &Emitter{
    make(map[string]*event),
    DEFAULT_MAX_LISTENERS,
    sync.Mutex{},
  }
}
