/*
If we know in advance the full set of exclusive resources
that our concurrent execution will use, we can use ordering to prevent deadlocks.
*/

/*
考虑下面简单的死锁代码，死锁发生是由于`red()`和`blue()`彼此获取mutexes以不同的次序。
 The red()
goroutine is using lock 1 and then lock 2, while blue() is using lock 2 and then lock 1.
If we change the listing so that they use the locks in the same order, as shown in the
following listing, the deadlock won’t occur.

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

*/

// 修复后的版本
// The deadlock doesn’t occur because we never get in a situation where
// both goroutines are holding different locks and requesting the other one.
// In this scenario, when
// they both try to obtain lock 1 at the same time, only one goroutine will succeed. The
// other one will be blocked until both locks are available again.
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
		l1.Lock()
		fmt.Println("Blue: Acquiring lock1")
		l2.Lock()
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
	time.Sleep(3 * time.Second)
	fmt.Println("Done")

}
