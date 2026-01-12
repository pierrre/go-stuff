//nolint:forbidigo // Print text.
package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	vs := []any{
		123,
		123.456,
		"test",
		[]int{1, 2, 3},
		myStruct{
			I: 123,
			S: "test",
		},
	}
	for _, v := range vs {
		rv := reflect.ValueOf(v)
		t := rv.Type()
		k := t.Kind()
		var s string
		switch k {
		case reflect.Int:
			s = strconv.FormatInt(rv.Int(), 10)
		case reflect.Float64:
			s = strconv.FormatFloat(rv.Float(), 'g', -1, 64)
		case reflect.String:
			s = rv.String()
		default:
			s = "unsupported"
		}
		fmt.Printf("type=%s, kind=%s, string=%s\n", t.String(), k.String(), s)
	}
}

type myStruct struct {
	I int
	S string
}
