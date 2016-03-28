// Package future provides implementations for futures & promises.
//
// Each Future/Promise creates a goroutinue to execute each one (or
// multiple when Promises are chained), and a channel to block on results.
// These are cleaned up when the original Func calls complete. Note that
// panics are not recovered explicitly, but you can recover then in your
// Func blocks.
package future

import (
	"errors"
	"time"
)

// Value type to allow returning arbitrary results.
type Value interface{}

// A Future is a result to an asynchronous call that cal be blocked on
// for a result when needed.
type Future interface {
	// Blocks on Future awaiting result
	Get() (Value, error)
	// Blocks on Future awaiting result, but returns a ErrTimeout if
	// the timeout Duration is hit before result returns.
	// Note that the execution still in Future after timeout.
	GetWithTimeout(timeout time.Duration) (Value, error)
}

// Returned when
var ErrTimeout = errors.New("Timed out")

// Creates a new Future. Func is asynchronously called and it is
// resolved with a Get or GetWithTimeout call on the Future.
func NewFuture(Func func() (Value, error)) Future {
	return newFutureResult(Func)
}

type futureResult struct {
	result chan *result
}

type result struct {
	value Value
	err   error
}

func newFutureResult(Func func() (Value, error)) *futureResult {
	f := &futureResult{
		result: make(chan *result),
	}
	go func() {
		defer close(f.result)
		value, err := Func()
		f.result <- &result{value, err}
	}()
	return f
}

func (f *futureResult) Get() (Value, error) {
	result := <-f.result
	return result.value, result.err
}

func (f *futureResult) GetWithTimeout(timeout time.Duration) (Value, error) {
	select {
	case result := <-f.result:
		return result.value, result.err
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}
