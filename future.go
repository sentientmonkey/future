package future

import (
	"errors"
	"time"
)

type Value interface{}

type Future struct {
	result chan *result
}

type result struct {
	value Value
	err   error
}

var ErrTimeout = errors.New("Timed out")

func NewFuture(Func func() (Value, error)) *Future {
	f := &Future{
		result: make(chan *result),
	}
	go func() {
		defer close(f.result)
		value, err := Func()
		f.result <- &result{value, err}
	}()
	return f
}

func (f *Future) Get() (Value, error) {
	result := <-f.result
	return result.value, result.err
}

func (f *Future) GetWithTimeout(timeout time.Duration) (Value, error) {
	select {
	case result := <-f.result:
		return result.value, result.err
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}
