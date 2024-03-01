package sync

import stdsync "sync"

type Semaphore struct {
	// Permits remaining on the semaphore
	permits int
	cond    *stdsync.Cond
}

func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{
		permits: permits,
		cond:    stdsync.NewCond(new(stdsync.Mutex)),
	}
}

// To implement the Acquire() function, we need to call wait() on a condition variable
//whenever the permits are 0 (or less). If there are enough permits, we simply subtract 1
//from the permit count

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()

	// Waits until there is an available permit
	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits++
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}
