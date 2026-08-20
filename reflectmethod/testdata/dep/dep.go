package dep

import "reflect"

func Call(name string) reflect.Value {
	return reflect.ValueOf(struct{}{}).MethodByName(name)
}

func First() reflect.Value {
	return reflect.ValueOf(struct{}{}).Method(0)
}
