# future
golang implementation of futures/promises

[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/sentientmonkey/future)

# Install

```
go get github.com/sentientmonkey/future
```

# Usage

```go
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/sentientmonkey/future"
)

func main() {
	f := future.NewFuture(func() (future.Value, error) {
		return http.Get("http://golang.org/")
	})

	result, err := f.Get()
	if err != nil {
		fmt.Printf("Got error: %s\n", err)
		return
	}

	response := result.(*http.Response)
	defer response.Body.Close()
	fmt.Printf("Got result: %d\n", response.StatusCode)
	// Output: Got result: 200
}
```

# License
MIT licensed. See the LICENSE file for details.
