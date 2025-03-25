package main

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
)

func TestFibonacci(t *testing.T) {
	v := fibonacci(50)
	assertauto.Equal(t, v)
}

func TestFibonacciString(t *testing.T) {
	s := fibonacciString(50)
	assertauto.Equal(t, s)
}

func BenchmarkFibonacci(b *testing.B) {
	for b.Loop() {
		fibonacci(50)
	}
}

func BenchmarkFibonacciString(b *testing.B) {
	for b.Loop() {
		fibonacciString(50)
	}
}
