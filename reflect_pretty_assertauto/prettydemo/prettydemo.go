//nolint:forbidigo // Print text.
package main

import (
	"fmt"

	"github.com/pierrre/pretty"
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
		fmt.Println(pretty.String(v))
	}
}

type myStruct struct {
	I int
	S string
}
