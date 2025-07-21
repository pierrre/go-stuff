// Package funcmock provides utilities for mocking functions.
package funcmock

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

// Func represents a function type.
type Func any

// MockList creates a mock function that calls the functions in the list in order.
// If the list is exhausted, it calls onTooManyCalls if it's not nil, otherwise it panics.
func MockList[FS ~[]F, F Func](fs FS, onTooManyCalls F) F {
	var count int64
	rOnTooManyCalls := reflect.ValueOf(onTooManyCalls)
	return MakeFunc[F](func(args []reflect.Value) []reflect.Value {
		c := atomic.AddInt64(&count, 1)
		var v reflect.Value
		switch {
		case c <= int64(len(fs)):
			f := fs[c-1]
			v = reflect.ValueOf(f)
		case !rOnTooManyCalls.IsNil():
			v = rOnTooManyCalls
		default:
			panic(fmt.Sprintf("too many calls: got %d, max %d", c, len(fs)))
		}
		return call(v, args)
	})
}

// MockCount creates a mock function that counts the number of calls.
// It returns the mock function and a function to get the count.
// The count is incremented after each call to the mocked function returns.
func MockCount[F Func](f F) (_ F, getCount func() int64) {
	var count int64
	v := reflect.ValueOf(f)
	f = MakeFunc[F](func(args []reflect.Value) []reflect.Value {
		defer atomic.AddInt64(&count, 1)
		return call(v, args)
	})
	getCount = func() int64 {
		return atomic.LoadInt64(&count)
	}
	return f, getCount
}

// MockSerial creates a mock function that serializes calls to the function.
// It uses a mutex to ensure that only one call is processed at a time.
func MockSerial[F Func](f F) F {
	mu := new(sync.Mutex)
	v := reflect.ValueOf(f)
	return MakeFunc[F](func(args []reflect.Value) []reflect.Value {
		mu.Lock()
		defer mu.Unlock()
		return call(v, args)
	})
}

// MakeFunc creates a mock function that calls the provided function fn.
func MakeFunc[F Func](fn func(args []reflect.Value) []reflect.Value) F {
	typ := checkType[F]()
	v := reflect.MakeFunc(typ, fn)
	f, _ := v.Interface().(F)
	return f
}

// MakeFuncResults creates a slice of zero values for the return values of the function type F.
func MakeFuncResults[F Func]() []reflect.Value {
	typ := checkType[F]()
	n := typ.NumOut()
	var results []reflect.Value
	if n > 0 {
		results = make([]reflect.Value, typ.NumOut())
		for i := range n {
			results[i] = reflect.Zero(typ.Out(i))
		}
	}
	return results
}

func checkType[F Func]() reflect.Type {
	typ := reflect.TypeFor[F]()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("unexpected kind: got %v, want %v", typ.Kind(), reflect.Func))
	}
	return typ
}

func call(v reflect.Value, args []reflect.Value) []reflect.Value {
	if v.Type().IsVariadic() {
		return v.CallSlice(args)
	} else {
		return v.Call(args)
	}
}
