package pipelining

/*
The steps in a simple pipeline all follow the same pattern:
accept input from an input channel of type X, process X,
and produce the result Y on an output channel of type Y.
*/

func AddOnPipe[X, Y any](quit <-chan struct{}, f func(X) Y, in <-chan X) chan Y {
	output := make(chan Y)
	go func() {
		defer close(output)
		for {
			select {
			case <-quit:
				return
			case input := <-in:
				output <- f(input)
			}
		}
	}()
	return output
}
