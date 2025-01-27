package main

import (
	"math/big"
	"testing"

	"github.com/pierrre/assert/assertauto"
)

const testIterations = 100

func TestFibonacci(t *testing.T) {
	v := Fibonacci(testIterations)
	assertauto.DeepEqual(t, v)
}

func TestFibonacciString(t *testing.T) {
	s := FibonacciString(testIterations)
	assertauto.Equal(t, s)
}

var benchRes any

func BenchmarkFibonacci(b *testing.B) {
	var res *big.Int
	for range b.N {
		res = Fibonacci(testIterations)
	}
	benchRes = res
}

func BenchmarkFibonacciString(b *testing.B) {
	var res string
	for range b.N {
		res = FibonacciString(testIterations)
	}
	benchRes = res
}
