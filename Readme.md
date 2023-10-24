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
  cMap := cmap.NewCMap[uint32]()

  // for 64 bit hash and bitmap
  cMap := cmap.NewCMap[uint64]()

  // insert key/val pair
  cMap.Put([]byte("hi"), []byte("world"))

  // retrieve value for key
  val := cMap.Get([]byte("hi"))

  // delete key/val pair
  cMap.Delete([]byte("hi"))
}
```

## Tests

```bash
go test -v ./tests
```


## godoc

For in depth definitions of types and functions, `godoc` can generate documentation from the formatted function comments. If `godoc` is not installed, it can be installed with the following:
```bash
go install golang.org/x/tools/cmd/godoc
```

To run the `godoc` server and view definitions for the package:
```bash
godoc -http=:6060
```

Then, in your browser, navigate to:
```
http://localhost:6060/pkg/github.com/sirgallo/cmap/
```


## Sources

[CMap](./docs/CMap.md)

[Murmur](./docs/Murmur.md)