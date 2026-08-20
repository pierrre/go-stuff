package main

import (
	"reflect"

	"example.com/dep"
	"example.com/helper"
)

// reflect.Value.Method in a comment must not be detected.
// reflect.Value.MethodByName in a comment must not be detected.
// reflect.Type.Method in a comment must not be detected.
// reflect.Type.MethodByName in a comment must not be detected.
const commentOnly = "reflect.Value.Method and reflect.Type.Method in a string"

type myType struct{}

func (myType) Method(i int) int {
	return 0
}

func (myType) MethodByName(name string) int {
	return 0
}

type fieldHolder struct {
	Method       int
	MethodByName int
}

func takesFunc(f func(int) reflect.Value) {}

var packageLevel = reflect.ValueOf(myType{}).Method(0)

var methodValue = reflect.ValueOf(myType{}).Method

func main() {
	dep.Call("A")
	dep.First()
	v := reflect.ValueOf(myType{})
	v.Method(0)
	v.MethodByName("A")
	t := reflect.TypeOf(myType{})
	t.Method(0)
	t.MethodByName("B")
	pv := reflect.ValueOf(myType{})
	(&pv).Method(0)
	reflect.Value.MethodByName(reflect.ValueOf(myType{}), "D")
	myType{}.Method(0)
	myType{}.MethodByName("C")
	func() {
		reflect.ValueOf(myType{}).Method(0)
	}()
	(reflect.Value.Method)(reflect.ValueOf(myType{}), 0)
	fh := fieldHolder{}
	_ = fh.Method
	_ = fh.MethodByName
	helper.Method(0)
	helper.MethodByName("x")
	takesFunc(v.Method)
	_ = v.Method
	_ = v.MethodByName
	_ = commentOnly
}
