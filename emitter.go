/**
 *	The MIT License (MIT)
 *
 *	Copyright (c) 2013 Chuck Preslar
 *
 *	Permission is hereby granted, free of charge, to any person obtaining a copy
 *	of this software and associated documentation files (the "Software"), to deal
 *	in the Software without restriction, including without limitation the rights
 *	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *	copies of the Software, and to permit persons to whom the Software is
 *	furnished to do so, subject to the following conditions:
 *
 *	The above copyright notice and this permission notice shall be included in
 *	all copies or substantial portions of the Software.
 *
 *	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *	THE SOFTWARE.
 */

package emission

import (
  "fmt"
  "sync"
)

const DEFAULT_MAX_LISTENERS = 10

type listener func(...interface{})

type Event struct {
  listeners []listener
}

type Emitter struct {
  events       map[string]*Event
  maxListeners int
}

/**
 * Adds the listener function `fn` to the `Emitter`'s event `e`
 * listener array.
 *
 * @receiver *Emitter
 */

func (emitter *Emitter) AddListener(e string, fn func(...interface{})) {
  if nil == fn {
    return
  }
  fn = listener(fn)
  _, ok := emitter.events[e]
  if !ok {
    emitter.events[e] = &Event{[]listener{}}
  }
  emitter.events[e].listeners = append(emitter.events[e].listeners, fn)
}

/**
 * Alias method for `#AddListener`.
 *
 * @receiver *Emitter
 */

func (emitter *Emitter) On(e string, fn func(...interface{})) {
  emitter.AddListener(e, fn)
}

/**
 * Loops through an `Emitter`'s events and listeners, comparing
 * the string value of the given listener function `fn` since go
 * does not allow you to compare functions.  If a match is found,
 * it is removed from the event's listeners array.
 *
 * @receiver *Emitter
 */

func (emitter *Emitter) RemoveListener(fn func(...interface{})) {
  for _, x := range emitter.events {
    for i, y := range x.listeners {
      if fmt.Sprintf("%v", y) == fmt.Sprintf("%v", fn) {
        x.listeners = append(x.listeners[:i], x.listeners[i+1:]...)
      }
    }
  }
}

/**
 * Alias method for `#RemoveListener`.
 *
 * @receiver *Emitter
 */

func (emitter *Emitter) Off(fn func(...interface{})) {
  emitter.RemoveListener(fn)
}

/**
 * Adds a listener function to an event that will run a maximum of one time
 * before being removed from it's listener array.
 *
 * @receiver *Emitter
 */

func (emitter *Emitter) Once(e string, fn func(...interface{})) {
  if nil == fn {
    return
  }
  emitter.AddListener(e, fn)
  emitter.AddListener(e, func(args ...interface{}) {
    emitter.RemoveListener(fn)
  })
}

/**
 * Triggers an event `e`, passing along `args` to each of the event's
 * listeners.  Each listener function is ran as a go routine.
 *
 * @receiver *Emitter
 */

func (emitter *Emitter) Emit(e string, args ...interface{}) {
  if _, ok := emitter.events[e]; !ok {
    return
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
}

/**
 * Returns a pointer to an `Emitter` struct.
 *
 * @returns *Emitter
 */

func NewEmitter() *Emitter {
  return &Emitter{
    make(map[string]*Event),
    DEFAULT_MAX_LISTENERS,
  }
}
