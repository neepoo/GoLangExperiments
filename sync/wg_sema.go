package sync

// this file contains the simply WaitGroup implementation by semaphore

type WaitGrp struct {
	sema *Semaphore
}

// NewWaitGrp init wg, The main difference
// between Go’s bundled waitgroup and our implementation is that we need to specify
// the size of the waitgroup at the start before we use it. In the waitgroup in Go’s sync
// package, we can increase the size of the group at any point—even when we have
// goroutines waiting on the work to be completed
func NewWaitGrp(size int) *WaitGrp {
	return &WaitGrp{sema: NewSemaphore(1 - size)}
}

func (wg *WaitGrp) Wait() {
	wg.sema.Acquire()
}

func (wg *WaitGrp) Done() {
	wg.sema.Release()
}
