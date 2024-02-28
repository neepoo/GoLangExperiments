package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// We’ll start by using an atomic
//variable as an indicator showing whether the mutex is locked. We can then use the
//CompareAndSwap() function to check and update the value of the indicator whenever
//we need to lock the mutex

// If the indicator is showing as free, CompareAndSwap(unlocked, locked) will succeed, and
//the indicator will be updated to locked. If the indicator is showing as locked, the
//CompareAndSwap(unlocked, locked) operation will fail, returning false. At this point, we
//can keep retrying until the indicator changes value and becomes unlocked. This type
//of mutex is called a spin lock.

const (
	// unlocked is the unlocked state of the mutex
	free int32 = 0
	// locked is the locked state of the mutex
	locked int32 = 1
)

// SpinLock is a type of lock in which an execution will go into a
// loop to try to get hold of a lock repeatedly until the lock becomes available.
type SpinLock int32

func (s *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32((*int32)(s), free, locked) {
		// Call the Go scheduler to give execution time to other goroutines
		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32((*int32)(s), free)
}

func NewSpinLock() sync.Locker {
	return new(SpinLock)
}
