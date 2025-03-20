package funcmock

import (
	"fmt"
	"reflect"
)

func Mock[FS ~[]F, F any](fs FS) F {
	typ := reflect.TypeFor[F]()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("unexpected kind: got %v, want %v", typ.Kind(), reflect.Func))
	}
	count := 0
	v := reflect.MakeFunc(typ, func(args []reflect.Value) (results []reflect.Value) {
		count++
		if count > len(fs) {
			panic(fmt.Sprintf("too many calls: got %d, max %d", count, len(fs)))
		}
		f := fs[count-1]
		v := reflect.ValueOf(f)
		return v.Call(args)
	})
	return v.Interface().(F) //nolint:forcetypeassert // We know the type.
}
