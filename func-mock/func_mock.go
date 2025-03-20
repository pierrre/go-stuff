package funcmock

import (
	"fmt"
	"reflect"
	"sync"
)

func mockList[FS ~[]F, F any](fs FS, onTooManyCalls F) F {
	count := 0
	rOnTooManyCalls := reflect.ValueOf(onTooManyCalls)
	return makeFunc[F](func(args []reflect.Value) []reflect.Value {
		count++
		var v reflect.Value
		if count <= len(fs) {
			f := fs[count-1]
			v = reflect.ValueOf(f)
		} else if !rOnTooManyCalls.IsNil() {
			v = rOnTooManyCalls
		} else {
			panic(fmt.Sprintf("too many calls: got %d, max %d", count, len(fs)))
		}
		return call(v, args)
	})
}

func mockCount[F any](f F) (F, func() int) {
	count := 0
	v := reflect.ValueOf(f)
	f = makeFunc[F](func(args []reflect.Value) []reflect.Value {
		defer func() {
			count++
		}()
		return call(v, args)
	})
	getCount := func() int {
		return count
	}
	return f, getCount
}

func mockSerial[F any](f F) F {
	mu := new(sync.Mutex)
	v := reflect.ValueOf(f)
	return makeFunc[F](func(args []reflect.Value) []reflect.Value {
		mu.Lock()
		defer mu.Unlock()
		return call(v, args)
	})
}

func makeFunc[F any](fn func(args []reflect.Value) []reflect.Value) F {
	typ := checkType[F]()
	v := reflect.MakeFunc(typ, fn)
	f := v.Interface().(F) //nolint:forcetypeassert // We know the type.
	return f
}

func checkType[F any]() reflect.Type {
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
