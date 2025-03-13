package main

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
)

type myStruct struct {
	Int     int
	Float64 float64
	String  string
	Map     map[string]any
}

func myFunction() *myStruct {
	return &myStruct{
		Int:     123,
		Float64: 123.456,
		String:  "test",
		Map: map[string]any{
			"foo": "bar",
		},
	}
}

func Test(t *testing.T) {
	v := myFunction()
	/// assert.Equal(t, v, expected)
	assertauto.Equal(t, v)
}
