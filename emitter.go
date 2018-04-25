// Package emission provides an event emitter.
package emission

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

// Default number of maximum listeners for an event.
const DefaultMaxListeners = 10

// Error presented when an invalid argument is provided as a listener function
var ErrNoneFunction = errors.New("Kind of Value for listener is not Func.")

// RecoveryListener ...
type RecoveryListener func(interface{}, interface{}, error)

// ListenerHandle is an opaque reference to a previously added listener. You need the handle to remove the listener.
type ListenerHandle uint32

type listenerRecord struct {
	fn			reflect.Value
	handle		ListenerHandle
	isOnce		bool
}

// Emitter ...
type Emitter struct {
	// Mutex to prevent race conditions within the Emitter.
	*sync.Mutex
	// Unique counter to allocate handles
	nextHandle ListenerHandle
	// Map of event to a slice of listener function's reflect Values.
	events map[interface{}][]listenerRecord
	// Optional RecoveryListener to call when a panic occurs.
	recoverer RecoveryListener
	// Maximum listeners for debugging potential memory leaks.
	maxListeners int
}

// AddListener appends the listener argument to the event arguments slice
// in the Emitter's events map. If the number of listeners for an event
// is greater than the Emitter's maximum listeners then a warning is printed.
// If the reflect Value of the listener does not have a Kind of Func then
// AddListener panics. If a RecoveryListener has been set then it is called
// recovering from the panic.
func (emitter *Emitter) AddListener(event, listener interface{}) ListenerHandle {
	return emitter.addListener(event,listener,false)
}

func (emitter *Emitter) addListener(event, listener interface{}, isOnce bool) ListenerHandle {
	emitter.Lock()
	defer emitter.Unlock()

	fn := reflect.ValueOf(listener)

	if reflect.Func != fn.Kind() {
		if nil == emitter.recoverer {
			panic(ErrNoneFunction)
		} else {
			emitter.recoverer(event, listener, ErrNoneFunction)
		}
	}

	if emitter.maxListeners != -1 && emitter.maxListeners < len(emitter.events[event])+1 {
		fmt.Fprintf(os.Stdout, "Warning: event `%v` has exceeded the maximum "+
			"number of listeners of %d.\n", event, emitter.maxListeners)
	}

	emitter.nextHandle = emitter.nextHandle + 1
	handle := emitter.nextHandle

	emitter.events[event] = append(emitter.events[event], listenerRecord{fn,handle, isOnce})

	return handle
}

// On is an alias for AddListener.
func (emitter *Emitter) On(event, listener interface{}) ListenerHandle {
	return emitter.AddListener(event, listener)
}

// RemoveListener removes the listener argument from the event arguments slice
// in the Emitter's events map.  If the reflect Value of the listener does not
// have a Kind of Func then RemoveListener panics. If a RecoveryListener has
// been set then it is called after recovering from the panic.
func (emitter *Emitter) RemoveListener(event interface{}, listenerHandle ListenerHandle) {
	emitter.Lock()
	defer emitter.Unlock()

	if events, ok := emitter.events[event]; ok {
		l := len(events)
		if l == 0 {
			return
		}

		newEvents := make([]listenerRecord,0,l - 1)

		for _, listenerRec := range events {
			if listenerHandle != listenerRec.handle {
				newEvents = append(newEvents, listenerRec)
			}
		}

		emitter.events[event] = newEvents
	}
}

// Off is an alias for RemoveListener.
func (emitter *Emitter) Off(event, listener ListenerHandle) {
	emitter.RemoveListener(event, listener)
}

// Once generates a new function which invokes the supplied listener
// only once before removing itself from the event's listener slice
// in the Emitter's events map. If the reflect Value of the listener
// does not have a Kind of Func then Once panics. If a RecoveryListener
// has been set then it is called after recovering from the panic.
func (emitter *Emitter) Once(event, listener interface{}) ListenerHandle {
	return emitter.addListener(event,listener,true)
}

