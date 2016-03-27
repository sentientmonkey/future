package future

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromise(t *testing.T) {
	p := NewPromise(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})
	value, err := p.Get()
	assert.Equal(t, 42, value)
	assert.NoError(t, err)
}

func TestPromiseThen(t *testing.T) {
	p := NewPromise(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})
	value, err := p.Then(func(value Value) (Value, error) {
		time.Sleep(1 * time.Second)
		return value.(int) + 3, nil
	}).Get()
	assert.Equal(t, 45, value)
	assert.NoError(t, err)
}

func TestPromiseErrorThen(t *testing.T) {
	p := NewPromise(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, errors.New("error!")
	})
	value, err := p.Then(func(value Value) (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return value.(int) + 3, nil
	}).Get()
	assert.Nil(t, value)
	assert.Error(t, err)
}

func TestPromiseChain(t *testing.T) {
	value, err := NewPromise(func() (Value, error) {
		time.Sleep(10 * time.Millisecond)
		return 20, nil
	}).Then(func(value Value) (Value, error) {
		time.Sleep(10 * time.Millisecond)
		return value.(int) - 10, nil
	}).Then(func(value Value) (Value, error) {
		time.Sleep(10 * time.Millisecond)
		return value.(int) * 3, nil
	}).Then(func(value Value) (Value, error) {
		time.Sleep(10 * time.Millisecond)
		return value.(int) / 5, nil
	}).Get()

	assert.Equal(t, 6, value)
	assert.NoError(t, err)
}

func TestPromiseChainDelay(t *testing.T) {
	start := time.Now()

	p1 := NewPromise(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	p2 := p1.Then(func(value Value) (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return value.(int) + 3, nil
	})

	assert.InDelta(t, 0.0, time.Since(start).Seconds(), 0.01)

	value, err := p2.Get()

	assert.Equal(t, 45, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.2, time.Since(start).Seconds(), 0.05)
}
