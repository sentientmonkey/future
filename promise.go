package future

// Promises are like Futures, but provide the additional functionality of being
// chained together.
type Promise interface {
	// Block on Promise awaiting result chain
	Get() (Value, error)
	// Then adds an additional async function and return a new Promise.
	// The result from the previous promise is passed along as value.
	// Func is not invoked if a previous Promise returns an error.
	Then(Func func(value Value) (Value, error)) Promise
}

// Create a new Promise. Func is asynchronously called and is resolved
// with a Get() call at an end of a Promise chain.
func NewPromise(Func func() (Value, error)) Promise {
	return newFutureResult(Func)
}

func (p *futureResult) Then(Func func(value Value) (Value, error)) Promise {
	next := &futureResult{
		result: make(chan *result),
	}

	go func() {
		res := <-p.result
		if res.err != nil {
			next.result <- res
			return
		}

		value, err := Func(res.value)
		next.result <- &result{value, err}
	}()
	return next
}
