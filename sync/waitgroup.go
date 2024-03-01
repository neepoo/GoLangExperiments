package sync

import stdsync "sync"

type WaitGroup struct {
	groupSize int
	cond      *stdsync.Cond
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		cond: stdsync.NewCond(new(stdsync.Mutex)),
	}
}

func (wg *WaitGroup) Add(n int) {
	wg.cond.L.Lock()
	wg.groupSize += n
	wg.cond.L.Unlock()
}

func (wg *WaitGroup) Done() {
	wg.cond.L.Lock()
	wg.groupSize--
	if wg.groupSize == 0 {
		// like std WaitGroup
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}

func (wg *WaitGroup) Wait() {
	wg.cond.L.Lock()
	for wg.groupSize > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}
