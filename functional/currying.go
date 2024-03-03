package functional

func threeSum(a, b, c int) int {
	return a + b + c
}

func threeSumCurried(a int) func(b int) func(c int) int {
	return func(b int) func(c int) int {
		return func(c int) int {
			return a + b + c
		}
	}
}
