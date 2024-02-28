package main

import (
	"fmt"
	"sync"
)

/*
	Go’s runtime checks to see which goroutine it should execute next, and if it finds that all of them are blocked while waiting

for a resource (such as a mutex), it will throw a fatal error. Unfortunately, this means
that it will only catch a deadlock if all the goroutines are blocked.
*/

func lockBoth(l1, l2 *sync.Mutex, wg *sync.WaitGroup) {
	for i := 0; i < 10000; i++ {
		l1.Lock()
		l2.Lock()
		l1.Unlock()
		l2.Unlock()
	}
	wg.Done()
}

func main() {
	l1, l2 := new(sync.Mutex), new(sync.Mutex)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go lockBoth(l1, l2, wg)
	go lockBoth(l2, l1, wg)
	wg.Wait()
	// 如果把wg.Wait()换成下面的四行，Go's runtime 不会检测到死锁，因为main goroutine没有被block而是在等待Sleep
	//go func() {
	//	wg.Wait()
	//	fmt.Println("Done waiting on waitgroup")
	//}()
	//time.Sleep(30 * time.Second)
	fmt.Println("Done")
}
