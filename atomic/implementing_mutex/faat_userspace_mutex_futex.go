package main

import "sync/atomic"

//Futex is short for fast userspace mutex. However, this definition is misleading, as futexes
//are not mutexes at all. A futex is a wait queue primitive that we can access from user space.

//It gives us the ability to suspend and awaken an execution on a specific address.
//Futexes come in handy when we need to implement efficient concurrency primitives
//such as mutexes, semaphores, and condition variables.

//When we call futex_wait(addr, value), we specify a memory address and a value. If
//the value at the memory address is equal to the specified parameter value, the execu-
//tion of the caller is suspended and placed at the back of a queue. The queue parks all
//the executions that have called futex_wait() on the same address value. The operat-
//ing system models a different queue for each memory address value.
//When we call futex_wait(addr, value) and the value of the memory address is dif-
//ferent from the parameter value, the function returns immediately, and the execution
//continues.

// The futex_wake(addr, count) wakes up suspended executions (threads and processes)
//that are waiting on the address specified. The operating system resumes a total of
//count executions, and it picks up the executions from the front of the queue. If the
//count parameter is 0, all the suspended executions are resumed.

// When there is resource contention, the executions will not loop
// needlessly, wasting CPU cycles. Instead, they will wait on a futex. They will be queued
// until the lock becomes available again.
type FutexLock int32

func (f *FutexLock) Lock() {
	for !atomic.CompareAndSwapInt32((*int32)(f), free, locked) {
		// if the lock is not available, waits, but only if the variables has a state of locked
		futex_wait((*int32)(f), locked)
	}
}

func (f *FutexLock) Unlock() {
	atomic.StoreInt32((*int32)(f), free)
	// wakes up 1 execution
	futex_wakeup((*int32)(f), 1)
}
