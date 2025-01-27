package main

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
)

func TestFibonacci(t *testing.T) {
	v := Fibonacci(50)
	assertauto.Equal(t, v)
}

func TestFibonacciString(t *testing.T) {
	s := FibonacciString(50)
	assertauto.Equal(t, s)
}

var benchRes any

func BenchmarkFibonacci(b *testing.B) {
	var res int64
	for range b.N {
		res = Fibonacci(50)
	}
	benchRes = res
}

func BenchmarkFibonacciString(b *testing.B) {
	var res string
	for range b.N {
		res = FibonacciString(50)
	}
	benchRes = res
}
