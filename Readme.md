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

import "github.com/sirgallo/ads/pkg/map"

func main() {
  // initialize c map

  // for 32 bit hash and bitmap
  cMap := cmap.NewLFMap[string, uint32]()

  // for 64 bit hash and bitmap
  cMap := cmap.NewLFMap[string, uint64]()

  // insert key/val pair
  cMap.Insert("hi", "world")

  // retrieve value for key
  val := cMap.Retrieve("hi")

  // delete key/val pair
  cMap.Delete("hi")
}
```

to test:
```bash
go test -v ./tests
```


## Sources

[CMap](./docs/CMap.md)

[Murmur](./docs/Murmur.md)