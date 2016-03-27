package future

type Promise struct {
	result chan *result
}

func NewPromise(Func func() (Value, error)) *Promise {
	p := &Promise{
		result: make(chan *result),
	}

	go func() {
		value, err := Func()
		p.result <- &result{value, err}
	}()
	return p
}

func (p *Promise) Get() (Value, error) {
	defer close(p.result)
	result := <-p.result
	return result.value, result.err
}

func (p *Promise) Then(Func func(value Value) (Value, error)) *Promise {
	next := &Promise{
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
