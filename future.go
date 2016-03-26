package future

import (
	"errors"
	"time"
)

type Future struct {
	result chan *result
}

type result struct {
	value interface{}
	err   error
}

var ErrTimeout = errors.New("Timed out")

func NewFuture(Func func() (interface{}, error)) *Future {
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

func (f *Future) Get() (interface{}, error) {
	ret := <-f.result
	return ret.value, ret.err
}

func (f *Future) GetWithTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case ret := <-f.result:
		return ret.value, ret.err
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}
