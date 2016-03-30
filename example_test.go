package future_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/sentientmonkey/future"
)

func ExampleFuture() {
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

func ExamplePromise() {
	p := future.NewPromise(func() (future.Value, error) {
		return http.Get("http://golang.org/")
	})

	p = p.Then(func(value future.Value) (future.Value, error) {
		response := value.(*http.Response)
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return string(b), nil
	})

	p = p.Then(func(value future.Value) (future.Value, error) {
		body := value.(string)
		r, err := regexp.Compile("<title>(.*)</title>")
		if err != nil {
			return nil, err
		}
		match := r.FindStringSubmatch(body)

		if len(match) < 1 {
			return nil, errors.New("Title not found")
		}

		return match[1], nil
	})

	result, err := p.Get()
	if err != nil {
		fmt.Printf("Got error: %s\n", err)
		return
	}

	s := result.(string)
	fmt.Printf("Got result: %s\n", s)
	// Output: Got result: The Go Programming Language
}
