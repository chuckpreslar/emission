// The MIT License (MIT)

// Copyright (c) 2013 Chuck Preslar

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package emission provides an event emitter.
package emission

import (
  "fmt"
  "sync"
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
}

// AddListener adds the listener function (fn) to the Emitter's event (e)
// listener array.
func (emitter *Emitter) AddListener(e string, fn func(...interface{})) *Emitter {
  if nil == fn {
    return emitter
  }
  fn = listener(fn)
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

// On is an alias method for Emmiter#AddListener.
func (emitter *Emitter) On(e string, fn func(...interface{})) *Emitter {
  return emitter.AddListener(e, fn)
}

// RemoveListener loops through an Emitter's events and listeners, comparing
// the string value of the given listener function (fn) since go
// does not allow you to compare functions.  If a match is found,
// it is removed from the event's listeners array.
func (emitter *Emitter) RemoveListener(fn func(...interface{})) *Emitter {
  for _, e := range emitter.events {
    for i, l := range e.listeners {
      if fmt.Sprintf("%v", l) == fmt.Sprintf("%v", fn) {
        e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
      }
    }
  }
  return emitter
}

// Off is an alias method for #RemoveListener.
func (emitter *Emitter) Off(fn func(...interface{})) *Emitter {
  return emitter.RemoveListener(fn)
}

// Once adds a listener function (fn) to an event (e) that will run a maximum of one time
// before being removed from it's listener array.
func (emitter *Emitter) Once(e string, fn func(...interface{})) *Emitter {
  if nil == fn {
    return emitter
  }
  var run listener
  run = func(args ...interface{}) {
    fn(args...)
    emitter.RemoveListener(run)
  }
  emitter.AddListener(e, run)
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

// Returns a pointer to an Emitter struct.
func NewEmitter() *Emitter {
  return &Emitter{
    make(map[string]*event),
    DEFAULT_MAX_LISTENERS,
  }
}
