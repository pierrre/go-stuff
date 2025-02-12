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

func BenchmarkFibonacci(b *testing.B) {
	for b.Loop() {
		Fibonacci(50)
	}
}

func BenchmarkFibonacciString(b *testing.B) {
	for b.Loop() {
		FibonacciString(50)
	}
}
