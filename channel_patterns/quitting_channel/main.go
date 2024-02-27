package main

import "fmt"

func printNumbers(numbers <-chan int, quit chan struct{}) {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-numbers)
		}
		close(quit)
	}()
}

func main() {
	numbers := make(chan int)
	quit := make(chan struct{})
	printNumbers(numbers, quit)
	next := 0
	for i := 1; ; i++ {
		next += i
		select {
		case numbers <- next:
		case <-quit:
			fmt.Println("quit program")
			return
		}
	}

}
