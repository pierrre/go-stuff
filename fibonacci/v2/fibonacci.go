//nolint:forbidigo // Print text.
package main

import (
	"fmt"
	"math/big"
)

func main() {
	s := FibonacciString(100)
	fmt.Println(s)
}

func Fibonacci(n int) *big.Int {
	a, b := big.NewInt(0), big.NewInt(1)
	for range n {
		a.Add(a, b)
		a, b = b, a
	}
	return a
}

func FibonacciString(n int) string {
	s := ""
	for i := range n {
		v := Fibonacci(i)
		vs := v.String()
		s += vs + "\n"
	}
	return s
}
