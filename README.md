## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get -u github.com/chuckpreslar/emission

## Usage

If you've ever used an event emitter before, using Emission should be very familiar.

```go
package main

import (
  "fmt"
  e "github.com/chuckpreslar/emission"
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
  emitter.On("test", b).
    On("hello", func(args ...interface{}) {
    for _, value := range args {
      fmt.Printf("Hello %v!\n", value)
    }
  })
  emitter.Emit("test", 1, 2, 3, 4).
    Emit("hello", "world").
    Emit("hello", "Chuck")

  /**
   * Hello from `a`! [1 2 3 4]
   * Hello from `b`! [1 2 3 4]
   * Hello world!
   * Hello Chuck!
   */

}

```

## Documentation

View godoc's or visit [godoc.org](http://godoc.org/github.com/chuckpreslar/emission).

    $ godoc emission
    
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

