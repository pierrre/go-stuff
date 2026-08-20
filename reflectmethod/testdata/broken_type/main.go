package main

import "reflect"

func main() {
	reflect.ValueOf(struct{}{}).Method(0)
	var i int = "not an int"
	_ = i
}
