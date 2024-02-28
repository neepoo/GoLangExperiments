package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	number := int32(17)
	// Change the value of the variable and return true
	result := atomic.CompareAndSwapInt32(&number, 17, 19)
	fmt.Printf("17 <- swap(17, 19): result: %t, value: %d\n", result, number)

	number = int32(23)
	// Compares and failed, leaving the value of the variable unchanged, and return false
	result = atomic.CompareAndSwapInt32(&number, 17, 19)
	fmt.Printf("23 <- swap(17, 19): result: %t, value: %d\n", result, number)

	// OUTPUT
	//17 <- swap(17, 19): result: true, value: 19
	//23 <- swap(17, 19): result: false, value: 23

}
