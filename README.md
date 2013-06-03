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


* __NewEmitter__

    Returns a pointer to an instance of an `Emitter` struct.

### Emitter Receiver Methods

* __AddListener__

    Takes arguments of a `string` (event name), and a `func` (event listener) to be called when the event occurs (or is emitted).

* __On__

    Alias method for `#AddListener`

* __RemoveListener__

    Takes a `func` argument, searches for a matching `func` listener within the Emitter's event listeners and removes it.

* __Off__

    Alias method for `#RemoveListener`

* __Once__

    Similar to `#AddListener` or `#On`, `#Once` takes the same parameters but it's listener will only be called a maximum of one time before it is removed.

* __Emit__

    Takes a `string` representing the event to emit, and arguments of type `interface` to pass along to any event listeners.  Each listener function is ran as a go routine.

## License

> The MIT License (MIT)

> Copyright (c) 2013 Chuck Preslar

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.

