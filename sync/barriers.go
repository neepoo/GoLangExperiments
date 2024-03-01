package sync

import stdsync "sync"

// Waitgroups are great for synchronizing after a task has been completed. But what if we
//need to coordinate our goroutines before we start a task? We might also need to align
//different executions at different points in time. Barriers give us the ability to synchro-
//nize groups of goroutines at specific points in our code.

// Let’s look at a simple analogy to help us compare waitgroups and barriers. A pri-
//vate plane will only leave when all the passengers arrive at the departure terminal.
//This represents a barrier. Everyone has to wait until every passenger arrives at this bar-
//rier (the airport terminal). When everyone has finally arrived, the passengers can pro-
//ceed and board the plane.

// IMPLEMENTING

// To start with, we need to know the size of the group of executions that will be using
//this barrier. In the implementation, we’ll call this the barrier size. We can use this size
//to know when enough goroutines are at the barrier.

type Barrier struct {
	size      int
	waitCount int
	cond      *stdsync.Cond
}

func NewBarrier(size int) *Barrier {
	condVar := stdsync.NewCond(new(stdsync.Mutex))
	return &Barrier{size, 0, condVar}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	b.waitCount++
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.cond.L.Unlock()
}
