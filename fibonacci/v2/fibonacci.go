//nolint:forbidigo // Print text.
package main

import (
	"fmt"
	"math/big"
)

func main() {
	s := fibonacciString(100)
	fmt.Println(s)
}

func fibonacci(n int) *big.Int {
	a, b := big.NewInt(0), big.NewInt(1)
	for range n {
		a.Add(a, b)
		a, b = b, a
	}
	return a
}

func fibonacciString(n int) string {
	s := ""
	for i := range n {
		v := fibonacci(i)
		vs := v.String()
		s += vs + "\n"
	}
	return s
}
