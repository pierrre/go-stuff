package main

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
)

const testIterations = 100

func TestFibonacci(t *testing.T) {
	v := fibonacci(testIterations)
	assertauto.Equal(t, v)
}

func TestFibonacciString(t *testing.T) {
	s := fibonacciString(testIterations)
	assertauto.Equal(t, s)
}

func BenchmarkFibonacci(b *testing.B) {
	for b.Loop() {
		fibonacci(testIterations)
	}
}

func BenchmarkFibonacciString(b *testing.B) {
	for b.Loop() {
		fibonacciString(testIterations)
	}
}
