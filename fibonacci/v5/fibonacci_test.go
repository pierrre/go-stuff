package main

import (
	"errors"
	"io"
	"math/big"
	"testing"

	"github.com/pierrre/assert"
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

func TestFibonacciWriteError(t *testing.T) {
	w := &errorWriter{
		err: errors.New("error"),
	}
	err := FibonacciWrite(w, testIterations)
	assert.Error(t, err)
}

type errorWriter struct {
	err error
}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, e.err
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

func BenchmarkFibonacciWrite(b *testing.B) {
	for range b.N {
		_ = FibonacciWrite(io.Discard, testIterations)
	}
}
