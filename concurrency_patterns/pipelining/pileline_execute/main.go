package main

import (
	"fmt"
	"time"

	"github.com/neepoo/GoLangExperiments/concurrency_patterns/pipelining"
)

func main() {
	now := time.Now()
	input := make(chan int)
	quit := make(chan struct{})
	defer close(quit)
	output := pipelining.AddOnPipe(quit, pipelining.Box,
		pipelining.AddOnPipe(quit, pipelining.AddToppings,
			pipelining.AddOnPipe(quit, pipelining.Bake,
				pipelining.AddOnPipe(quit, pipelining.Mixture,
					pipelining.AddOnPipe(quit, pipelining.PrepareTray, input)))))
	go func() {
		for i := 0; i < 10; i++ {
			input <- i
		}
	}()
	for i := 0; i < 10; i++ {
		fmt.Println(<-output, "received")
	}
	fmt.Println("elapsed:", time.Since(now).String())
}