// Emit attempts to use the reflect package to Call each listener stored
// in the Emitter's events map with the supplied arguments. Each listener
// is called within its own go routine. The reflect package will panic if
// the agruments supplied do not align the parameters of a listener function.
// If a RecoveryListener has been set then it is called after recovering from
// the panic.
func (emitter *Emitter) Emit(event interface{}, arguments ...interface{}) *Emitter {
	var (
		listeners []listenerRecord
		ok        bool
	)

	// Lock the mutex when reading from the Emitter's
	// events map.
	emitter.Lock()

	if listeners, ok = emitter.events[event]; !ok {
		// If the Emitter does not include the event in its
		// event map, it has no listeners to Call yet.
		emitter.Unlock()
		return emitter
	}

	// Unlock the mutex immediately following the read
	// instead of deferring so that listeners registered
	// with Once can aquire the mutex for removal.
	emitter.Unlock()

	var wg sync.WaitGroup

	wg.Add(len(listeners))

	for _, listenerRec := range listeners {
		go func(listenerRec listenerRecord) {
			defer wg.Done()

			fn := listenerRec.fn

			// Recover from potential panics, supplying them to a
			// RecoveryListener if one has been set, else allowing
			// the panic to occur.
			if nil != emitter.recoverer {
				defer func() {
					if r := recover(); nil != r {
						err := fmt.Errorf("%v", r)
						emitter.recoverer(event, fn.Interface(), err)
					}
				}()
			}

			var values []reflect.Value

			for i := 0; i < len(arguments); i++ {
				if arguments[i] == nil {
					values = append(values, reflect.New(fn.Type().In(i)).Elem())
				} else {
					values = append(values, reflect.ValueOf(arguments[i]))
				}
			}

			if listenerRec.isOnce {
				emitter.RemoveListener(event,listenerRec.handle)
			}

			fn.Call(values)
		}(listenerRec)
	}

	wg.Wait()
	return emitter
}

// EmitSync attempts to use the reflect package to Call each listener stored
// in the Emitter's events map with the supplied arguments. Each listener
// is called synchronously. The reflect package will panic if
// the agruments supplied do not align the parameters of a listener function.
// If a RecoveryListener has been set then it is called after recovering from
// the panic.
func (emitter *Emitter) EmitSync(event interface{}, arguments ...interface{}) *Emitter {
	var (
		listeners []listenerRecord
		ok        bool
	)

	// Lock the mutex when reading from the Emitter's
	// events map.
	emitter.Lock()

	if listeners, ok = emitter.events[event]; !ok {
		// If the Emitter does not include the event in its
		// event map, it has no listeners to Call yet.
		emitter.Unlock()
		return emitter
	}

	// Unlock the mutex immediately following the read
	// instead of deferring so that listeners registered
	// with Once can aquire the mutex for removal.
	emitter.Unlock()

	for _, listenerRec := range listeners {
		fn := listenerRec.fn

		// Recover from potential panics, supplying them to a
		// RecoveryListener if one has been set, else allowing
		// the panic to occur.
		if nil != emitter.recoverer {
			defer func() {
				if r := recover(); nil != r {
					err := fmt.Errorf("%v", r)
					emitter.recoverer(event, fn.Interface(), err)
				}
			}()
		}

		var values []reflect.Value

		for i := 0; i < len(arguments); i++ {
			if arguments[i] == nil {
				values = append(values, reflect.New(fn.Type().In(i)).Elem())
			} else {
				values = append(values, reflect.ValueOf(arguments[i]))
			}
		}

		if listenerRec.isOnce {
			emitter.RemoveListener(event,listenerRec.handle)
		}

		fn.Call(values)
	}

	return emitter
}

// RecoverWith sets the listener to call when a panic occurs, recovering from
// panics and attempting to keep the application from crashing.
func (emitter *Emitter) RecoverWith(listener RecoveryListener) *Emitter {
	emitter.recoverer = listener
	return emitter
}

// SetMaxListeners sets the maximum number of listeners per
// event for the Emitter. If -1 is passed as the maximum,
// all events may have unlimited listeners. By default, each
// event can have a maximum number of 10 listeners which is
// useful for finding memory leaks.
func (emitter *Emitter) SetMaxListeners(max int) *Emitter {
	emitter.Lock()
	defer emitter.Unlock()

	emitter.maxListeners = max
	return emitter
}

// GetListenerCount gets count of listeners for a given event.
func (emitter *Emitter) GetListenerCount(event interface{}) (count int) {
	emitter.Lock()
	if listeners, ok := emitter.events[event]; ok {
		count = len(listeners)
	}
	emitter.Unlock()
	return
}

// NewEmitter returns a new Emitter object, defaulting the
// number of maximum listeners per event to the DefaultMaxListeners
// constant and initializing its events map.
func NewEmitter() (emitter *Emitter) {
	emitter = new(Emitter)
	emitter.Mutex = new(sync.Mutex)
	emitter.events = make(map[interface{}][]listenerRecord)
	emitter.maxListeners = DefaultMaxListeners
	return
}
