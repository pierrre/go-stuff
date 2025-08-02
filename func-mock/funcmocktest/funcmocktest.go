// Package funcmocktest provides utilities for mocking functions in tests.
package funcmocktest

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	funcmock "github.com/pierrre/go-stuff/func-mock"
)

// MockList creates a mock function that calls the functions in the list in order.
// If the list is exhausted, the test fails.
// If not all functions are called, the test fails.
func MockList[FS ~[]F, F funcmock.Func](tb testing.TB, fs FS, opts ...assert.Option) F {
	tb.Helper()
	onTooManyCalls := funcmock.MakeFunc[F](func(args []reflect.Value) []reflect.Value {
		tb.Helper()
		assert.Fail(tb, "funcmocktest.MockList", fmt.Sprintf("too many calls: max %d", len(fs)), 1, opts...)
		return funcmock.MakeFuncResults[F]()
	})
	f := funcmock.MockList(fs, onTooManyCalls)
	f, getCount := funcmock.MockCount(f)
	tb.Cleanup(func() {
		tb.Helper()
		count := getCount()
		if count < int64(len(fs)) {
			assert.Fail(tb, "funcmocktest.MockList", fmt.Sprintf("remaining calls: got %d, want %d", count, len(fs)), 1, opts...)
		}
	})
	return f
}

// MockCount creates a mock function that counts the number of calls.
// If the number of calls does not match the expected count, the test fails.
func MockCount[F funcmock.Func](tb testing.TB, f F, expected int64, opts ...assert.Option) F {
	tb.Helper()
	f, getCount := funcmock.MockCount(f)
	tb.Cleanup(func() {
		tb.Helper()
		count := getCount()
		if count != expected {
			assert.Fail(tb, "funcmocktest.MockCount", fmt.Sprintf("unexpected calls: got %d, want %d", count, expected), 1, opts...)
		}
	})
	return f
}
