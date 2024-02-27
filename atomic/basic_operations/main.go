package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func stringy(money *int32) {
	for i := 0; i < 1000_000; i++ {
		atomic.AddInt32(money, 10)
	}
	fmt.Println("Sringy Done")
}

func spendy(money *int32) {
	for i := 0; i < 1000_000; i++ {
		atomic.AddInt32(money, -10)
	}
	fmt.Println("Spendy Done")
}

func main() {
	wg := new(sync.WaitGroup)
	money := new(int32)
	wg.Add(2)
	go func() {
		stringy(money)
		wg.Done()
	}()
	go func() {
		spendy(money)
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("Money in account: ", atomic.LoadInt32(money))
}
