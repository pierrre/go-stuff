package main

import (
	"errors"
	"io"
	"testing"

	"github.com/pierrre/assert"
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

func TestFibonacciWriteError(t *testing.T) {
	w := &errorWriter{
		err: errors.New("error"),
	}
	err := fibonacciWrite(w, testIterations)
	assert.Error(t, err)
}

type errorWriter struct {
	err error
}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, e.err
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

func BenchmarkFibonacciWrite(b *testing.B) {
	for b.Loop() {
		_ = fibonacciWrite(io.Discard, testIterations)
	}
}
