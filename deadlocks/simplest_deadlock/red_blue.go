package main

import (
	"fmt"
	"sync"
	"time"
)

func red(l1, l2 *sync.Mutex) {
	for {
		fmt.Println("Red: Acquiring lock1")
		l1.Lock()
		fmt.Println("Red: Acquiring lock2")
		l2.Lock()
		fmt.Println("Red: Both locks Acquired")
		l1.Unlock()
		l2.Unlock()
		fmt.Println("Red: Locks Released")
	}
}

func blue(l1, l2 *sync.Mutex) {
	for {
		fmt.Println("Blue: Acquiring lock2")
		l2.Lock()
		fmt.Println("Blue: Acquiring lock1")
		l1.Lock()
		fmt.Println("Blue: Both locks Acquired")
		l1.Unlock()
		l2.Unlock()
		fmt.Println("Blue: locks Released")
	}
}

func main() {
	l1, l2 := new(sync.Mutex), new(sync.Mutex)
	go red(l1, l2)
	go blue(l1, l2)
	time.Sleep(20 * time.Second)
	fmt.Println("Done")

}
