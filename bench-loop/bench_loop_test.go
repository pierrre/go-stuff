package benchloop

import (
	"slices"
	"testing"
)

var (
	benchRes any
	s        = slices.Repeat([]int{123}, 100000)
)

func BenchmarkNew(b *testing.B) {
	it := slices.Values(s)
	for b.Loop() {
		for v := range it {
			_ = v
		}
	}
}

func BenchmarkOld(b *testing.B) {
	it := slices.Values(s)
	b.ResetTimer()
	var res int
	for range b.N {
		for v := range it {
			res = v
		}
	}
	benchRes = res
}
