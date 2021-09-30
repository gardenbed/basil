[![Go Doc][godoc-image]][godoc-url]

# Factory

This is a helper package for testing purposes.
It offers a set of functions for generating random values for standard built-in data types.
It can also generate random values for user-defined structs through reflection.

## Quick Start

```go
package main

import (
  "fmt"

  "github.com/gardenbed/basil/factory"
)

func main() {
  name := factory.Name()
  email := factory.Email()
  fmt.Printf("%s <%s>\n", name, email)
}
```

```go
package main

import (
  "fmt"
  "log"
  "net/url"
  "time"

  "github.com/gardenbed/basil/factory"
)

{
  object := struct {
    String     string
    Bool       bool
    Int        int
    Uint       uint
    Float64    float64
    Complex128 complex128
    Nested     struct {
      Duration time.Duration
      Time     *time.Time
      URL      *url.URL
    }
  }{}

  if err := factory.Populate(&object, false); err != nil {
    log.Fatalf("populate error: %s", err)
  }

  fmt.Printf("%+v\n", object)
}
```


[godoc-url]: https://pkg.go.dev/github.com/gardenbed/basil/factory
[godoc-image]: https://pkg.go.dev/badge/github.com/gardenbed/basil/factory
