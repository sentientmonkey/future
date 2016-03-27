package future

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFutureError(t *testing.T) {
	f1 := NewFuture(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, errors.New("test error")
	})

	value, err := f1.Get()
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestFutureAsync(t *testing.T) {
	start := time.Now()
	f1 := NewFuture(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	f2 := NewFuture(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return 43, nil
	})
	value, err := f1.Get()
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	value, err = f2.Get()
	assert.Equal(t, 43, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureWithTimeout(t *testing.T) {
	start := time.Now()
	f1 := NewFuture(func() (Value, error) {
		time.Sleep(1 * time.Second)
		return 42, nil
	})

	value, err := f1.GetWithTimeout(100 * time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, ErrTimeout, err)
	assert.Nil(t, value)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureWithTimeoutComplete(t *testing.T) {
	start := time.Now()
	f1 := NewFuture(func() (Value, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	value, err := f1.GetWithTimeout(1 * time.Second)
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}
