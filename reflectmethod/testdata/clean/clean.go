package clean

import "reflect"

func Kind(v any) reflect.Kind {
	return reflect.ValueOf(v).Kind()
}
