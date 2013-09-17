// Package emission provides an event emitter.
package emission

import (
  "errors"
  "fmt"
  "os"
  "reflect"
  "sync"
)

const (
  DefaultMaxListeners = 10
)

var (
  ErrNoneFunction = errors.New("Kind of Value for listener is not Func.")
)

type Emitter struct {
  events       map[interface{}][]reflect.Value
  maxListeners int
}

// AddListener appends the listener argument to the event arguments slice
// in the Emitter's events map. If the number of listeners for an event
// is greater than the Emitter's maximum listeners then a warning is printed.
// If the relect Value of the listener does not have a Kind of Func then
// AddListener panics.
func (emitter *Emitter) AddListener(event, listener interface{}) *Emitter {
  fn := reflect.ValueOf(listener)

  if reflect.Func != fn.Kind() {
    panic(ErrNoneFunction)
  }

  if emitter.maxListeners != -1 && emitter.maxListeners < len(emitter.events[event])+1 {
    fmt.Fprintf(os.Stdout, "Warning: event `%v` has exceeded the maximum "+
      "number of listeners of %d.\n", event, emitter.maxListeners)
  }

  emitter.events[event] = append(emitter.events[event], fn)

  return emitter
}

// On is an alias for AddListener.
func (emitter *Emitter) On(event, listener interface{}) *Emitter {
  return emitter.AddListener(event, listener)
}

// RemoveListener removes the listener argument from the event arguments slice
// in the Emitter's events map.  If the reflect Value of the listener does not
// have a Kind of Func then RemoveListener panics.
func (emitter *Emitter) RemoveListener(event, listener interface{}) *Emitter {
  fn := reflect.ValueOf(listener)

  if reflect.Func != fn.Kind() {
    panic(ErrNoneFunction)
  }

  if events, ok := emitter.events[event]; ok {
    for i, listener := range events {
      if fn == listener {
        // Do not break here to ensure the listener has not been
        // added more than once.
        emitter.events[event] = append(emitter.events[event][:i], emitter.events[event][i+1:]...)
      }
    }
  }

  return emitter
}

// Off is an alias for RemoveListener.
func (emitter *Emitter) Off(event, listener interface{}) *Emitter {
  return emitter.RemoveListener(event, listener)
}

// Once generates a new function which invokes the supplied listener
// only once before removing itself from the event's listener slice
// in the Emitter's events map. If the reflect Value of the listener
// does not have a Kind of Func then Once panics.
func (emitter *Emitter) Once(event, listener interface{}) *Emitter {
  fn := reflect.ValueOf(listener)

  if reflect.Func != fn.Kind() {
    panic(ErrNoneFunction)
  }

  var run func(...interface{})

  run = func(arguments ...interface{}) {
    var values []reflect.Value

    for i := 0; i < len(arguments); i++ {
      values = append(values, reflect.ValueOf(arguments[i]))
    }

    fn.Call(values)
    emitter.RemoveListener(event, run)
  }

  emitter.AddListener(event, run)
  return emitter
}

// Emit attempts to use the reflect package to Call each listener stored
// in the Emitter's events map with the supplied arguments. Each listener
// is called within its own go routine. The reflect package will panic if
// the agruments supplied do not align the parameters of a listener function.
func (emitter *Emitter) Emit(event interface{}, arguments ...interface{}) *Emitter {
  var (
    listeners []reflect.Value
    ok        bool
  )

  if listeners, ok = emitter.events[event]; !ok {
    // If the Emitter does not include the event in its
    // event map, it has no listeners to Call yet.
    return emitter
  }

  var (
    wg     sync.WaitGroup
    values []reflect.Value
  )

  for i := 0; i < len(arguments); i++ {
    values = append(values, reflect.ValueOf(arguments[i]))
  }

  wg.Add(len(listeners))

  for _, fn := range listeners {
    go func(fn reflect.Value) {
      defer wg.Done()
      fn.Call(values)
    }(fn)
  }

  wg.Wait()
  return emitter
}

// SetMaxListeners sets the maximum number of listeners per
// event for the Emitter. If -1 is passed as the maximum,
// all events may have unlimited listeners. By default, each
// event can have a maximum number of 10 listeners which is
// useful for finding memory leaks.
func (emitter *Emitter) SetMaxListeners(max int) *Emitter {
  emitter.maxListeners = max
  return emitter
}

// NewEmitter returns a new Emitter object, defaulting the
// number of maximum listeners per event to the DefaultMaxListeners
// constant and initializing its events map.
func NewEmitter() (emitter *Emitter) {
  emitter = new(Emitter)
  emitter.events = make(map[interface{}][]reflect.Value)
  emitter.maxListeners = DefaultMaxListeners
  return
}
