package sync

import "sync"

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
	cond         *sync.Cond
}

func NewRWMutex() *RWMutex {
	return &RWMutex{cond: sync.NewCond(new(sync.Mutex))}
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
