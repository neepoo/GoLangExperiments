package functional

import (
	"fmt"
	"testing"
)

func Example_createLargerThanPredicate() {
	out := filter([]int{1, 2, 3, 4, 5, 6}, createLargerThanPredicate(3))
	fmt.Println(out)

	// Output:
	// [4 5 6]
}

func Example_partial() {
	bucky := maleHavaneseSpawner("bucky")
	rocky := maleHavaneseSpawner("rocky")
	tipsy := femalePoodleSpawner("tipsy")
	fmt.Printf("%v\n", bucky)
	fmt.Printf("%v\n", rocky)
	fmt.Printf("%v\n", tipsy)

	// Output:
	// {bucky 1 0}
	// {rocky 1 0}
	// {tipsy 3 1}
}

func Test_addThree(t *testing.T) {
	got := threeSumCurried(1)(2)(3)
	if threeSum(1, 2, 3) != got {
		t.Fatalf("got: %d, expect 6", got)
	}
}
