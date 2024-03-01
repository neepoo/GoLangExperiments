package sync

import stdsync "sync"

// We could block new readers from acquiring
//the read lock as soon as a writer calls the WriteLock() function

// To design a write-preferred lock, we need a few properties:
// 1. Readers’ counter—Initially set to 0, this tells us how many reader goroutines are actively accessing the shared resources.
// 2. Writers’ waiting counter—Initially set to 0, this tells us how many writer goroutines
//are suspended waiting to access the shared resource.
// 3. Writer active indicator—Initially set to false, this flag tells us if the resource is currently being updated by a writer goroutine.
// 4. Condition variable with mutex—This allows us to set various conditions on the
//preceding properties, suspending execution when the conditions aren’t met.

type RWMutex struct {
	// Stores the number of readers currently holding the read lock
	readersCounter int
	// Store the number of writers currently waiting
	writersWaiting int
	// Indicates if a writer is holding the write lock
	writerActive bool
	cond         *stdsync.Cond
}

func NewRWMutex() *RWMutex {
	return &RWMutex{cond: stdsync.NewCond(new(stdsync.Mutex))}
}

func (rw *RWMutex) ReadLock() {
	rw.cond.L.Lock()
	// waits on condition variable while writers are waiting or active
	for rw.writersWaiting > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	rw.readersCounter++
	rw.cond.L.Unlock()
}

func (rw *RWMutex) WriteLock() {
	rw.cond.L.Lock()
	rw.writersWaiting++
	// wait as long as reader or a writer are active
	for rw.readersCounter > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	// once the wait is over, decrements the writers waiting counter
	rw.writersWaiting--
	rw.writerActive = true
	rw.cond.L.Unlock()
}

func (rw *RWMutex) ReadUnlock() {
	rw.cond.L.Lock()
	rw.readersCounter--
	//  Since there can only ever be one writer active
	//at any point in time, we can send a broadcast every time we unlock.
	// This will wake up
	//any writers or readers that are currently waiting on the condition variable
	//  If there are
	//both readers and writers waiting, a writer will be preferred since the readers will go
	//back into suspension when the writers’ waiting counter is above
	if rw.readersCounter == 0 {
		rw.cond.Broadcast()
	}
	rw.cond.L.Unlock()
}

func (rw *RWMutex) WriteUnlock() {
	rw.cond.L.Lock()
	rw.writerActive = false
	rw.cond.Broadcast()
	rw.cond.L.Unlock()
}
