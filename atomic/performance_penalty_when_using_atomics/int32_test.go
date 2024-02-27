package performance_penalty_when_using_atomics

import (
	"sync/atomic"
	"testing"
)

var total = int64(0)

/*
RESULT:
goos: windows
goarch: amd64
pkg: github.com/neepoo/GoLangExperiments/atomic/performance_penalty_when_using_atomics
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkNormal
BenchmarkNormal-24      849596016                1.429 ns/op
BenchmarkAtomic
BenchmarkAtomic-24      342733562                3.507 ns/op
PASS
*/
func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		total++
	}
}

func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&total, 1)
	}
}
