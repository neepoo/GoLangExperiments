package main

import (
	"sort"
	"sync"
)

type Flight struct {
	Origin, Dest string
	SeatsLeft    int
	Locker       sync.Locker
}

//returning true
//only if all flights on the input slice contain enough seats for the booking request. The
//implementation starts by sorting the input list of flights in alphabetical order using
//the origin and destination. This ordering is done to avoid deadlocks

func Book(flights []*Flight, seatsToBook int) bool {
	bookable := true
	sort.Slice(flights, func(a, b int) bool {
		flightA := flights[a].Origin + flights[a].Dest
		flightB := flights[b].Origin + flights[b].Dest
		return flightA < flightB
	})
	for _, f := range flights {
		f.Locker.Lock()
	}

	for i := 0; i < len(flights) && bookable; i++ {
		if flights[i].SeatsLeft < seatsToBook {
			bookable = false
		}
	}
	for i := 0; i < len(flights) && bookable; i++ {
		flights[i].SeatsLeft -= seatsToBook
	}
	for _, f := range flights {
		f.Locker.Unlock()
	}
	return bookable
}

// Comparing and swapping
//The CompareAndSwap() function can be used to check and set a flag indicating that a
//resource is locked. This function works by accepting a value pointer and old and new
//parameters. If the old parameter is equal to the value stored at the pointer, the value is
//updated to match that of the new parameter. This operation (like all operations in the
//atomic package) is atomic and thus cannot be interrupted by another execution.

// The value of the variable is what we expect, equal to the old
// parameter. When this happens, the value is updated to that of the new parameter, and
// the function returns true.

// when we
//call the function on a value not equal to the old parameter. In this case, the update is
//not applied, and the function returns false.
