## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get -u github.com/chuckpreslar/emission

## Usage

If you've ever used an event emitter before, using Emission should be very familiar.

```go
package main

import(
  e "github.com/chuckpreslar/emission"
  "fmt"
)

func main() {
  emitter := e.NewEmitter()
  
  a := func(args ...interface{}) {
    fmt.Println("Hello from `a`!", args)
  }
  
  b := func(args ...interface{}) {
    fmt.Println("Hello from `b`!", args)
  }
  
  emitter.On("test", a) // or emitter.AddListener("test", a)
  emitter.On("test", b) //  ...
  emitter.Emit("test", 1, 2, 3, 4)
  
  /**
   * Hello from `a`! [1 2 3 4]
   * Hello from `b`! [1 2 3 4]
   */
}
```

## Methods


#### NewEmitter

Returns a pointer to an instance of an `Emitter` struct.

### Emitter Receiver Methods

#### AddListener

Takes arguments of a `string` (event name), and a `func` (event listener) to be called when the event occurs (or is emitted).

#### On

Alias method for `#AddListener`

#### RemoveListener

Takes a `func` argument, searches for a matching `func` listener within the Emitter's event listeners and removes it.

#### Off

Alias method for `#RemoveListener`

#### Once

Similar to `#AddListener` or `#On`, `#Once` takes the same parameters but it's listener will only be called a maximum of one time before it is removed.

#### Emit

Takes a `string` representing the event to emit, and arguments of type `interface` to pass along to any event listeners.  Each listener function is ran as a go routine.
