package functional

import "fmt"

func Example_createLargerThanPredicate() {
	out := filter([]int{1, 2, 3, 4, 5, 6}, createLargerThanPredicate(3))
	fmt.Println(out)

	// Output:
	// [4 5 6]
}
