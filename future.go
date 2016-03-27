package future

import (
	"errors"
	"time"
)

type Value interface{}

type Future interface {
	Get() (Value, error)
	GetWithTimeout(timeout time.Duration) (Value, error)
}

type futureResult struct {
	result chan *result
}

type result struct {
	value Value
	err   error
}

var ErrTimeout = errors.New("Timed out")

func NewFuture(Func func() (Value, error)) Future {
	return newFutureResult(Func)
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
