//nolint:forbidigo // Print text.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	s := fibonacciString(50)
	fmt.Println(s)
}

func fibonacci(n int) int64 {
	a, b := int64(0), int64(1)
	for range n {
		a, b = b, a+b
	}
	return a
}

func fibonacciString(n int) string {
	s := ""
	sb := new(strings.Builder)
	for i := range n {
		v := fibonacci(i)
		sb.WriteString(strconv.FormatInt(v, 10) + "\n")
	}
	s += sb.String()
	return s
}
