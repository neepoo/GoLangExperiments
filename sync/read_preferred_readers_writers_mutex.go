package sync

import "sync"

//To implement our readers–writer mutex, we need a system that, when a goroutine
//calls ReadLock(), blocks any access to the write part while allowing other goroutines to
//still call ReadLock() without blocking. We’ll block the write part by making sure that a
//goroutine calling WriteLock() suspends execution. Only when all the read goroutines
//call ReadUnlock() will we allow another goroutine to unblock from WriteLock().

type ReadPreferredRWMutex struct {
	// count the number of reader goroutines currently in the critical section
	readersCounter int
	// Mutex for synchronizing readers access
	readersLock sync.Mutex
	// Mutex for blocking writers access
	globalLock sync.Mutex
}

func (rw *ReadPreferredRWMutex) ReadLock() {
	rw.readersLock.Lock()
	rw.readersCounter++
	// if a reader goroutine is the first one in,
	// it attempts to lock globalLock
	if rw.readersCounter == 1 {
		rw.globalLock.Lock()
	}
	rw.readersLock.Unlock()
}

func (rw *ReadPreferredRWMutex) WriteLock() {
	rw.globalLock.Lock()
}

func (rw *ReadPreferredRWMutex) ReadUnlock() {
	rw.readersLock.Unlock()
	rw.readersCounter--
	if rw.readersCounter == 0 {
		rw.globalLock.Unlock()
	}
	rw.readersLock.Unlock()
}

func (rw *ReadPreferredRWMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}
