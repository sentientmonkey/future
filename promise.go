package future

type Promise interface {
	Get() (Value, error)
	Then(Func func(value Value) (Value, error)) Promise
}

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
