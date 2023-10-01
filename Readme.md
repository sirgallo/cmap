# CMap


## Installation

In your `Go` project main directory (where the `go.mod` file is located)
```bash
go get github.com/sirgallo/cmap
go mod tidy
```

Make sure to run go mod tidy to install dependencies.


## Usage

```go
package main

import "github.com/sirgallo/cmap"

func main() {
  // initialize c map

  // for 32 bit hash and bitmap
  cMap := cmap.NewCMap[string, uint32]()

  // for 64 bit hash and bitmap
  cMap := cmap.NewCMap[string, uint64]()

  // insert key/val pair
  cMap.Put("hi", "world")

  // retrieve value for key
  val := cMap.Get("hi")

  // delete key/val pair
  cMap.Delete("hi")
}
```

## Tests

```bash
go test -v ./tests
```


## Sources

[CMap](./docs/CMap.md)

[Murmur](./docs/Murmur.md)