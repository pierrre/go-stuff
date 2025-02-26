package main

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
)

const testIterations = 100

func TestFibonacci(t *testing.T) {
	v := Fibonacci(testIterations)
	assertauto.Equal(t, v)
}

func TestFibonacciString(t *testing.T) {
	s := FibonacciString(testIterations)
	assertauto.Equal(t, s)
}

func BenchmarkFibonacci(b *testing.B) {
	for b.Loop() {
		Fibonacci(testIterations)
	}
}

func BenchmarkFibonacciString(b *testing.B) {
	for b.Loop() {
		FibonacciString(testIterations)
	}
}
