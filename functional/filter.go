package functional

import "golang.org/x/exp/constraints"

func filter[I any](is []I, predicate func(I) bool) []I {
	var out []I
	for _, i := range is {
		if predicate(i) {
			out = append(out, i)
		}
	}
	return out
}

// createLargerThanPredicate 泛型版本，接受一个阈值并返回一个判断函数
func createLargerThanPredicate[T constraints.Ordered](threshold T) func(T) bool {
	return func(i T) bool {
		return i > threshold
	}
}
