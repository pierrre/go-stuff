package main

import "reflect"

func other() reflect.Value {
	return reflect.ValueOf(struct{}{}).Method(0)
}
